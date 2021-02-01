package main

import "math"

type texture interface {
	value(u float64, v float64, p Vec3) Color3
}

type solidColor struct {
	colorValue Color3
}

func (s solidColor) value(u float64, v float64, p Vec3) Color3 {
	return s.colorValue
}

type checkerTexture struct {
	odd  texture
	even texture
}

func (s checkerTexture) value(u float64, v float64, p Vec3) Color3 {

	sines := math.Sin(10.0*p.X()) * math.Sin(10.0*p.Y()) * math.Sin(10.0*p.Z())
	if sines < 0 {
		return s.odd.value(u, v, p)
	}
	return s.even.value(u, v, p)
}

type noiseTexture struct {
	noise perlin
	scale float64
}

func (s noiseTexture) value(u float64, v float64, p Vec3) Color3 {
	return Color3{1, 1, 1}.Mult(0.5).Mult(1 + math.Sin(p.Z()*s.scale+10*s.noise.turb(p, 7)))
}
