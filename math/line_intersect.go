// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

func linesIntersect(p0, p1, q0, q1 Vec2) bool {
	const epsilon = 0.0001
	// http://stackoverflow.com/a/565282
	// t = (q − p) × s / (r × s)
	r, s := p1.Sub(p0), q1.Sub(q0)
	PmQ := p0.Sub(q0)
	QmP := q0.Sub(p0)
	ScR := s.Cross(r)
	if Absf(ScR) < 0.0001 {
		// Lines are parallel
		QmPcR := QmP.Cross(r)
		if Absf(QmPcR) < 0.0001 {
			// Colinear
			QmPdR := QmP.Dot(r)
			PmQdS := PmQ.Dot(s)
			if (0 <= QmPdR && QmPdR < r.Dot(r)) ||
				(0 <= PmQdS && PmQdS < s.Dot(s)) {
				return true
			}
		}
		return false
	}
	t := PmQ.Cross(s) / ScR
	u := PmQ.Cross(r) / ScR
	return t > epsilon && u > epsilon && t < 1-epsilon && u < 1-epsilon
}
