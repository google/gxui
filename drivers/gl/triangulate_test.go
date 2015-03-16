// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import test "github.com/google/gxui/testing"
import (
	"github.com/google/gxui/math"
	"testing"
)

func v(a, b float32) math.Vec2 {
	return math.Vec2{X: a, Y: b}
}

func TestIsConcave(t *testing.T) {
	/*
	   B
	   A  C
	*/
	A := v(0, 1)
	B := v(0, 0)
	C := v(1, 1)
	edges := []math.Vec2{A, B, C}
	test.AssertEquals(t, false, isConcave(edges, 0, 1, 2))
	test.AssertEquals(t, false, isConcave(edges, 1, 2, 0))
	test.AssertEquals(t, false, isConcave(edges, 2, 0, 1))
	test.AssertEquals(t, true, isConcave(edges, 2, 1, 0))
	test.AssertEquals(t, true, isConcave(edges, 0, 2, 1))
	test.AssertEquals(t, true, isConcave(edges, 1, 0, 2))
}

func TestTriangluate3Verts(t *testing.T) {
	/*
	   B
	   A  C
	*/
	A := v(0, 1)
	B := v(0, 0)
	C := v(1, 1)
	edges := []math.Vec2{A, B, C}
	test.AssertEquals(t, edges, triangulate(edges))
}

func TestTriangluateQuad(t *testing.T) {
	/*
	   A  B
	   D  C
	*/
	A := v(0, 0)
	B := v(1, 0)
	C := v(1, 1)
	D := v(0, 1)
	edges := []math.Vec2{A, B, C, D}
	tris := []math.Vec2{
		A, B, C,
		A, C, D,
	}
	test.AssertEquals(t, tris, triangulate(edges))
}

func TestTriangluateDupeVertex(t *testing.T) {
	/*
	   A  B
	   D  C
	*/
	A := v(0, 0)
	B := v(1, 0)
	C := v(1, 1)
	D := v(0, 1)
	edges := []math.Vec2{A, B, B, C, C, D}
	tris := []math.Vec2{
		A, B, C,
		A, C, D,
	}
	test.AssertEquals(t, tris, triangulate(edges))
}
func TestTriangluateConcave(t *testing.T) {
	/*
	         A
	     H       B

	   G           C

	     F       D
	         E
	*/
	A := v(0, -3)
	B := v(4, -2)
	C := v(6, -0)
	D := v(4, 2)
	E := v(0, 3)
	F := v(-4, 2)
	G := v(-6, -0)
	H := v(-4, -2)
	edges := []math.Vec2{A, B, C, D, E, F, G, H}
	tris := []math.Vec2{
		A, B, C,
		A, C, D,
		A, D, E,
		A, E, F,
		A, F, G,
		A, G, H,
	}
	test.AssertEquals(t, tris, triangulate(edges))
}

func TestTriangluateConvex(t *testing.T) {
	/*
	       D-------E
	       |       |
	   B---C   G   |
	   |       | \ |
	   A-------H   F
	*/
	A := v(0, 2)
	B := v(0, 1)
	C := v(1, 1)
	D := v(1, 0)
	E := v(3, 0)
	F := v(3, 2)
	G := v(2, 1)
	H := v(2, 2)

	edges := []math.Vec2{A, B, C, D, E, F, G, H}
	tris := []math.Vec2{
		A, B, C,
		C, D, E,
		E, F, G,
		G, H, A,
		G, A, C,
		G, C, E,
	}
	test.AssertEquals(t, tris, triangulate(edges))
}

func TestTriangluateConvex2(t *testing.T) {
	/*
	   A-----------------B
	   |                 |
	   |                 |
	   F-------------E   C
	                  \ /
	                   D
	*/

	A := v(462, 531)
	B := v(576, 531)
	C := v(576, 557)
	D := v(573, 572)
	E := v(566, 557)
	F := v(462, 557)

	edges := []math.Vec2{A, B, C, C, D, E, F}
	tris := []math.Vec2{
		A, B, C,
		C, D, E,
		E, F, A,
		E, A, C,
	}
	test.AssertEquals(t, tris, triangulate(edges))
}

func TestTriangluateConvex3(t *testing.T) {
	/*
	    B---------------C
	   A                 D
	   |            F    |
	   |           G  \  |
	   J             \   E
	    I--------------H
	*/

	A := v(8.000002, 132)
	B := v(11.514719, 123.514)
	C := v(404, 120)
	D := v(412.4853, 123.51472)
	E := v(416, 176.98395)
	F := v(300.47543, 139.39261)
	G := v(290.65604, 143.82726)
	H := v(392.73557, 208)
	I := v(20, 208)
	J := v(11.514719, 204.48528)

	edges := []math.Vec2{A, B, C, C, D, E, F, G, H, I, J}
	tris := []math.Vec2{
		A, B, C,
		A, C, D,
		D, E, F,
		G, H, I,
		G, I, J,
		G, J, A,
		A, D, F,
		A, F, G,
	}
	test.AssertEquals(t, tris, triangulate(edges))
}

func TestTriangluateConvex4(t *testing.T) {
	/*
	    J---------------A
	   I                 B
	   |            D    |
	   |           E  \  |
	   H             \   C
	    G--------------F
	*/

	A := v(404, 120)
	B := v(412.4853, 123.51472)
	C := v(416, 177.90715)
	D := v(280.82632, 146.21126)
	E := v(271.39908, 151.5048)
	F := v(384.94205, 208)
	G := v(20, 208)
	H := v(11.514719, 204.48528)
	I := v(8.000002, 132)
	J := v(11.514719, 123.51472)

	edges := []math.Vec2{A, B, C, C, D, E, F, G, H, I, J}
	tris := []math.Vec2{
		A, B, C,
		A, C, D,
		E, F, G,
		E, G, H,
		E, H, I,
		E, I, J,
		J, A, D,
		J, D, E,
	}
	test.AssertEquals(t, tris, triangulate(edges))
}

func TestTriangulateChevron(t *testing.T) {
	// A-----B
	//  \     \
	//   \     \
	//    F     C
	//   /     /
	//  /     /
	// E-----D
	A := v(0, 0)
	B := v(10, 0)
	C := v(15, 5)
	D := v(10, 10)
	E := v(0, 10)
	F := v(5, 5)

	edges := []math.Vec2{A, B, C, D, E, F}
	tris := []math.Vec2{
		A, B, C,
		C, D, E,
		C, E, F,
		C, F, A,
	}
	test.AssertEquals(t, tris, triangulate(edges))
}
