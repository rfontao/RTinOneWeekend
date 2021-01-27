package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {

	//Image
	const aspectRatio = 16.0 / 9.0
	const imageWidth int = 400
	const imageHeight int = int(float64(imageWidth) / aspectRatio)
	const samplesPerPixel int = 100
	const maxDepth int = 50

	// World/Camera
	world := threeBallScene()

	//Camera
	lookFrom := point3{-2, 2, 1}
	lookAt := point3{0, 0, -1}
	up := vec3{0, 1, 0}
	distToFocus := 5.0
	aperture := 0.1

	c := initCamera(lookFrom, lookAt, up, 20, aspectRatio, aperture, distToFocus)

	//Render

	t0 := time.Now()

	upLeft := image.Point{0, 0}
	lowRight := image.Point{imageWidth - 1, imageHeight - 1}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for y := imageHeight - 1; y >= 0; y-- {
		fmt.Printf("%d/%d lines\n", imageHeight-y, imageHeight)
		for x := 0; x < imageWidth; x++ {
			ch := make(chan color3, samplesPerPixel)

			pixelColor := color3{0, 0, 0}
			sendRays(&world, &c, x, y, imageWidth, imageHeight, samplesPerPixel, maxDepth, ch)

			for i := 0; i < samplesPerPixel; i++ {
				pixelColor = pixelColor.Add(<-ch)
			}
			// Colors are defined by Red, Green, Blue, Alpha uint8 values.
			img.Set(x, imageHeight-y, color3ToRGBA(pixelColor, samplesPerPixel))
		}
	}

	// Encode as PNG.
	f, _ := os.Create("a.png")
	png.Encode(f, img)

	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func sendRays(world *hittableList, c *camera, x int, y int, imageWidth int, imageHeight int, samplesPerPixel int, maxDepth int, ch chan color3) {

	for s := 0; s < samplesPerPixel; s++ {
		go func() {
			rnd := rand.New(rand.NewSource(rand.Int63()))

			//Horizontal ratio?
			u := (float64(x) + randomDouble(rnd)) / float64(imageWidth-1)
			//Vertical ratio?
			v := (float64(y) + randomDouble(rnd)) / float64(imageHeight-1)

			currentRay := c.getRay(u, v, rnd)
			rayColor := currentRay.RayColor(world, maxDepth, rnd)
			// pixelColor = pixelColor.Add(rayColor)
			ch <- rayColor
		}()
	}
}
