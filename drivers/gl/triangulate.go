// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"container/list"
	"fmt"
	"github.com/google/gxui/math"
)

const debugTriangulate = false
const epsilon = 0.00001

func isConcave(edges []math.Vec2, a, b, c int) bool {
	return edges[b].Sub(edges[a]).Cross(edges[b].Sub(edges[c])) > -epsilon
}

func isEar(edges []math.Vec2, a, b, c int) bool {
	if isConcave(edges, a, b, c) {
		if debugTriangulate {
			fmt.Printf("concave: %c, %c, %c\n", 'A'+a, 'A'+b, 'A'+c)
		}
		return false
	}

	plane := [3]math.Vec3{}
	for i := 0; i < 3; i++ {
		p := [3]int{a, b, c}[i]
		q := [3]int{b, c, a}[i]
		normal := edges[q].Sub(edges[p]).Normalize().Tangent()
		plane[i] = normal.Vec3(-normal.Dot(edges[p]))
	}

	for i := 0; i < len(edges); i++ {
		if i == a || i == b || i == c {
			continue
		}
		v := edges[i].Vec3(1)

		if v.Dot(plane[0]) > -epsilon &&
			v.Dot(plane[1]) > -epsilon &&
			v.Dot(plane[2]) > -epsilon {
			if debugTriangulate {
				fmt.Printf("non-ear: %c, %c, %c (%c %f:%f:%f)\n",
					'A'+a, 'A'+b, 'A'+c,
					'A'+i,
					v.Dot(plane[0]), v.Dot(plane[1]), v.Dot(plane[2]),
				)
			}
			return false
		}
	}
	if debugTriangulate {
		fmt.Printf("ear: %c, %c, %c\n",
			'A'+a, 'A'+b, 'A'+c,
		)
	}
	return true
}

func pruneEdgeDuplicates(edges []math.Vec2) []math.Vec2 {
	pruned := make([]math.Vec2, 0, len(edges))
	last := math.Vec2{}
	for i, v := range edges {
		if i == 0 || last.Sub(v).Len() > 0.0001 {
			pruned = append(pruned, v)
		}
		last = v
	}
	return pruned
}

func triangulate(edges []math.Vec2) []math.Vec2 {
	if debugTriangulate {
		fmt.Printf("triangulate()\n")
	}

	edges = pruneEdgeDuplicates(edges)
	if len(edges) < 3 {
		return []math.Vec2{}
	}

	l := list.New()
	for i := range edges {
		l.PushBack(i)
	}

	out := []math.Vec2{}

	pruned := true
	for pruned {
		pruned = false
		for a := l.Front(); a != nil; {
			b := a.Next()
			if b == nil {
				b = l.Front()
			}
			c := b.Next()
			if c == nil {
				c = l.Front()
			}

			ia, ib, ic := a.Value.(int), b.Value.(int), c.Value.(int)

			if isEar(edges, ia, ib, ic) {
				l.Remove(b)
				out = append(out, edges[ia], edges[ib], edges[ic])
				pruned = true
				if l.Len() < 3 {
					return out
				}
			} else {
				a = a.Next()
			}
		}
	}

	if debugTriangulate {
		fmt.Printf("Failed to prune an ear! edges: %#v, out: %v\n", edges, out)
	}

	// assert.True(l.Len() < 3, "Failed to prune an ear! edges: %#v, out: %v", edges, out)
	return out
}
