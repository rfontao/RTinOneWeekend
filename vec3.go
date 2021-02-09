package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
)

//Vec3 -> X Y Z
type Vec3 [3]float64

//Color3 -> R G B
type Color3 = Vec3

//Point3 -> Same as vec3
type Point3 = Vec3

//Print the contents of the Vec3
func (v Vec3) Print() {
	fmt.Printf("Vec3[% 0.3f, % 0.3f, % 0.3f]\n", v[0], v[1], v[2])
}

//Access functions

//X -> index 0
func (v Vec3) X() float64 {
	return v[0]
}

//Y -> index 1
func (v Vec3) Y() float64 {
	return v[1]
}

//Z -> index 2
func (v Vec3) Z() float64 {
	return v[2]
}

//Add return the addition of v with v2
func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v[0] + v2[0], v[1] + v2[1], v[2] + v2[2]}
}

//Sub returns v - v2
func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v[0] - v2[0], v[1] - v2[1], v[2] - v2[2]}
}

//Mult return the scalar multiplication v * c
func (v Vec3) Mult(c float64) Vec3 {
	return Vec3{v[0] * c, v[1] * c, v[2] * c}
}

//MultEach is for each member v * v2
func (v Vec3) MultEach(v2 Vec3) Vec3 {
	return Vec3{v[0] * v2[0], v[1] * v2[1], v[2] * v2[2]}
}

//Div is v / c
func (v Vec3) Div(c float64) Vec3 {
	return v.Mult(1.0 / c)
}

//Length of v
func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

//LengthSquared of v
func (v Vec3) LengthSquared() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

//Dot product -> v.v2
func (v Vec3) Dot(v2 Vec3) float64 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

//Normalize v
func (v Vec3) Normalize() Vec3 {
	return v.Div(v.Length())
}

//Cross is v x v2
func (v Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{v[1]*v2[2] - v[2]*v2[1], v[2]*v2[0] - v[0]*v2[2], v[0]*v2[1] - v[1]*v2[0]}
}

func (v Vec3) Copy() Vec3 {
	return Vec3{v[0], v[1], v[2]}
}

// Lerp (1−t)⋅startValue+t⋅endValue
func Lerp(start Vec3, end Vec3, t float64) Vec3 {
	return start.Mult(1.0 - t).Add(end.Mult(t))
}

//RandomVec3 gives a vec3 with random values in  [0.0, 1.0)
func RandomVec3(rnd *rand.Rand) Vec3 {
	return Vec3{RandomDouble(rnd), RandomDouble(rnd), RandomDouble(rnd)}
}

//RandomRangeVec3 is RandomVec3 but with range [min, max)
func RandomRangeVec3(min float64, max float64, rnd *rand.Rand) Vec3 {
	return Vec3{RandomDoubleRange(min, max, rnd), RandomDoubleRange(min, max, rnd), RandomDoubleRange(min, max, rnd)}
}

//RandomInUnitSphere is a random vec3 with length inferior to 1
func RandomInUnitSphere(rnd *rand.Rand) Vec3 {
	for {
		p := RandomRangeVec3(-1, 1, rnd)
		if p.LengthSquared() >= 1 {
			continue
		}
		return p
	}
}

//RandomUnitVector is RandomInUnitSphere but normalized
func RandomUnitVector(rnd *rand.Rand) Vec3 {
	return RandomInUnitSphere(rnd).Normalize()
}

//RandomInHemisphere not sure
func RandomInHemisphere(normal Vec3, rnd *rand.Rand) Vec3 {
	inUnitSphere := RandomInUnitSphere(rnd)

	if inUnitSphere.Dot(normal) > 0.0 {
		return inUnitSphere
	}
	return inUnitSphere.Mult(-1)

}

//RandomInUnitDisk is a random vec3 in a circunference with z = 0 and length < 1
func RandomInUnitDisk(rnd *rand.Rand) Vec3 {
	for {
		p := Vec3{RandomDoubleRange(-1, 1, rnd), RandomDoubleRange(-1, 1, rnd), 0}
		if p.LengthSquared() >= 1 {
			continue
		}
		return p
	}
}

//NearZero checks if all values of a vec3 are close to zero
func (v Vec3) NearZero() bool {
	const s = 1e-18
	return (math.Abs(v[0]) < s) && (math.Abs(v[1]) < s) && (math.Abs(v[2]) < s)
}

//Reflect v with normal n
func Reflect(v Vec3, n Vec3) Vec3 {
	return v.Sub(n.Mult(v.Dot(n) * 2.0))
}

//Refract refracts uv with normal n
func Refract(uv Vec3, n Vec3, etaRatio float64) Vec3 {
	cosTheta := math.Min(uv.Mult(-1).Dot(n), 1.0)
	rOutPerp := (uv.Add(n.Mult(cosTheta))).Mult(etaRatio)
	rOutParallel := n.Mult(-math.Sqrt(math.Abs(1 - rOutPerp.LengthSquared())))

	return rOutParallel.Add(rOutPerp)
}

//Color3ToRGBA return the RGBA equivalent of a Color3
func Color3ToRGBA(c Color3, samplesPerPixel int) color.RGBA {
	const rgbMax float64 = 255.0

	r := c.X()
	g := c.Y()
	b := c.Z()

	if r != r {
		r = 0.0
	}
	if g != g {
		g = 0.0
	}
	if b != b {
		b = 0.0
	}

	scale := 1.0 / float64(samplesPerPixel)
	r = math.Sqrt(scale * r)
	g = math.Sqrt(scale * g)
	b = math.Sqrt(scale * b)

	return color.RGBA{uint8(256.0 * Clamp(r, 0.0, 0.999)),
		uint8(256.0 * Clamp(g, 0.0, 0.999)),
		uint8(256.0 * Clamp(b, 0.0, 0.999)),
		0xff}
}

//RGBAToColor3 .
func RGBAToColor3(c color.Color) Color3 {

	const rgbMax float64 = 255.0

	r, g, b, _ := c.RGBA()

	return Color3{float64(r>>8) / rgbMax, float64(g>>8) / rgbMax, float64(b>>8) / rgbMax}

	// r := c.X()
	// g := c.Y()
	// b := c.Z()

	// scale := 1.0 / float64(samplesPerPixel)
	// r = math.Sqrt(scale * r)
	// g = math.Sqrt(scale * g)
	// b = math.Sqrt(scale * b)

	// return color.RGBA{uint8(256.0 * Clamp(r, 0.0, 0.999)),
	// 	uint8(256.0 * Clamp(g, 0.0, 0.999)),
	// 	uint8(256.0 * Clamp(b, 0.0, 0.999)),
	// 	0xff}
}
