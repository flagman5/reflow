// Copyright 2017 GRAIL, Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package types

import "testing"

var (
	ty1 = Struct(
		&Field{Name: "a", T: Int},
		&Field{Name: "b", T: String},
		&Field{Name: "c", T: Tuple(&Field{T: Int}, &Field{T: String})})
	ty2 = Struct(
		&Field{Name: "a", T: Int},
		&Field{Name: "d", T: String},
		&Field{Name: "c", T: Tuple(&Field{T: Int}, &Field{T: String})})
	ty12 = Struct(
		&Field{Name: "a", T: Int},
		&Field{Name: "b", T: String},
		&Field{Name: "c", T: Tuple(&Field{T: Int}, &Field{T: String})},
		&Field{Name: "d", T: String})

	mty1 = Module(
		[]*Field{
			&Field{Name: "a", T: Int},
			&Field{Name: "b", T: String},
			&Field{Name: "c", T: Tuple(&Field{T: Int}, &Field{T: String})},
		},
		nil,
	)
	mty2 = Module(
		[]*Field{
			&Field{Name: "a", T: Int},
			&Field{Name: "d", T: String},
			&Field{Name: "c", T: Tuple(&Field{T: Int}, &Field{T: String})},
		},
		nil,
	)
	mty12 = Module(
		[]*Field{
			&Field{Name: "a", T: Int},
			&Field{Name: "b", T: String},
			&Field{Name: "c", T: Tuple(&Field{T: Int}, &Field{T: String})},
			&Field{Name: "d", T: String},
		},
		nil,
	)

	mapTy  = Map(Int, String)
	listTy = List(String)

	types = []*T{ty1, ty2, ty12, mty1, mty2, mty12}
)

func TestUnify(t *testing.T) {
	u := ty1.Unify(ty2)
	if got, want := u.String(), "{a int, c (int, string)}"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	u = mty1.Unify(mty2)
	if got, want := u.String(), "module{a int, c (int, string)}"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := listTy.Unify(List(Bottom)), listTy; !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := mapTy.Unify(Map(Top, Bottom)), mapTy; !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := Map(Int, ty1).Unify(Map(Int, ty2)).String(), `[int:{a int, c (int, string)}]`; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSub(t *testing.T) {
	if ty1.Sub(ty2) {
		t.Errorf("%s is not a subtype of %s", ty1, ty2)
	}
	if !ty12.Sub(ty1) {
		t.Errorf("%s is a subtype of %s", ty12, ty1)
	}
	if !ty12.Sub(ty2) {
		t.Errorf("%s is a subtype of %s", ty12, ty2)
	}
	if mty1.Sub(mty2) {
		t.Errorf("%s is not a subtype of %s", ty1, ty2)
	}
	if !mty12.Sub(mty1) {
		t.Errorf("%s is a subtype of %s", ty12, ty1)
	}
	if !mty12.Sub(mty2) {
		t.Errorf("%s is a subtype of %s", ty12, ty2)
	}
	for _, typ := range types {
		if !Bottom.Sub(typ) {
			t.Errorf("bottom is a subtype of %v", typ)
		}
	}
	for _, typ := range types {
		if !typ.Sub(Top) {
			t.Errorf("%s is a subtype of top", typ)
		}
	}
}

func TestEnv(t *testing.T) {
	e := NewEnv()
	e.Bind("a", Int)
	e.Bind("b", String)
	e, save := e.Push(), e
	e.Bind("a", String)
	if got, want := e.Type("a"), String; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := e.Type("b"), String; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	e = save
	if got, want := e.Type("a"), Int; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := e.Type("b"), String; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
