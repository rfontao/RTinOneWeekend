package main

import "math"

type camera struct {
	origin          point3
	lowerLeftCorner point3
	horizontal      vec3
	vertical        vec3
	u, v, w         vec3
	lensRadius      float64
}

//vfov in degrees
func initCamera(lookFrom point3, lookAt point3, up vec3, vfov float64, aspectRatio float64, aperture float64, focusDist float64) (c camera) {
	theta := degToRad(vfov)
	h := math.Tan(theta / 2.0)
	viewportHeight := 2.0 * h
	viewportWidth := aspectRatio * viewportHeight

	c.w = (lookFrom.Sub(lookAt)).Normalize()
	c.u = (up.Cross(c.w)).Normalize()
	c.v = c.w.Cross(c.u)

	c.origin = lookFrom
	c.horizontal = c.u.Mult(viewportWidth * focusDist)
	c.vertical = c.v.Mult(viewportHeight * focusDist)
	c.lowerLeftCorner = c.origin.Sub(c.horizontal.Div(2)).Sub(c.vertical.Div(2)).Sub(c.w.Mult(focusDist))

	c.lensRadius = aperture / 2.0

	return c
}

func (c camera) getRay(s float64, t float64) *ray {

	rd := randomInUnitDisk().Mult(c.lensRadius)
	offset := (c.u.Mult(rd.X())).Add(c.v.Mult(rd.Y()))

	return &ray{c.origin.Add(offset), c.lowerLeftCorner.Add(c.horizontal.Mult(s)).Add(c.vertical.Mult(t)).Sub(c.origin).Sub(offset)}
}
