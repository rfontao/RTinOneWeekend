package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {

	//Image
	const imageHeight int = 256
	const imageWidth int = 256

	//Render

	upLeft := image.Point{0, 0}
	lowRight := image.Point{imageWidth - 1, imageHeight - 1}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	var rgbMax = 255.0
	// Set color for each pixel.
	for y := 0; y < imageHeight; y++ {
		fmt.Printf("%d/%d lines\n", y, imageHeight-1)
		for x := 0; x < imageWidth; x++ {

			var r = uint8((float64(x) / float64(imageHeight)) * rgbMax)
			var g = uint8((float64(imageHeight-y) / float64(imageWidth)) * rgbMax)
			var b = uint8(0.25 * rgbMax)

			// Colors are defined by Red, Green, Blue, Alpha uint8 values.
			img.Set(x, y, color.RGBA{r, g, b, 0xff})
			// img.Set(x, y, color.RGBA{uint8(y), 255, 255, 0xff})
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)

	var a vec3 = vec3{1, 2, 3}

	a.Print()
}
