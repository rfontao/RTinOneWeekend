package main

import (
	"math/rand"
)

type ray struct {
	origin    Point3
	direction Vec3
	time      float64
}

func (r *ray) At(t float64) Point3 {
	return r.origin.Add(r.direction.Mult(t))
}

//"Background" color (colors can be changed)
// func (r *ray) RayColor(world hittable, background Color3, maxDepth int, rnd *rand.Rand) Color3 {
// 	if maxDepth <= 0 {
// 		//No more light gathered
// 		return Color3{0, 0, 0}
// 	}

// 	rec, hit := world.hit(r, 0.001, infinity)
// 	if !hit {
// 		return background
// 	}

// 	emitted := rec.mat.emitted(r, rec, rec.u, rec.v, rec.p)
// 	scattered, albedo, pdf, scatter := rec.mat.scatter(r, rec, rnd)
// 	// attenuation.Print()
// 	if !scatter {
// 		return emitted
// 	}
// 	return emitted.Add(albedo.MultEach(scattered.RayColor(world, background, maxDepth-1, rnd).Mult(rec.mat.scatteringPdf(r, rec, scattered)).Div(pdf)))
// }

func (r *ray) RayColor(world hittable, background Color3, maxDepth int, rnd *rand.Rand, lights hittable) Color3 {
	if maxDepth <= 0 {
		//No more light gathered
		return Color3{0, 0, 0}
	}

	rec, hit := world.hit(r, 0.001, infinity)
	if !hit {
		return background
	}

	emitted := rec.mat.emitted(r, rec, rec.u, rec.v, rec.p)
	scattered, albedo, pdfVal, scatter := rec.mat.scatter(r, rec, rnd)
	// attenuation.Print()
	if !scatter {
		return emitted
	}

	// onLight := Point3{RandomDoubleRange(213, 343, rnd), 554, RandomDoubleRange(227, 332, rnd)}
	// toLight := onLight.Sub(rec.p)
	// distanceSquared := toLight.LengthSquared()
	// toLight = toLight.Normalize()

	// if toLight.Dot(rec.normal) < 0 {
	// 	return emitted
	// }

	// lightArea := (343 - 213) * (332 - 227)
	// lightCosine := math.Abs(toLight.Y())
	// if lightCosine < 0.0000001 {
	// 	return emitted
	// }

	// pdf = distanceSquared / (lightCosine * float64(lightArea))
	// scattered = &ray{rec.p, toLight, r.time}

	p0 := hittablePdf{lights, rec.p}
	p1 := newCosinePdf(rec.normal)
	mixedPdf := mixturePdf{[2]pdf{p0, p1}}

	scattered = &ray{rec.p, mixedPdf.generate(rnd), r.time}
	pdfVal = mixedPdf.value(scattered.direction)

	return emitted.Add(albedo.MultEach(scattered.RayColor(world, background, maxDepth-1, rnd, lights).Mult(rec.mat.scatteringPdf(r, rec, scattered)).Div(pdfVal)))
}
