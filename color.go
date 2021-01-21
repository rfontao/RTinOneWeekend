package main

import (
	"image/color"
)

func color3ToRGBA(c color3) color.RGBA {
	const rgbMax float64 = 255.0

	return color.RGBA{uint8(c.X() * rgbMax), uint8(c.Y() * rgbMax), uint8(c.Z() * rgbMax), 0xff}
}
