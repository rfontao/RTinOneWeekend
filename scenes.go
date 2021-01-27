package main

import (
	"math"
	"math/rand"
)

func threeBallScene() (world hittableList) {
	materialGround := lambertian{color3{0.8, 0.8, 0.0}}
	materialCenter := lambertian{color3{0.1, 0.2, 0.5}}
	// materialLeft := metal{color3{0.8, 0.8, 0.8}, 0.3}
	// materialCenter := dielectric{1.5}
	materialLeft := dielectric{1.5}
	materialRight := metal{color3{0.8, 0.6, 0.2}, 0.0}

	world.Add(&sphere{point3{0, -100.5, -1}, 100.0, materialGround})
	world.Add(&sphere{point3{0, 0, -1}, 0.5, materialCenter})
	world.Add(&sphere{point3{-1, 0, -1}, 0.5, materialLeft})
	world.Add(&sphere{point3{-1, 0, -1}, -0.4, materialLeft}) //TGlass ball
	world.Add(&sphere{point3{1, 0, -1}, 0.5, materialRight})

	return world
}

func testWideViewScene() (world hittableList) {
	//Test of wide view
	R := math.Cos(math.Pi / 4.0)
	materialLeft := lambertian{color3{0, 0, 1}}
	materialRight := lambertian{color3{1, 0, 0}}

	world.Add(&sphere{point3{-R, 0, -1}, R, materialLeft})
	world.Add(&sphere{point3{R, 0, -1}, R, materialRight})

	return world
}

func randomScene() (world hittableList) {

	groundMaterial := lambertian{color3{0.5, 0.5, 0.5}}
	world.Add(&sphere{point3{0, -1000, 0}, 1000, groundMaterial})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()

			center := point3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if center.Sub(point3{4, 0.2, 0}).Length() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					albedo := color3{rand.Float64() * rand.Float64(), rand.Float64() * rand.Float64(), rand.Float64() * rand.Float64()}
					sphereMaterial := lambertian{albedo}
					world.Add(&sphere{center, 0.2, sphereMaterial})
				} else if chooseMat < 0.95 {
					// metal
					albedo := color3{0.5 * (1 + rand.Float64()), 0.5 * (1 + rand.Float64()), 0.5 * (1 + rand.Float64())}
					fuzz := 0.5 * rand.Float64()
					sphereMaterial := metal{albedo, fuzz}
					world.Add(&sphere{center, 0.2, sphereMaterial})
				} else {
					// glass
					sphereMaterial := dielectric{1.5}
					world.Add(&sphere{center, 0.2, sphereMaterial})
				}
			}

		}
	}

	material1 := dielectric{1.5}
	world.Add(&sphere{point3{0, 1, 0}, 1.0, material1})

	material2 := lambertian{color3{0.4, 0.2, 0.1}}
	world.Add(&sphere{point3{-4, 1, 0}, 1.0, material2})

	material3 := metal{color3{0.7, 0.6, 0.5}, 0.0}
	world.Add(&sphere{point3{4, 1, 0}, 1.0, material3})

	return world

}
