package main

import "math"

type onb [3]Vec3

func (b onb) u() Vec3 {
	return b[0]
}

func (b onb) v() Vec3 {
	return b[1]
}

func (b onb) w() Vec3 {
	return b[2]
}

func (b onb) local(a Vec3) Vec3 {
	return b.u().Mult(a.X()).Add(b.v().Mult(a.Y())).Add(b.w().Mult(a.Z()))
}

func buildFromW(n Vec3) onb {
	var b onb

	b[2] = n.Normalize()
	var a Vec3
	if math.Abs(b.w().X()) > 0.9 {
		a = Vec3{0, 1, 0}
	} else {
		a = Vec3{1, 0, 0}
	}

	b[1] = (b.w().Cross(a)).Normalize()
	b[0] = b.w().Cross(b.v())

	return b
}
