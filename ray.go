package main

type ray struct {
	origin    point3
	direction vec3
}

func (r ray) At(t float64) point3 {
	return r.origin.Add(r.direction.Mult(t))
}

//"Background" color (colors can be changed)
func (r ray) RayColor(world hittable) color3 {
	rec, hit := world.hit(r, 0, infinity)
	if hit {
		// return rec.normal.Add(vec3{1, 1, 1}).Mult(0.5)
		// return color3{rec.normal.X() + 1, rec.normal.Y() + 1, rec.normal.Z() + 1}.Mult(0.5)
		return color3{1 - rec.normal.X(), 1 - rec.normal.Y(), 1 - rec.normal.Z()}.Mult(0.5)
	}
	rec.t = 0.5 * (r.direction.Normalize().Y() + 1.0)
	return Lerp(color3{1.0, 1.0, 1.0}, color3{0.5, 0.7, 1.0}, rec.t)
}
