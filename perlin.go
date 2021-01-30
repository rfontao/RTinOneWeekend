package main

import (
	"math"
	"math/rand"
)

type perlin struct {
	pointCount int
	ranFloat   []float64
	permX      []int
	permY      []int
	permZ      []int
}

func newPerlin() perlin {
	var p perlin
	p.pointCount = 256

	p.ranFloat = make([]float64, 256)
	for i := 0; i < p.pointCount; i++ {
		p.ranFloat[i] = rand.Float64()
	}

	p.permX = p.perlinGeneratePerm()
	p.permY = p.perlinGeneratePerm()
	p.permZ = p.perlinGeneratePerm()

	return p
}

func (p perlin) noise(point Point3) float64 {
	// i := int(4.0*point.X()) & 255
	// j := int(4.0*point.Y()) & 255
	// k := int(4.0*point.Z()) & 255

	// return p.ranFloat[p.permX[i]^p.permY[j]^p.permZ[k]]

	u := point.X() - math.Floor(point.X())
	v := point.Y() - math.Floor(point.Y())
	w := point.Z() - math.Floor(point.Z())
	u = u * u * (3 - 2*u)
	v = v * v * (3 - 2*v)
	w = w * w * (3 - 2*w)

	i := int(math.Floor(point.X()))
	j := int(math.Floor(point.Y()))
	k := int(math.Floor(point.Z()))

	c := [2][2][2]float64{
		{{0.0, 0.0}, {0.0, 0.0}},
		{{0.0, 0.0}, {0.0, 0.0}},
	}

	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = p.ranFloat[p.permX[(i+di)&255]^p.permY[(j+dj)&255]^p.permZ[(k+dk)&255]]
			}
		}
	}

	return trilinearInterp(c, u, v, w)
}

func (p perlin) perlinGeneratePerm() []int {
	points := make([]int, p.pointCount)

	for i := 0; i < p.pointCount; i++ {
		points[i] = i
	}

	permute(points, p.pointCount)

	return points
}

func permute(p []int, n int) {
	for i := n - 1; i > 0; i-- {
		target := rand.Intn(i + 1)
		p[i], p[target] = p[target], p[i]
	}
}

func trilinearInterp(c [2][2][2]float64, u, v, w float64) float64 {
	accum := 0.0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				accum += (float64(i)*u + (1-float64(i))*(1-u)) *
					(float64(j)*v + (1-float64(j))*(1-v)) *
					(float64(k)*w + (1-float64(k))*(1-w)) * c[i][j][k]
			}
		}
	}

	return accum
}
