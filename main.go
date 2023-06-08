package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type options struct {
	aspectRatio     float64
	imageWidth      int
	imageHeight     int
	samplesPerPixel int
	maxDepth        int
	background      Color3
}

func main() {

	//Image
	aspectRatio := 16.0 / 9.0
	imageWidth := 800
	imageHeight := int(float64(imageWidth) / aspectRatio)

	opts := options{
		aspectRatio:     aspectRatio,
		imageWidth:      imageWidth,
		imageHeight:     imageHeight,
		samplesPerPixel: 500,
		maxDepth:        5,
		background:      Color3{0, 0, 0},
	}

	var world hittable

	var lookAt, lookFrom Point3
	vfov := 40.0
	aperture := 0.0

	var lights hittableList
	lights.Add(&xzRect{lambertian{}, 213, 343, 227, 332, 554})
	lights.Add(&sphere{Point3{190, 90, 190}, 90, metal{}})

	switch 6 {
	case 1:
		world = randomScene()
		opts.background = Color3{0.7, 0.8, 1.00}

		//Camera
		lookFrom = Point3{13, 2, 3}
		lookAt = Point3{0, 0, 0}
		aperture = 0.1
		vfov = 20
	case 2:
		world = twoSpheres()
		opts.background = Color3{0.7, 0.8, 1.00}

		//Camera
		lookFrom = Point3{13, 2, 3}
		lookAt = Point3{0, 0, 0}
		aperture = 0.1
		vfov = 20
	case 3:
		world = twoPerlinSpheres()
		opts.background = Color3{0.7, 0.8, 1.00}

		//Camera
		lookFrom = Point3{13, 2, 3}
		lookAt = Point3{0, 0, 0}
		vfov = 20
	case 4:
		world = imageTextureTest()
		opts.background = Color3{0.7, 0.8, 1.00}

		//Camera
		lookFrom = Point3{13, 2, 3}
		lookAt = Point3{0, 0, 0}
		vfov = 20
	case 5:
		world = simpleLight()
		opts.samplesPerPixel = 50
		opts.background = Color3{0, 0, 0}

		//Camera
		lookFrom = Point3{26, 3, 6}
		lookAt = Point3{0, 2, 0}
		vfov = 20
	case 6:
		//next week chapter 6.11 -> 2min
		world = cornellBox()
		opts.aspectRatio = 1.0
		opts.imageWidth = 600
		opts.imageHeight = int(float64(opts.imageWidth) / opts.aspectRatio)
		// opts.samplesPerPixel = 2500
		opts.samplesPerPixel = 2000
		// opts.maxDepth = 50
		opts.maxDepth = 5
		opts.background = Color3{0, 0, 0}

		//Camera
		lookFrom = Point3{278, 278, -800}
		lookAt = Point3{278, 278, 0}
		vfov = 40

	case 7:
		world = cornellSmoke()
		opts.aspectRatio = 1.0
		opts.imageWidth = 600
		opts.imageHeight = int(float64(opts.imageWidth) / opts.aspectRatio)
		opts.samplesPerPixel = 50
		opts.background = Color3{0, 0, 0}

		//Camera
		lookFrom = Point3{278, 278, -800}
		lookAt = Point3{278, 278, 0}
		vfov = 40

	case 8:
		world = finalScene()
		opts.aspectRatio = 1.0
		opts.imageWidth = 800
		opts.imageHeight = int(float64(opts.imageWidth) / opts.aspectRatio)
		opts.samplesPerPixel = 100
		opts.maxDepth = 3
		opts.background = Color3{0, 0, 0}

		//Camera
		lookFrom = Point3{478, 278, -600}
		lookAt = Point3{278, 278, 0}
		vfov = 40
	}

	// World/Camera

	up := Vec3{0, 1, 0}
	distToFocus := 10.0 //lookAt.Sub(lookFrom).Length()
	c := initCamera(lookFrom, lookAt, up, vfov, opts.aspectRatio, aperture, distToFocus, 0.0, 1.0)

	//Render

	t0 := time.Now()

	upLeft := image.Point{0, 0}
	lowRight := image.Point{opts.imageWidth - 1, opts.imageHeight - 1}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	wg := sync.WaitGroup{}
	for y := opts.imageHeight - 1; y >= 0; y-- {
		//Draw each pixel of line
		wg.Add(1)
		go func(row int) {
			var bar *progressbar.ProgressBar
			if row == 0 {
				bar = progressbar.Default(int64(opts.imageWidth))
			}
			for x := 0; x < opts.imageWidth; x++ {
				ch := make(chan Color3, opts.samplesPerPixel)

				pixelColor := Color3{0, 0, 0}
				sendRays(world, &c, x, row, &opts, ch, &lights)

				for i := 0; i < opts.samplesPerPixel; i++ {
					pixelColor = pixelColor.Add(<-ch)
				}
				// Colors are defined by Red, Green, Blue, Alpha uint8 values.
				img.Set(x, opts.imageHeight-row, Color3ToRGBA(pixelColor, opts.samplesPerPixel))
				if row == 0 {
					bar.Add(1)
				}
			}
			//Maybe change later
			fmt.Printf("%d/%d lines\n", opts.imageHeight-row, opts.imageHeight)
			wg.Done()
		}(y)
	}
	wg.Wait()

	// Encode as PNG.
	f, _ := os.Create("images/out.png")
	png.Encode(f, img)

	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func sendRays(world hittable, c *camera, x int, y int, opts *options, ch chan Color3, lights hittable) {
	go func() {
		for s := 0; s < opts.samplesPerPixel; s++ {
			// go func() {
			rnd := rand.New(rand.NewSource(rand.Int63()))

			//Horizontal ratio?
			u := (float64(x) + RandomDouble(rnd)) / float64(opts.imageWidth-1)
			//Vertical ratio?
			v := (float64(y) + RandomDouble(rnd)) / float64(opts.imageHeight-1)

			currentRay := c.getRay(u, v, rnd)
			rayColor := currentRay.RayColor(world, opts.background, opts.maxDepth, rnd, lights)
			// pixelColor = pixelColor.Add(rayColor)
			ch <- rayColor
			// }()
		}
	}()
}
