package main

import "math/rand"

type ray struct {
	origin    Point3
	direction Vec3
}

func (r *ray) At(t float64) Point3 {
	return r.origin.Add(r.direction.Mult(t))
}

//"Background" color (colors can be changed)
func (r *ray) RayColor(world *hittableList, maxDepth int, rnd *rand.Rand) Color3 {
	if maxDepth <= 0 {
		//No more light gathered
		return Color3{0, 0, 0}
	}

	rec, hit := world.hit(r, 0.001, infinity)
	if hit {

		scattered, attenuation, scatter := rec.mat.scatter(r, rec, rnd)
		// attenuation.Print()
		if scatter {
			return attenuation.MultEach(scattered.RayColor(world, maxDepth-1, rnd))
		}
		return Color3{0, 0, 0}
	}

	t := 0.5 * (r.direction.Normalize().Y() + 1.0)
	return Lerp(Color3{1.0, 1.0, 1.0}, Color3{0.5, 0.7, 1.0}, t)
}
