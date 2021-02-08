package main

import (
	"math"
	"math/rand"
)

const infinity float64 = math.MaxFloat64

//DegToRad converts degrees to radians
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

//RandomDouble return a random float64 in [0.0, 1.0)
func RandomDouble(rnd *rand.Rand) float64 {
	return rnd.Float64()
}

//RandomDoubleRange return random vlaue in [min, max)
func RandomDoubleRange(min float64, max float64, rnd *rand.Rand) float64 {
	return min + (max-min)*RandomDouble(rnd)
}

//Clamp clamps x between min and max
func Clamp(x float64, min float64, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func RandomCosineDirection(rnd *rand.Rand) Vec3 {
	r1 := RandomDouble(rnd)
	r2 := RandomDouble(rnd)
	z := math.Sqrt(1 - r2)

	phi := 2 * math.Pi * r1
	x := math.Cos(phi) * math.Sqrt(r2)
	y := math.Sin(phi) * math.Sqrt(r2)

	return Vec3{x, y, z}
}
