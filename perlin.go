package main

import (
	"math"
	"math/rand"
)

type perlin struct {
	pointCount int
	ranVec     []Vec3
	permX      []int
	permY      []int
	permZ      []int
}

func newPerlin() perlin {
	var p perlin
	p.pointCount = 256

	p.ranVec = make([]Vec3, 256)
	for i := 0; i < p.pointCount; i++ {
		p.ranVec[i] = Vec3{rand.Float64()*2 - 1, rand.Float64()*2 - 1, rand.Float64()*2 - 1}.Normalize()
	}

	p.permX = p.perlinGeneratePerm()
	p.permY = p.perlinGeneratePerm()
	p.permZ = p.perlinGeneratePerm()

	return p
}

func (p perlin) noise(point Point3) float64 {

	u := point.X() - math.Floor(point.X())
	v := point.Y() - math.Floor(point.Y())
	w := point.Z() - math.Floor(point.Z())
	u = u * u * (3 - 2*u)
	v = v * v * (3 - 2*v)
	w = w * w * (3 - 2*w)

	i := int(math.Floor(point.X()))
	j := int(math.Floor(point.Y()))
	k := int(math.Floor(point.Z()))

	var c [2][2][2]Vec3
	// c := [2][2][2]Vec3{
	// 	{{Vec3{0.0, 0.0}, 0.0}, {0.0, 0.0}},
	// 	{{0.0, 0.0}, {0.0, 0.0}},
	// }

	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = p.ranVec[p.permX[(i+di)&255]^p.permY[(j+dj)&255]^p.permZ[(k+dk)&255]]
			}
		}
	}

	return perlinInterp(c, u, v, w)
}

func (p perlin) perlinGeneratePerm() []int {
	points := make([]int, p.pointCount)

	for i := 0; i < p.pointCount; i++ {
		points[i] = i
	}

	permute(points, p.pointCount)

	return points
}

func (p perlin) turb(point Point3, depth int) float64 {
	accum := 0.0
	tempPoint := point
	weight := 1.0

	for i := 0; i < depth; i++ {
		accum += weight * p.noise(tempPoint)
		weight *= 0.5
		tempPoint = tempPoint.Mult(2.0)
	}

	return math.Abs(accum)
}

func permute(p []int, n int) {
	for i := n - 1; i > 0; i-- {
		target := rand.Intn(i + 1)
		p[i], p[target] = p[target], p[i]
	}
}

func perlinInterp(c [2][2][2]Vec3, u, v, w float64) float64 {

	uu := u * u * (3 - 2*u)
	vv := v * v * (3 - 2*v)
	ww := w * w * (3 - 2*w)
	accum := 0.0

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weightV := Vec3{u - float64(i), v - float64(j), w - float64(k)}
				accum += (float64(i)*uu + (1-float64(i))*(1-uu)) *
					(float64(j)*vv + (1-float64(j))*(1-vv)) *
					(float64(k)*ww + (1-float64(k))*(1-ww)) * c[i][j][k].Dot(weightV)
			}
		}
	}

	return accum
}
