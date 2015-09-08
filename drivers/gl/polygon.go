// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
)

func appendVec2(arr []float32, vecs ...math.Vec2) []float32 {
	for _, v := range vecs {
		arr = append(arr, v.X, v.Y)
	}
	return arr
}

func pruneDuplicates(p gxui.Polygon) gxui.Polygon {
	pruned := make(gxui.Polygon, 0, len(p))
	last := gxui.PolygonVertex{}
	for i, v := range p {
		if i == 0 || last.Position.Sub(v.Position).Vec2().Len() > 0.001 {
			pruned = append(pruned, v)
		}
		last = v
	}
	return pruned
}

func segment(penWidth, r float32, a, b, c math.Vec2, aIsLast bool, vsEdgePos []float32, fillEdge []math.Vec2) ([]float32, []math.Vec2) {
	ba, ca := a.Sub(b), a.Sub(c)
	baLen, caLen := ba.Len(), ca.Len()
	baDir, caDir := ba.DivS(baLen), ca.DivS(caLen)
	dp := baDir.Dot(caDir)
	if dp < -0.99999 {
		// Straight lines cause DBZs, special case
		inner := a.Sub(caDir.Tangent().MulS(penWidth))
		vsEdgePos = appendVec2(vsEdgePos, a, inner)
		if fillEdge != nil /*&& i != 0*/ {
			fillEdge = append(fillEdge, inner)
		}
		return vsEdgePos, fillEdge
	}
	α := math.Acosf(dp) / 2
	// ╔═══════════════════════════╦════════════════╗
	// ║                           ║                ║
	// ║             A             ║                ║
	// ║            ╱:╲            ║                ║
	// ║           ╱α:α╲           ║   A            ║
	// ║          ╱  :  ╲          ║   |╲           ║
	// ║         ╱ . d . ╲         ║   |α╲          ║
	// ║        .    :    .        ║   |  ╲         ║
	// ║       .P    :    Q.       ║   |   ╲        ║
	// ║      ╱      X      ╲      ║   |    ╲       ║
	// ║     ╱ .     ┊     . ╲     ║   |     ╲      ║
	// ║    ╱   .    r    .   ╲    ║   |      ╲     ║
	// ║   ╱       . ┊ .       ╲   ║   |┐     β╲    ║
	// ║  B          ┊          C  ║   P————————X   ║
	// ║                           ║                ║
	// ║             ^             ║                ║
	// ║             ┊v            ║                ║
	// ║             ┊  u          ║                ║
	// ║             ┊—————>       ║                ║
	// ║                           ║                ║
	// ╚═══════════════════════════╩════════════════╝
	v := baDir.Add(caDir).Normalize()
	u := v.Tangent()
	//
	// cos(2 • α) = dp
	//
	//      cos⁻¹(dp)
	// α = ───────────
	//          2
	//
	//           r
	// sin(α) = ───
	//           d
	//
	//       r
	// d = ──────
	//     sin(α)
	//
	d := r / math.Sinf(α)

	// X cannot be futher than half way along ab or ac
	dMax := math.Minf(baLen, caLen) / (2 * math.Cosf(α))
	if d > dMax {
		// Adjust d and r to compensate
		d = dMax
		r = d * math.Sinf(α)
	}

	x := a.Sub(v.MulS(d))

	convex := baDir.Tangent().Dot(caDir) <= 0

	w := penWidth
	β := math.Pi/2 - α

	// Special case for convex vertices where the pen width is greater than
	// the rounding. Without dealing with this, we'd end up with the inner
	// vertices overlapping. Instead use a point calculated much the same as
	// x, but using the pen width.
	useFixedInnerPoint := convex && w > r
	fixedInnerPoint := a.Sub(v.MulS(math.Minf(w/math.Sinf(α), dMax)))

	// Concave vertices behave much the same as convex, but we have to flip
	// β as the sweep is reversed and w as we're extruding.
	if !convex {
		w, β = -w, -β
	}

	steps := 1 + int(d*α)

	if aIsLast {
		// No curvy edge required for the last vertex.
		// This is already done by the first vertex.
		steps = 1
	}

	for j := 0; j < steps; j++ {
		γ := float32(0)
		if steps > 1 {
			γ = math.Lerpf(-β, β, float32(j)/float32(steps-1))
		}

		dir := v.MulS(math.Cosf(γ)).Add(u.MulS(math.Sinf(γ)))
		va := x.Add(dir.MulS(r))
		vb := va.Sub(dir.MulS(w))
		if useFixedInnerPoint {
			vb = fixedInnerPoint
		}

		vsEdgePos = appendVec2(vsEdgePos, va, vb)
		if fillEdge != nil {
			fillEdge = append(fillEdge, vb)
		}
	}

	return vsEdgePos, fillEdge
}

func closedPolyToShape(p gxui.Polygon, penWidth float32) (fillShape, edgeShape *shape) {
	p = pruneDuplicates(p)

	fillEdge := []math.Vec2{}
	vsEdgePos := []float32{}

	for i, cnt := 0, len(p); i < cnt; i++ {
		r := p[i].RoundedRadius
		a := p[i].Position.Vec2()
		b := p[(i+cnt-1)%cnt].Position.Vec2()
		c := p[(i+1)%cnt].Position.Vec2()
		vsEdgePos, fillEdge = segment(penWidth, r, a, b, c, i == len(p), vsEdgePos, fillEdge)
	}

	// Close the edge
	if len(vsEdgePos) >= 4 {
		vsEdgePos = append(vsEdgePos, vsEdgePos[:4]...)
	}

	fillTris := triangulate(fillEdge)
	if len(fillTris) > 0 {
		fillPos := make([]float32, len(fillTris)*2)
		for i, t := range fillTris {
			fillPos[i*2+0] = t.X
			fillPos[i*2+1] = t.Y
		}
		fillShape = newShape(newVertexBuffer(
			newVertexStream("aPosition", stFloatVec2, fillPos),
		), nil, dmTriangles)
	}

	if len(vsEdgePos) > 0 {
		edgeShape = newShape(newVertexBuffer(
			newVertexStream("aPosition", stFloatVec2, vsEdgePos),
		), nil, dmTriangleStrip)
	}

	return fillShape, edgeShape
}

func openPolyToShape(p gxui.Polygon, penWidth float32) *shape {
	p = pruneDuplicates(p)
	if len(p) < 2 {
		return nil
	}

	vsEdgePos := []float32{}

	{ // p[0] -> p[1]
		a, c := p[0].Position.Vec2(), p[1].Position.Vec2()
		caDir := a.Sub(c).Normalize()
		inner := a.Sub(caDir.Tangent().MulS(penWidth))
		vsEdgePos = appendVec2(vsEdgePos, a, inner)
	}
	for i := 1; i < len(p)-1; i++ {
		r := p[i].RoundedRadius
		a := p[i].Position.Vec2()
		b := p[i-1].Position.Vec2()
		c := p[i+1].Position.Vec2()
		vsEdgePos, _ = segment(penWidth, r, a, b, c, false, vsEdgePos, nil)
	}
	{ // p[N-2] -> p[N-1]
		a, c := p[len(p)-2].Position.Vec2(), p[len(p)-1].Position.Vec2()
		caDir := a.Sub(c).Normalize()
		inner := c.Sub(caDir.Tangent().MulS(penWidth))
		vsEdgePos = appendVec2(vsEdgePos, c, inner)
	}
	if len(vsEdgePos) > 0 {
		return newShape(newVertexBuffer(
			newVertexStream("aPosition", stFloatVec2, vsEdgePos),
		), nil, dmTriangleStrip)
	}
	return nil
}
