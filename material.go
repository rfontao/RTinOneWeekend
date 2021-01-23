package main

import "math"

type material interface {
	scatter(rayIn ray, rec hitRecord) (scattered ray, attenuation color3, scatter bool)
}

type lambertian struct {
	albedo color3
}

func (lamb lambertian) scatter(rayIn ray, rec hitRecord) (scattered ray, attenuation color3, scatter bool) {

	scatterDirection := rec.normal.Add(randomUnitVector())
	// Catch degenerate scatter direction
	if scatterDirection.nearZero() {
		scatterDirection = rec.normal
	}

	scattered = ray{rec.p, scatterDirection}
	attenuation = lamb.albedo
	return scattered, attenuation, true
}

type metal struct {
	albedo color3
	fuzz   float64 //Radius of sphere
}

func (m metal) scatter(rayIn ray, rec hitRecord) (scattered ray, attenuation color3, scatter bool) {

	reflected := reflect(rayIn.direction.Normalize(), rec.normal)

	scattered = ray{rec.p, reflected.Add(randomInUnitSphere().Mult(m.fuzz))}
	attenuation = m.albedo
	return scattered, attenuation, scattered.direction.Dot(rec.normal) > 0
}

type dielectric struct {
	ir float64 //Index of refraction
}

func (m dielectric) scatter(rayIn ray, rec hitRecord) (scattered ray, attenuation color3, scatter bool) {

	attenuation = color3{1, 1, 1}
	var refractionRatio float64
	if rec.frontFace {
		refractionRatio = 1.0 / m.ir
	} else {
		refractionRatio = m.ir
	}

	unitDirection := rayIn.direction.Normalize()
	cosTheta := math.Min(unitDirection.Mult(-1).Dot(rec.normal), 1.0)
	sinTheta := math.Sqrt(1.0 - math.Pow(cosTheta, 2))

	cannotRefract := refractionRatio*sinTheta > 1.0

	var direction vec3
	if cannotRefract || reflectance(cosTheta, refractionRatio) > randomDouble() {
		direction = reflect(unitDirection, rec.normal)
	} else {
		direction = refract(unitDirection, rec.normal, refractionRatio)

	}

	scattered = ray{rec.p, direction}

	return scattered, attenuation, true
}

func reflectance(cosine float64, refIndex float64) float64 {
	// Use Schlick's approximation for reflectance.
	r0 := math.Pow((1.0-refIndex)/(1.0+refIndex), 2)
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}
