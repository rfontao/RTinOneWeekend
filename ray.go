package main

type ray struct {
	origin    point3
	direction vec3
}

func (r ray) At(t float64) point3 {
	return r.origin.Add(r.direction.Mult(t))
}

//"Background" color (colors can be changed)
func (r ray) RayColor(world hittable, maxDepth int) color3 {

	if maxDepth <= 0 {
		//No more light gathered
		return color3{0, 0, 0}
	}

	rec, hit := world.hit(r, 0.001, infinity)
	if hit == true {
		// target := rec.p.Add(rec.normal).Add(randomInUnitSphere())
		target := rec.p.Add(randomInHemisphere(rec.normal))

		res := ray{rec.p, target.Sub(rec.p)}.RayColor(world, maxDepth-1).Mult(0.5)
		// res.Print()
		return res
	}
	rec.t = 0.5 * (r.direction.Normalize().Y() + 1.0)
	return Lerp(color3{1.0, 1.0, 1.0}, color3{0.5, 0.7, 1.0}, rec.t)
}
