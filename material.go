package main

import (
	"math"
	"math/rand"
)

type material interface {
	scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, scatter bool)
}

type lambertian struct {
	albedo Color3
}

func (lamb lambertian) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, scatter bool) {

	scatterDirection := rec.normal.Add(RandomUnitVector(rnd))
	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.normal
	}

	scattered = &ray{rec.p, scatterDirection}
	attenuation = &lamb.albedo
	return scattered, attenuation, true
}

type metal struct {
	albedo Color3
	fuzz   float64 //Radius of sphere
}

func (m metal) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, scatter bool) {

	Reflected := Reflect(rayIn.direction.Normalize(), rec.normal)

	scattered = &ray{rec.p, Reflected.Add(RandomInUnitSphere(rnd).Mult(m.fuzz))}
	attenuation = &m.albedo
	return scattered, attenuation, scattered.direction.Dot(rec.normal) > 0
}

type dielectric struct {
	ir float64 //Index of Refraction
}

func (m dielectric) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (scattered *ray, attenuation *Color3, scatter bool) {

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
	scattered = &ray{rec.p, direction}

	return scattered, attenuation, true
}

func reflectance(cosine float64, refIndex float64) float64 {
	// Use Schlick's approximation for Reflectance.
	r0 := math.Pow((1.0-refIndex)/(1.0+refIndex), 2)
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}
