package main

import (
	"math"
	"math/rand"
)

type scatterRecord struct {
	specularRay ray
	isSpecular  bool
	attenuation Color3
	pdf         pdf
}

type material interface {
	scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (sRec *scatterRecord, scatter bool)
	emitted(rayIn *ray, rec *hitRecord, u float64, v float64, p Point3) Color3
	scatteringPdf(rayIn *ray, rec *hitRecord, scattered *ray) float64
}

type lambertian struct {
	albedo texture
}

func (lamb lambertian) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (sRec *scatterRecord, scatter bool) {

	var sRecord scatterRecord
	sRecord.isSpecular = false
	sRecord.attenuation = lamb.albedo.value(rec.u, rec.v, rec.p)
	sRecord.pdf = newCosinePdf(rec.normal)

	return &sRecord, true
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

func (m metal) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (sRec *scatterRecord, scatter bool) {

	reflected := Reflect(rayIn.direction.Normalize(), rec.normal)
	var sRecord scatterRecord
	sRecord.specularRay = ray{rec.p, reflected.Add(RandomInUnitSphere(rnd).Mult(m.fuzz)), 0}
	sRecord.attenuation = m.albedo
	sRecord.isSpecular = true
	sRecord.pdf = nil
	return &sRecord, true
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

func (m dielectric) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (sRec *scatterRecord, scatter bool) {

	var sRecord scatterRecord
	sRecord.isSpecular = true
	sRecord.pdf = nil

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

	sRecord.attenuation = Color3{1, 1, 1}
	sRecord.specularRay = ray{rec.p, direction, rayIn.time}

	return &sRecord, true
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

func (m diffuseLight) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (sRec *scatterRecord, scatter bool) {
	return nil, false
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

func (m isotropic) scatter(rayIn *ray, rec *hitRecord, rnd *rand.Rand) (sRec *scatterRecord, scatter bool) {
	var sRecord scatterRecord
	sRecord.isSpecular = false
	sRecord.pdf = newSpherePdf()
	sRecord.attenuation = m.albedo.value(rec.u, rec.v, rec.p)
	return &sRecord, true
}

func (m isotropic) emitted(rayIn *ray, rec *hitRecord, u float64, v float64, p Point3) Color3 {
	return Color3{0, 0, 0}
}

func (m isotropic) scatteringPdf(rayIn *ray, rec *hitRecord, scattered *ray) float64 {
	return 1.0 / (4.0 * math.Pi)
}
