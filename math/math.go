// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"math"
)

func R2D(r float32) float32 {
	return 180.0 * r / Pi
}

func D2R(r float32) float32 {
	return Pi * r / 180.0
}

func Absf(v float32) float32 {
	if v < 0 {
		return -v
	} else {
		return v
	}
}

func Round(v float32) int {
	if v < 0 {
		return int(v - 0.5)
	} else {
		return int(v + 0.4999999)
	}
}

func Sinf(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func Cosf(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

func Tanf(v float32) float32 {
	return float32(math.Tan(float64(v)))
}

func Asinf(v float32) float32 {
	return float32(math.Asin(float64(v)))
}

func Acosf(v float32) float32 {
	return float32(math.Acos(float64(v)))
}

func Atanf(v float32) float32 {
	return float32(math.Atan(float64(v)))
}

func Sqrtf(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}

func Powf(v, e float32) float32 {
	return float32(math.Pow(float64(v), float64(e)))
}

func Lerp(a, b int, s float32) int {
	r := float32(b - a)
	return a + int(r*s)
}

func Lerpf(a, b float32, s float32) float32 {
	r := b - a
	return a + r*s
}

func Ramp(s float32, a, b float32) float32 {
	return (s - a) / (b - a)
}

func RampSat(s float32, a, b float32) float32 {
	return Saturate((s - a) / (b - a))
}

func Saturate(x float32) float32 {
	return Clampf(x, 0, 1)
}

func SmoothStep(s float32, a, b float32) float32 {
	x := RampSat(s, a, b)
	return x * x * (3 - 2*x)
}

func Clamp(x, min, max int) int {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}

func Clampf(x, min, max float32) float32 {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}

func Min(values ...int) int {
	m := MaxInt
	for _, v := range values {
		if v < m {
			m = v
		}
	}
	return m
}

func Minf(values ...float32) float32 {
	m := float32(math.MaxFloat32)
	for _, v := range values {
		if v < m {
			m = v
		}
	}
	return m
}

func Max(values ...int) int {
	m := MinInt
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}

func Maxf(values ...float32) float32 {
	m := float32(-math.MaxFloat32)
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}

func Mod(a, b int) int {
	x := a % b
	if x < 0 {
		return x + b
	} else {
		return x
	}
}
