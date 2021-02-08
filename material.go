package main

import (
	"math"
	"math/rand"
)

type material interface {
	scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, pdf float64, scatter bool)
	emitted(rayIn *ray, rec *hitRecord, u float64, v float64, p Point3) Color3
	scatteringPdf(rayIn *ray, rec *hitRecord, scattered *ray) float64
}

type lambertian struct {
	albedo texture
}

func (lamb lambertian) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, pdf float64, scatter bool) {

	uvw := buildFromW(rec.normal)
	direction := uvw.local(RandomCosineDirection(rnd))
	scattered = &ray{rec.p, direction.Normalize(), rayIn.time}
	at := lamb.albedo.value(rec.u, rec.v, rec.p)
	pdf = uvw.w().Dot(scattered.direction) / math.Pi
	return scattered, &at, pdf, true
}

func (lamb lambertian) emitted(rayIn *ray, rec *hitRecord, u float64, v float64, p Point3) Color3 {
	return Color3{0, 0, 0}
}

func (lamb lambertian) scatteringPdf(rayIn *ray, rec *hitRecord, scattered *ray) float64 {
	cosine := rec.normal.Dot(scattered.direction.Normalize())
	if cosine < 0 {
		return 0
	}
	return cosine / math.Pi
}

type metal struct {
	albedo Color3
	fuzz   float64 //Radius of sphere
}

func (m metal) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, pdf float64, scatter bool) {

	Reflected := Reflect(rayIn.direction.Normalize(), rec.normal)

	scattered = &ray{rec.p, Reflected.Add(RandomInUnitSphere(rnd).Mult(m.fuzz)), rayIn.time}
	attenuation = &m.albedo
	return scattered, attenuation, 0, scattered.direction.Dot(rec.normal) > 0
}

func (m metal) emitted(rayIn *ray, rec *hitRecord, u float64, v float64, p Point3) Color3 {
	return Color3{0, 0, 0}
}

func (m metal) scatteringPdf(rayIn *ray, rec *hitRecord, scattered *ray) float64 {
	return 0
}

type dielectric struct {
	ir float64 //Index of Refraction
}

func (m dielectric) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, pdf float64, scatter bool) {

	var RefractionRatio float64
	if rec.frontFace {
		RefractionRatio = 1.0 / m.ir
	} else {
		RefractionRatio = m.ir
	}

	unitDirection := rayIn.direction.Normalize()
	cosTheta := math.Min(unitDirection.Mult(-1).Dot(rec.normal), 1.0)
	sinTheta := math.Sqrt(1.0 - math.Pow(cosTheta, 2))

	cannotRefract := RefractionRatio*sinTheta > 1.0

	var direction Vec3
	if cannotRefract || reflectance(cosTheta, RefractionRatio) > RandomDouble(rnd) {
		direction = Reflect(unitDirection, rec.normal)
	} else {
		direction = Refract(unitDirection, rec.normal, RefractionRatio)

	}

	attenuation = &Color3{1, 1, 1}
	scattered = &ray{rec.p, direction, rayIn.time}

	return scattered, attenuation, 0, true
}

func (m dielectric) emitted(rayIn *ray, rec *hitRecord, u float64, v float64, p Point3) Color3 {
	return Color3{0, 0, 0}
}

func (m dielectric) scatteringPdf(rayIn *ray, rec *hitRecord, scattered *ray) float64 {
	return 0
}

func reflectance(cosine float64, refIndex float64) float64 {
	// Use Schlick's approximation for Reflectance.
	r0 := math.Pow((1.0-refIndex)/(1.0+refIndex), 2)
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}

type diffuseLight struct {
	emit texture
}

func (m diffuseLight) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, pdf float64, scatter bool) {
	return nil, nil, 0, false
}

func (m diffuseLight) emitted(rayIn *ray, rec *hitRecord, u float64, v float64, p Point3) Color3 {
	if rec.frontFace {
		return m.emit.value(u, v, p)
	}
	return m.emit.value(u, v, p)
}

func (m diffuseLight) scatteringPdf(rayIn *ray, rec *hitRecord, scattered *ray) float64 {
	return 0
}

type isotropic struct {
	albedo texture
}

func (m isotropic) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, pdf float64, scatter bool) {
	s := ray{rec.p, RandomInUnitSphere(rnd), rayIn.time}
	a := m.albedo.value(rec.u, rec.v, rec.p)
	return &s, &a, 0, true
}

func (m isotropic) emitted(rayIn *ray, rec *hitRecord, u float64, v float64, p Point3) Color3 {
	return Color3{0, 0, 0}
}

func (m isotropic) scatteringPdf(rayIn *ray, rec *hitRecord, scattered *ray) float64 {
	return 0
}
