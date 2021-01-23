package main

import (
	"image/color"
	"math"
)

func color3ToRGBA(c color3, samplesPerPixel int) color.RGBA {
	const rgbMax float64 = 255.0

	r := c.X()
	g := c.Y()
	b := c.Z()

	scale := 1.0 / float64(samplesPerPixel)
	r = math.Sqrt(scale * r)
	g = math.Sqrt(scale * g)
	b = math.Sqrt(scale * b)

	return color.RGBA{uint8(256.0 * clamp(r, 0.0, 0.999)),
		uint8(256.0 * clamp(g, 0.0, 0.999)),
		uint8(256.0 * clamp(b, 0.0, 0.999)),
		0xff}
}
