package main

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
}

func (m metal) scatter(rayIn ray, rec hitRecord) (scattered ray, attenuation color3, scatter bool) {

	reflected := reflect(rayIn.direction.Normalize(), rec.normal)

	scattered = ray{rec.p, reflected}
	attenuation = m.albedo
	return scattered, attenuation, scattered.direction.Dot(rec.normal) > 0
}
