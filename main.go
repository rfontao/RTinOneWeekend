package main

import (
	"image"
	"image/png"
	"os"
)

func main() {

	//Image
	const aspectRatio = 16.0 / 9.0
	const imageWidth int = 400
	const imageHeight int = int(float64(imageWidth) / aspectRatio)

	//Camera

	var viewportHeight = 2.0
	var viewportWidth = aspectRatio * viewportHeight
	var focalLength = 1.0

	var origin = point3{0, 0, 0}
	var horizontal = vec3{viewportWidth, 0, 0}
	var vertical = vec3{0, viewportHeight, 0}
	var lowerLeftCorner = origin.Sub(horizontal.Div(2)).Sub(vertical.Div(2)).Sub(vec3{0, 0, focalLength})
	lowerLeftCorner.Print()

	//Render

	upLeft := image.Point{0, 0}
	lowRight := image.Point{imageWidth - 1, imageHeight - 1}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for y := 0; y < imageHeight; y++ {
		// fmt.Printf("%d/%d lines\n", y, imageHeight-1)
		for x := 0; x < imageWidth; x++ {

			//Horizontal ratio?
			var u float64 = float64(x) / float64(imageWidth-1)
			//Vertical ratio?
			var v float64 = float64(y) / float64(imageHeight-1)

			var currentRay = ray{origin, lowerLeftCorner.Add(horizontal.Mult(u)).Add(vertical.Mult(v)).Sub(origin)}
			// Colors are defined by Red, Green, Blue, Alpha uint8 values.
			img.Set(x, imageHeight-y, color3ToRGBA(currentRay.RayColor()))
		}
	}

	// Encode as PNG.
	f, _ := os.Create("a.png")
	png.Encode(f, img)

	// var a ray = ray{point3{0, 0, 0}, vec3{1, 0, 3}}

	// // fmt.Print(a.Dot(b))
	// a.At(2).Print()
}
