package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"sync"
	"time"
)

type options struct {
	aspectRatio     float64
	imageWidth      int
	imageHeight     int
	samplesPerPixel int
	maxDepth        int
}

func main() {

	//Image
	aspectRatio := 16.0 / 9.0
	imageWidth := 1200
	imageHeight := int(float64(imageWidth) / aspectRatio)

	opts := options{
		aspectRatio:     aspectRatio,
		imageWidth:      imageWidth,
		imageHeight:     imageHeight,
		samplesPerPixel: 100,
		maxDepth:        50,
	}

	// World/Camera
	world := randomScene()

	//Camera
	lookFrom := Point3{13, 2, 3}
	lookAt := Point3{0, 0, 0}
	up := Vec3{0, 1, 0}
	distToFocus := 10.0 //lookAt.Sub(lookFrom).Length()
	aperture := 0.1

	c := initCamera(lookFrom, lookAt, up, 20, aspectRatio, aperture, distToFocus, 0.0, 1.0)

	//Render

	t0 := time.Now()

	upLeft := image.Point{0, 0}
	lowRight := image.Point{imageWidth - 1, imageHeight - 1}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	wg := sync.WaitGroup{}
	for y := imageHeight - 1; y >= 0; y-- {
		//Draw each pixel of line
		go func(row int) {
			wg.Add(1)
			for x := 0; x < imageWidth; x++ {
				ch := make(chan Color3, opts.samplesPerPixel)

				pixelColor := Color3{0, 0, 0}
				sendRays(world, &c, x, row, &opts, ch)

				for i := 0; i < opts.samplesPerPixel; i++ {
					pixelColor = pixelColor.Add(<-ch)
				}
				// Colors are defined by Red, Green, Blue, Alpha uint8 values.
				img.Set(x, imageHeight-row, Color3ToRGBA(pixelColor, opts.samplesPerPixel))
			}
			//Maybe change later
			fmt.Printf("%d/%d lines\n", imageHeight-row, imageHeight)
			wg.Done()
		}(y)
	}
	wg.Wait()

	// Encode as PNG.
	f, _ := os.Create("images/a.png")
	png.Encode(f, img)

	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func sendRays(world hittable, c *camera, x int, y int, opts *options, ch chan Color3) {

	for s := 0; s < opts.samplesPerPixel; s++ {
		go func() {
			rnd := rand.New(rand.NewSource(rand.Int63()))

			//Horizontal ratio?
			u := (float64(x) + RandomDouble(rnd)) / float64(opts.imageWidth-1)
			//Vertical ratio?
			v := (float64(y) + RandomDouble(rnd)) / float64(opts.imageHeight-1)

			currentRay := c.getRay(u, v, rnd)
			rayColor := currentRay.RayColor(world, opts.maxDepth, rnd)
			// pixelColor = pixelColor.Add(rayColor)
			ch <- rayColor
		}()
	}
}
