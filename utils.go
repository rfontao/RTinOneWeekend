package main

import (
	"math"
	"math/rand"
)

const infinity float64 = math.MaxFloat64

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func randomDouble() float64 {
	return rand.Float64()
}

func randomDoubleRange(min float64, max float64) float64 {
	return min + (max-min)*randomDouble()
}

func clamp(x float64, min float64, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
