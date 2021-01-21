package main

type ray struct {
	origin    point3
	direction vec3
}

func (r ray) At(t float64) point3 {
	return r.origin.Add(r.direction.Mult(t))
}

//"Background" color (colors can be changed)
func (r ray) RayColor() color3 {
	var t float64 = 0.5 * (r.direction.Normalize().Y() + 1.0)
	return Lerp(color3{1.0, 1.0, 1.0}, color3{0.5, 0.7, 1.0}, t)
}
