package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
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
		fmt.Printf("%d/%d lines\n", imageHeight-y, imageHeight-1)
		for x := 0; x < imageWidth; x++ {

			pixelColor := color3{0, 0, 0}
			for s := 0; s < samplesPerPixel; s++ {
				//Horizontal ratio?
				u := (float64(x) + randomDouble()) / float64(imageWidth-1)
				//Vertical ratio?
				v := (float64(y) + randomDouble()) / float64(imageHeight-1)
				currentRay := c.getRay(u, v)
				rayColor := currentRay.RayColor(world, maxDepth)
				pixelColor = pixelColor.Add(rayColor)
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

/*
HitSphere checks if a ray r hits a sphere with center center and radius r
If so return t

t2b⋅b+2tb⋅(A−C)+(A−C)⋅(A−C)−r2=0 => a = b.b; b =
A = ray origin
b = ray direction

*/
func HitSphere(center point3, radius float64, r ray) float64 {
	rToCenter := r.origin.Sub(center) //A - C
	a := r.direction.LengthSquared()  // r dir DOT r dir
	h := rToCenter.Dot(r.direction)
	c := rToCenter.LengthSquared() - math.Pow(radius, 2)

	discriminant := math.Pow(h, 2) - a*c

	if discriminant < 0 {
		return -1.0
	}
	return (-h - math.Sqrt(discriminant)/a)
}

// func HitSphere(center point3, radius float64, r ray) float64 {
// 	rToCenter := r.origin.Sub(center) //A - C
// 	a := r.direction.Dot(r.direction)
// 	b := rToCenter.Dot(r.direction) * 2.0
// 	c := rToCenter.Dot(rToCenter) - math.Pow(radius, 2)

// 	discriminat := math.Pow(b, 2) - 4.0*a*c

// 	if discriminat < 0 {
// 		return -1.0
// 	}
// 	return (-b - math.Sqrt(discriminat)/(2.0*a))
// }
