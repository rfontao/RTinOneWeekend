package main

import (
	"math"
	"math/rand"
)

type camera struct {
	origin          Point3
	lowerLeftCorner Point3
	horizontal      Vec3
	vertical        Vec3
	u, v, w         Vec3
	lensRadius      float64
	time0, time1    float64
}

//vfov in degrees
func initCamera(lookFrom Point3, lookAt Point3, up Vec3, vfov float64, aspectRatio float64, aperture float64, focusDist float64, t0 float64, t1 float64) (c camera) {
	theta := DegToRad(vfov)
	h := math.Tan(theta / 2.0)
	viewportHeight := 2.0 * h
	viewportWidth := aspectRatio * viewportHeight

	c.time0 = t0
	c.time1 = t1

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

func (c camera) getRay(s float64, t float64, rnd *rand.Rand) *ray {

	rd := RandomInUnitDisk(rnd).Mult(c.lensRadius)
	offset := (c.u.Mult(rd.X())).Add(c.v.Mult(rd.Y()))

	return &ray{c.origin.Add(offset), c.lowerLeftCorner.Add(c.horizontal.Mult(s)).Add(c.vertical.Mult(t)).Sub(c.origin).Sub(offset), RandomDoubleRange(c.time0, c.time1, rnd)}
}
