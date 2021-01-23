package main

import (
	"fmt"
	"math"
)

type vec3 [3]float64
type color3 = vec3
type point3 = vec3

func (v vec3) Print() {
	fmt.Printf("vec3[% 0.3f, % 0.3f, % 0.3f]\n", v[0], v[1], v[2])
}

//Access functions

func (v vec3) X() float64 {
	return v[0]
}

func (v vec3) Y() float64 {
	return v[1]
}

func (v vec3) Z() float64 {
	return v[2]
}

func (v vec3) Add(v2 vec3) vec3 {
	return vec3{v[0] + v2[0], v[1] + v2[1], v[2] + v2[2]}
}

func (v vec3) Sub(v2 vec3) vec3 {
	return vec3{v[0] - v2[0], v[1] - v2[1], v[2] - v2[2]}
}

func (v vec3) Mult(c float64) vec3 {
	return vec3{v[0] * c, v[1] * c, v[2] * c}
}

func (v vec3) Div(c float64) vec3 {
	return v.Mult(1.0 / c)
}

func (v vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v vec3) LengthSquared() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v vec3) Dot(v2 vec3) float64 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

func (v vec3) Normalize() vec3 {
	return v.Div(v.Length())
}

func (v vec3) Cross(v2 vec3) vec3 {
	return vec3{v[1]*v2[2] - v[2]*v2[1], v[2]*v2[0] - v[0]*v2[2], v[0]*v2[1] - v[1]*v2[0]}
}

// Lerp (1−t)⋅startValue+t⋅endValue
func Lerp(start vec3, end vec3, t float64) vec3 {
	return start.Mult(1.0 - t).Add(end.Mult(t))
}

func randomVec3() vec3 {
	return vec3{randomDouble(), randomDouble(), randomDouble()}
}

func randomRangeVec3(min float64, max float64) vec3 {
	return vec3{randomDoubleRange(min, max), randomDoubleRange(min, max), randomDoubleRange(min, max)}
}

func randomInUnitSphere() vec3 {
	for {
		p := randomRangeVec3(-1, 1)
		if p.LengthSquared() >= 1 {
			continue
		}
		return p
	}
}
