package main

import (
	"fmt"
	"math"
)

type vec3 [3]float64
type color3 = vec3
type point3 = vec3

func (v vec3) Print() {
	fmt.Printf("vec3[% 0.3f, % 0.3f, % 0.3f]", v[0], v[1], v[2])
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

func (v vec3) Mul(c float64) vec3 {
	return vec3{v[0] * c, v[1] * c}
}

func (v vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v vec3) LengthSquared() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}
