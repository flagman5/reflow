// Copyright 2017 GRAIL, Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package reflow_test

import (
	"reflect"
	"testing"

	"github.com/grailbio/reflow"
)

func TestResources(t *testing.T) {
	r1 := reflow.Resources{10, 5, 1}
	r2 := reflow.Resources{5, 2, 3}
	if got, want := r1.Sub(r2), (reflow.Resources{5, 3, 0}); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r1.Add(r2), (reflow.Resources{15, 7, 4}); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r1.Add(r2), r2.Add(r1); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if r1.Available(r2) {
		t.Errorf("expected %v to be unavailable in %v", r2, r1)
	}
	r3 := reflow.Resources{3, 1, 1}
	if !r1.Available(r3) {
		t.Errorf("expected %v to be available in %v", r3, r1)
	}

	if got, want := r1.Units(r2), 3; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r2.Units(r1), 3; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := reflow.MaxResources.Add(reflow.MaxResources), reflow.MaxResources; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r1.Min(r2), (reflow.Resources{5, 2, 1}); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r1.Max(r2), (reflow.Resources{10, 5, 3}); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestRequirements(t *testing.T) {
	var req reflow.Requirements
	req = req.Add(reflow.Resources{10, 5, 1})
	req = req.Add(reflow.Resources{20, 3, 1})
	if got, want := req.Min, (reflow.Resources{20, 5, 1}); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

var (
	file1 = reflow.File{reflow.Digester.FromString("foo"), 3}
	file2 = reflow.File{reflow.Digester.FromString("bar"), 3}
	file3 = reflow.File{reflow.Digester.FromString("a/b/c"), 5}

	v1 = reflow.Fileset{Map: map[string]reflow.File{
		"foo": file1,
		"bar": file2,
	}}
	v2 = reflow.Fileset{Map: map[string]reflow.File{
		"a/b/c": file3,
		"bar":   file2,
	}}
	vlist = reflow.Fileset{List: []reflow.Fileset{v1, v2}}
)

const vlistSHA256 = "sha256:d60e67ce9e89548b502a5ad7968e99caed0d388f0a991b906f41a7ba65adb31f"

func TestValueDigest(t *testing.T) {
	if v1.Digest() == v2.Digest() {
		t.Errorf("did not expect v1, v2 to have same digest")
	}
	if got, want := vlist.Digest().String(), vlistSHA256; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestValueFile(t *testing.T) {
	files := vlist.Files()
	expected := map[reflow.File]bool{file1: true, file2: true, file3: true}
	for _, f := range files {
		if !expected[f] {
			t.Errorf("unexpected file %v", f)
		}
		delete(expected, f)
	}
	if len(expected) != 0 {
		t.Errorf("expected additional files %v", expected)
	}
}

func TestValuePullup(t *testing.T) {
	got := vlist.Pullup()
	want := reflow.Fileset{Map: map[string]reflow.File{"foo": file1, "bar": file2, "a/b/c": file3}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestValueN(t *testing.T) {
	if got, want := vlist.N(), 4; got != want {
		t.Errorf("got %v, want %v")
	}
}

func TestEmpty(t *testing.T) {
	empty := []reflow.Fileset{
		{},
		{List: make([]reflow.Fileset, 1)},
		{List: []reflow.Fileset{{}, {Map: map[string]reflow.File{}}, {List: make([]reflow.Fileset, 100)}}},
	}
	for i, fs := range empty {
		if !fs.Empty() {
			t.Errorf("expected empty %d %v", i, fs)
		}
	}
}

func TestAnyEmpty(t *testing.T) {
	empty := []reflow.Fileset{
		{},
		{List: make([]reflow.Fileset, 1)},
		{List: []reflow.Fileset{{}, {Map: map[string]reflow.File{}}, {List: make([]reflow.Fileset, 100)}}},
		{List: []reflow.Fileset{{Map: map[string]reflow.File{".": reflow.File{}}}, {}}},
	}
	for i, fs := range empty {
		if !fs.AnyEmpty() {
			t.Errorf("expected anyempty %d %v", i, fs)
		}
	}
}
