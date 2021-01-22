package main

type camera struct {
	origin          point3
	lowerLeftCorner point3
	horizontal      vec3
	vertical        vec3
}

func initCamera() (c camera) {
	const aspectRatio = 16.0 / 9.0
	viewportHeight := 2.0
	viewportWidth := aspectRatio * viewportHeight
	focalLength := 1.0

	c.origin = point3{0, 0, 0}
	c.horizontal = vec3{viewportWidth, 0, 0}
	c.vertical = vec3{0, viewportHeight, 0}
	c.lowerLeftCorner = c.origin.Sub(c.horizontal.Div(2)).Sub(c.vertical.Div(2)).Sub(vec3{0, 0, focalLength})
	return c
}

func (c camera) getRay(u float64, v float64) ray {
	return ray{c.origin, c.lowerLeftCorner.Add(c.horizontal.Mult(u)).Add(c.vertical.Mult(v)).Sub(c.origin)}
}
