package main

import (
	"fmt"
	"math"
	"math/rand"
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

func (v vec3) MultEach(v2 vec3) vec3 {
	return vec3{v[0] * v2[0], v[1] * v2[1], v[2] * v2[2]}
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

func randomVec3(rnd *rand.Rand) vec3 {
	return vec3{randomDouble(rnd), randomDouble(rnd), randomDouble(rnd)}
}

func randomRangeVec3(min float64, max float64, rnd *rand.Rand) vec3 {
	return vec3{randomDoubleRange(min, max, rnd), randomDoubleRange(min, max, rnd), randomDoubleRange(min, max, rnd)}
}

func randomInUnitSphere(rnd *rand.Rand) vec3 {
	for {
		p := randomRangeVec3(-1, 1, rnd)
		if p.LengthSquared() >= 1 {
			continue
		}
		return p
	}
}

//True lambertian reflection
func randomUnitVector(rnd *rand.Rand) vec3 {
	return randomInUnitSphere(rnd).Normalize()
}

func randomInHemisphere(normal vec3, rnd *rand.Rand) vec3 {
	inUnitSphere := randomInUnitSphere(rnd)

	if inUnitSphere.Dot(normal) > 0.0 {
		return inUnitSphere
	}
	return inUnitSphere.Mult(-1)

}

func randomInUnitDisk(rnd *rand.Rand) vec3 {
	for {
		p := vec3{randomDoubleRange(-1, 1, rnd), randomDoubleRange(-1, 1, rnd), 0}
		if p.LengthSquared() >= 1 {
			continue
		}
		return p
	}
}

func (v vec3) nearZero() bool {
	const s = 1e-18
	return (math.Abs(v[0]) < s) && (math.Abs(v[1]) < s) && (math.Abs(v[2]) < s)
}

func reflect(v vec3, n vec3) vec3 {
	return v.Sub(n.Mult(v.Dot(n) * 2.0))
}

func refract(uv vec3, n vec3, etaRatio float64) vec3 {
	cosTheta := math.Min(uv.Mult(-1).Dot(n), 1.0)
	rOutPerp := (uv.Add(n.Mult(cosTheta))).Mult(etaRatio)
	rOutParallel := n.Mult(-math.Sqrt(math.Abs(1 - rOutPerp.LengthSquared())))

	return rOutParallel.Add(rOutPerp)
}
