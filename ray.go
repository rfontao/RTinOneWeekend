package main

import "math/rand"

type ray struct {
	origin    Point3
	direction Vec3
	time      float64
}

func (r *ray) At(t float64) Point3 {
	return r.origin.Add(r.direction.Mult(t))
}

//"Background" color (colors can be changed)
func (r *ray) RayColor(world hittable, background Color3, maxDepth int, rnd *rand.Rand) Color3 {
	if maxDepth <= 0 {
		//No more light gathered
		return Color3{0, 0, 0}
	}

	rec, hit := world.hit(r, 0.001, infinity)
	if !hit {
		return background
	}

	emitted := rec.mat.emitted(rec.u, rec.v, rec.p)
	scattered, attenuation, scatter := rec.mat.scatter(r, rec, rnd)
	// attenuation.Print()
	if !scatter {
		return emitted
	}
	return emitted.Add(attenuation.MultEach(scattered.RayColor(world, background, maxDepth-1, rnd)))
}
