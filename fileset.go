// Copyright 2018 GRAIL, Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package reflow

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/grailbio/base/data"
	"github.com/grailbio/base/digest"
)

// File represents a File inside of Reflow. A file is said to be
// resolved if it contains the digest of the file's contents.
// Otherwise, a File is said to be a reference, in which case it must
// contain a source and etag.
type File struct {
	// The digest of the contents of the file.
	ID digest.Digest

	// The size of the file.
	Size int64

	// Source stores a URL for the file from which it may
	// be retrieved.
	Source string `json:",omitempty"`

	// ETag stores an optional entity tag for the Source file.
	ETag string `json:",omitempty"`

	// LastModified stores the file's last modified time.
	LastModified time.Time
}

// Digest returns the file's digest: if the file is a reference, the
// digest comprises the reference, source, and etag. Resolved files
// return the file's digest.
func (f File) Digest() digest.Digest {
	if !f.IsRef() {
		return f.ID
	}
	w := Digester.NewWriter()
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(f.Size))
	w.Write(b[:])
	io.WriteString(w, f.Source)
	io.WriteString(w, f.ETag)
	return w.Digest()
}

// Equal returns whether files f and g represent the same content.
func (f File) Equal(g File) bool {
	if f.IsRef() || g.IsRef() {
		// When we compare references, we require nonempty etags.
		return f.Size == g.Size && f.Source == g.Source && f.ETag != "" && f.ETag == g.ETag
	}
	return f.ID == g.ID
}

// IsRef returns whether this file is a file reference.
func (f File) IsRef() bool {
	return f.ID.IsZero()
}

func (f File) String() string {
	var b strings.Builder
	if !f.IsRef() {
		fmt.Fprintf(&b, "id: %s", f.ID)
	}
	maybeComma(&b)
	fmt.Fprintf(&b, "size: %d", f.Size)
	if f.Source != "" {
		maybeComma(&b)
		fmt.Fprintf(&b, "source: %s", f.Source)
	}
	if f.ETag != "" {
		maybeComma(&b)
		fmt.Fprintf(&b, "etag: %vs", f.ETag)
	}
	return b.String()
}

func (f File) Short() string {
	if f.IsRef() {
		return f.Source
	} else {
		return f.ID.Short()
	}
}

// Fileset is the result of an evaluated flow. Values may either be
// lists of values or Filesets. Filesets are a map of paths to Files.
type Fileset struct {
	List []Fileset       `json:",omitempty"`
	Map  map[string]File `json:"Fileset,omitempty"`
}

// Equal reports whether v is equal to w.
func (v Fileset) Equal(w Fileset) bool {
	if len(v.List) != len(w.List) {
		return false
	}
	for i := range v.List {
		if !v.List[i].Equal(w.List[i]) {
			return false
		}
	}
	if len(v.Map) != len(w.Map) {
		return false
	}
	for key, f := range v.Map {
		g, ok := w.Map[key]
		if !ok || !f.Equal(g) {
			return false
		}
	}
	return true
}

// Empty tells whether this value is empty, that is, it contains
// no files.
func (v Fileset) Empty() bool {
	for _, fs := range v.List {
		if !fs.Empty() {
			return false
		}
	}
	return len(v.Map) == 0
}

// AnyEmpty tells whether this value, or any of its constituent
// values contain no files.
func (v Fileset) AnyEmpty() bool {
	for _, fs := range v.List {
		if fs.AnyEmpty() {
			return true
		}
	}
	return len(v.List) == 0 && len(v.Map) == 0
}

// Flatten is a convenience function to flatten (shallowly) the value
// v, returning a list of Values. If the value is a list value, the
// list is returned; otherwise a unary list of the value v is
// returned.
func (v Fileset) Flatten() []Fileset {
	switch {
	case v.List != nil:
		return v.List
	default:
		return []Fileset{v}
	}
}

// Files returns the set of Files that comprise the value.
func (v Fileset) Files() []File {
	fs := map[File]bool{}
	v.files(fs)
	files := make([]File, len(fs))
	i := 0
	for f := range fs {
		files[i] = f
		i++
	}
	return files
}

// N returns the number of files (not necessarily unique) in this value.
func (v Fileset) N() int {
	var n int
	for _, v := range v.List {
		n += v.N()
	}
	n += len(v.Map)
	return n
}

// Size returns the total size of this value.
func (v Fileset) Size() int64 {
	var s int64
	for _, v := range v.List {
		s += v.Size()
	}
	for _, f := range v.Map {
		s += f.Size
	}
	return s
}

// Subst fileset v with references files substituted using the
// provided map. Subst returns whether the fileset is fully resolved
// after substitution.
func (v Fileset) Subst(sub map[File]File) (out Fileset, resolved bool) {
	resolved = true
	if v.List != nil {
		out.List = make([]Fileset, len(v.List))
		for i := range out.List {
			var ok bool
			out.List[i], ok = v.List[i].Subst(sub)
			if !ok {
				resolved = false
			}
		}
	}
	if v.Map != nil {
		out.Map = make(map[string]File, len(v.Map))
		for path, file := range v.Map {
			if f, ok := sub[file]; ok {
				out.Map[path] = f
			} else {
				out.Map[path] = file
				if file.IsRef() {
					resolved = false
				}
			}
		}
	}
	return
}

func (v Fileset) files(fs map[File]bool) {
	for i := range v.List {
		v.List[i].files(fs)
	}
	if v.Map != nil {
		for _, f := range v.Map {
			fs[f] = true
		}
	}
}

// Short returns a short, human-readable string representing the
// value. Its intended use is for pretty-printed output. In
// particular, hashes are abbreviated, and lists display only the
// first member, followed by ellipsis. For example, a list of values
// is printed as:
//	list<val<sample.fastq.gz=f2c59c40>, ...50MB>
func (v Fileset) Short() string {
	switch {
	case v.List != nil:
		s := "list<"
		if len(v.List) != 0 {
			s += v.List[0].Short()
			if len(v.List) > 1 {
				s += ", ..."
			} else {
				s += " "
			}
		}
		s += data.Size(v.Size()).String() + ">"
		return s
	case len(v.Map) == 0:
		return "val<>"
	default:
		paths := make([]string, len(v.Map))
		i := 0
		for path := range v.Map {
			paths[i] = path
			i++
		}
		sort.Strings(paths)
		path := paths[0]
		file := v.Map[path]
		s := fmt.Sprintf("val<%s=%s", path, file.Short())
		if len(paths) > 1 {
			s += ", ..."
		} else {
			s += " "
		}
		s += data.Size(v.Size()).String() + ">"
		return s
	}
}

// String returns a full, human-readable string representing the value v.
// Unlike Short, string is fully descriptive: it contains the full digest and
// lists are complete. For example:
//	list<sample.fastq.gz=sha256:f2c59c40a1d71c0c2af12d38a2276d9df49073c08360d72320847efebc820160>,
//	  sample2.fastq.gz=sha256:59eb82c49448e349486b29540ad71f4ddd7f53e5a204d50997f054d05c939adb>>
func (v Fileset) String() string {
	switch {
	case v.List != nil:
		vals := make([]string, len(v.List))
		for i := range v.List {
			vals[i] = v.List[i].String()
		}
		return fmt.Sprintf("list<%s>", strings.Join(vals, ", "))
	case len(v.Map) == 0:
		return "void"
	default:
		// TODO(marius): should we include the bindings here?
		paths := make([]string, len(v.Map))
		i := 0
		for path := range v.Map {
			paths[i] = path
			i++
		}
		sort.Strings(paths)
		binds := make([]string, len(paths))
		for i, path := range paths {
			binds[i] = fmt.Sprintf("%s=<%s>", path, v.Map[path])
		}
		return fmt.Sprintf("obj<%s>", strings.Join(binds, ", "))
	}
}

// Digest returns a digest representing the value. Digests preserve
// semantics: two values with the same digest are considered to be
// equivalent.
func (v Fileset) Digest() digest.Digest {
	w := Digester.NewWriter()
	v.WriteDigest(w)
	return w.Digest()
}

// WriteDigest writes the digestible material for v to w. The
// io.Writer is assumed to be produced by a Digester, and hence
// infallible. Errors are not checked.
func (v Fileset) WriteDigest(w io.Writer) {
	switch {
	case v.List != nil:
		for i := range v.List {
			v.List[i].WriteDigest(w)
		}
	default:
		paths := make([]string, len(v.Map))
		i := 0
		for path := range v.Map {
			paths[i] = path
			i++
		}
		sort.Strings(paths)
		for _, path := range paths {
			io.WriteString(w, path)
			digest.WriteDigest(w, v.Map[path].Digest())
		}
	}
}

// Pullup merges this value (tree) into a single toplevel fileset.
func (v Fileset) Pullup() Fileset {
	if v.List == nil {
		return v
	}
	p := Fileset{Map: map[string]File{}}
	v.pullup(p.Map)
	return p
}

func (v Fileset) pullup(m map[string]File) {
	for k, f := range v.Map {
		m[k] = f
	}
	for _, v := range v.List {
		v.pullup(m)
	}
}

func maybeComma(b *strings.Builder) {
	if b.Len() > 0 {
		b.WriteString(", ")
	}
}
