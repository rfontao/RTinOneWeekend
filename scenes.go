package main

import (
	"math"
	"math/rand"
)

func threeBallScene() hittable {

	var world hittableList

	materialGround := lambertian{solidColor{Color3{0.8, 0.8, 0.0}}}
	materialCenter := lambertian{solidColor{Color3{0.1, 0.2, 0.5}}}
	// materialLeft := metal{Color3{0.8, 0.8, 0.8}, 0.3}
	// materialCenter := dielectric{1.5}
	materialLeft := dielectric{1.5}
	materialRight := metal{Color3{0.8, 0.6, 0.2}, 0.0}

	world.Add(&sphere{Point3{0, -100.5, -1}, 100.0, materialGround})
	world.Add(&sphere{Point3{0, 0, -1}, 0.5, materialCenter})
	world.Add(&sphere{Point3{-1, 0, -1}, 0.5, materialLeft})
	world.Add(&sphere{Point3{-1, 0, -1}, -0.4, materialLeft}) //TGlass ball
	world.Add(&sphere{Point3{1, 0, -1}, 0.5, materialRight})

	return newBvhNode(world.objects, 0.0, 1.0)
}

func testWideViewScene() hittable {

	var world hittableList
	//Test of wide view
	R := math.Cos(math.Pi / 4.0)
	materialLeft := lambertian{solidColor{Color3{0, 0, 1}}}
	materialRight := lambertian{solidColor{Color3{1, 0, 0}}}

	world.Add(&sphere{Point3{-R, 0, -1}, R, materialLeft})
	world.Add(&sphere{Point3{R, 0, -1}, R, materialRight})

	return newBvhNode(world.objects, 0.0, 1.0)
}

func randomScene() hittable {
	var world hittableList

	groundMaterial := lambertian{checkerTexture{solidColor{Color3{0.2, 0.3, 0.1}}, solidColor{Color3{0.9, 0.9, 0.9}}}}
	world.Add(&sphere{Point3{0, -1000, 0}, 1000, groundMaterial})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()

			center := Point3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if center.Sub(Point3{4, 0.2, 0}).Length() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					albedo := Color3{rand.Float64() * rand.Float64(), rand.Float64() * rand.Float64(), rand.Float64() * rand.Float64()}
					sphereMaterial := lambertian{solidColor{albedo}}
					world.Add(&sphere{center, 0.2, sphereMaterial})
				} else if chooseMat < 0.95 {
					// metal
					albedo := Color3{0.5 * (1 + rand.Float64()), 0.5 * (1 + rand.Float64()), 0.5 * (1 + rand.Float64())}
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
	world.Add(&sphere{Point3{0, 1, 0}, 1.0, material1})

	material2 := lambertian{solidColor{Color3{0.4, 0.2, 0.1}}}
	world.Add(&sphere{Point3{-4, 1, 0}, 1.0, material2})

	material3 := metal{Color3{0.7, 0.6, 0.5}, 0.0}
	world.Add(&sphere{Point3{4, 1, 0}, 1.0, material3})

	return newBvhNode(world.objects, 0.0, 1.0)

}

func randomSceneMoving() hittable {

	var world hittableList

	groundMaterial := lambertian{solidColor{Color3{0.5, 0.5, 0.5}}}
	world.Add(&sphere{Point3{0, -1000, 0}, 1000, groundMaterial})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()

			center := Point3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if center.Sub(Point3{4, 0.2, 0}).Length() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					albedo := Color3{rand.Float64() * rand.Float64(), rand.Float64() * rand.Float64(), rand.Float64() * rand.Float64()}
					center2 := center.Add(Vec3{0, 0.5 * rand.Float64(), 0})
					sphereMaterial := lambertian{solidColor{albedo}}
					world.Add(&movingSphere{center, center2, 0.0, 1.0, 0.2, sphereMaterial})
				} else if chooseMat < 0.95 {
					// metal
					albedo := Color3{0.5 * (1 + rand.Float64()), 0.5 * (1 + rand.Float64()), 0.5 * (1 + rand.Float64())}
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
	world.Add(&sphere{Point3{0, 1, 0}, 1.0, material1})

	material2 := lambertian{solidColor{Color3{0.4, 0.2, 0.1}}}
	world.Add(&sphere{Point3{-4, 1, 0}, 1.0, material2})

	material3 := metal{Color3{0.7, 0.6, 0.5}, 0.0}
	world.Add(&sphere{Point3{4, 1, 0}, 1.0, material3})

	return newBvhNode(world.objects, 0.0, 1.0)

}

func twoSpheres() hittable {
	var world hittableList

	checker := lambertian{checkerTexture{solidColor{Color3{0.2, 0.3, 0.1}}, solidColor{Color3{0.9, 0.9, 0.9}}}}

	world.Add(&sphere{Point3{0, -10, 0}, 10, checker})
	world.Add(&sphere{Point3{0, 10, 0}, 10, checker})

	return newBvhNode(world.objects, 0.0, 1.0)
}

func twoPerlinSpheres() hittable {
	var world hittableList

	noise := lambertian{noiseTexture{newPerlin(), 4}}

	world.Add(&sphere{Point3{0, -1000, 0}, 1000, noise})
	world.Add(&sphere{Point3{0, 2, 0}, 2, noise})

	return newBvhNode(world.objects, 0.0, 1.0)
}

func imageTextureTest() hittable {
	var world hittableList

	// imTex := lambertian{newImageTexture("unknown.png")}
	imTex := lambertian{newImageTexture("earthmap.jpg")}

	world.Add(&sphere{Point3{0, 0, 0}, 2, imTex})

	return newBvhNode(world.objects, 0.0, 1.0)
}

func simpleLight() hittable {
	var world hittableList

	noise := lambertian{noiseTexture{newPerlin(), 4}}

	world.Add(&sphere{Point3{0, -1000, 0}, 1000, noise})
	world.Add(&sphere{Point3{0, 2, 0}, 2, noise})

	diffLight := diffuseLight{solidColor{Color3{4, 4, 4}}}
	world.Add(&xyRect{diffLight, 3, 5, 1, 3, -2})

	return newBvhNode(world.objects, 0.0, 1.0)
}

func cornellBox() hittable {
	var world hittableList

	red := lambertian{solidColor{Color3{0.65, 0.05, 0.05}}}
	white := lambertian{solidColor{Color3{0.73, 0.73, 0.73}}}
	green := lambertian{solidColor{Color3{0.12, 0.45, 0.15}}}
	light := diffuseLight{solidColor{Color3{15, 15, 15}}}

	world.Add(&yzRect{green, 0, 555, 0, 555, 555})
	world.Add(&yzRect{red, 0, 555, 0, 555, 0})
	world.Add(&flipFace{&xzRect{light, 213, 343, 227, 332, 554}})
	world.Add(&xzRect{white, 0, 555, 0, 555, 0})
	world.Add(&xzRect{white, 0, 555, 0, 555, 555})
	world.Add(&xyRect{white, 0, 555, 0, 555, 555})

	var box1 hittable = newBox(Point3{0, 0, 0}, Point3{165, 330, 165}, white)
	box1 = newRotateY(box1, 15)
	box1 = &translate{box1, Vec3{265, 0, 295}}
	world.Add(box1)

	var box2 hittable = newBox(Point3{0, 0, 0}, Point3{165, 165, 165}, white)
	box2 = newRotateY(box2, -18)
	box2 = &translate{box2, Vec3{130, 0, 65}}
	world.Add(box2)

	return newBvhNode(world.objects, 0.0, 1.0)
}

func cornellSmoke() hittable {
	var world hittableList

	red := lambertian{solidColor{Color3{0.65, 0.05, 0.05}}}
	white := lambertian{solidColor{Color3{0.73, 0.73, 0.73}}}
	green := lambertian{solidColor{Color3{0.12, 0.45, 0.15}}}
	light := diffuseLight{solidColor{Color3{15, 15, 15}}}

	world.Add(&yzRect{green, 0, 555, 0, 555, 555})
	world.Add(&yzRect{red, 0, 555, 0, 555, 0})
	world.Add(&xzRect{light, 113, 443, 127, 432, 554})
	world.Add(&xzRect{white, 0, 555, 0, 555, 0})
	world.Add(&xzRect{white, 0, 555, 0, 555, 555})
	world.Add(&xyRect{white, 0, 555, 0, 555, 555})

	var box1 hittable = newBox(Point3{0, 0, 0}, Point3{165, 330, 165}, white)
	box1 = newRotateY(box1, 15)
	box1 = &translate{box1, Vec3{265, 0, 295}}

	var box2 hittable = newBox(Point3{0, 0, 0}, Point3{165, 165, 165}, white)
	box2 = newRotateY(box2, -18)
	box2 = &translate{box2, Vec3{130, 0, 65}}

	world.Add(newConstantMedium(box1, 0.01, solidColor{Color3{0, 0, 0}}))
	world.Add(newConstantMedium(box2, 0.01, solidColor{Color3{1, 1, 1}}))

	return newBvhNode(world.objects, 0.0, 1.0)
}

func finalScene() hittable {
	var boxes1 hittableList
	ground := lambertian{solidColor{Color3{0.48, 0.83, 0.53}}}

	boxesPerSide := 20
	for i := 0; i < boxesPerSide; i++ {
		for j := 0; j < boxesPerSide; j++ {
			w := 100.0
			x0 := -1000.0 + float64(i)*w
			z0 := -1000.0 + float64(j)*w
			y0 := 0.0
			x1 := x0 + w
			y1 := 1 + 100.0*rand.Float64() //1--100
			z1 := z0 + w

			boxes1.Add(newBox(Point3{x0, y0, z0}, Point3{x1, y1, z1}, ground))
		}
	}

	var objects hittableList
	objects.Add(newBvhNode(boxes1.objects, 0, 1))

	light := diffuseLight{solidColor{Color3{7, 7, 7}}}
	objects.Add(&xzRect{light, 123, 423, 147, 412, 554})

	center1 := Point3{400, 400, 200}
	center2 := center1.Add(Vec3{30, 0, 0})
	movingSphereMaterial := lambertian{solidColor{Color3{0.7, 0.3, 0.1}}}
	objects.Add(&movingSphere{center1, center2, 0, 1, 50, movingSphereMaterial})

	objects.Add(&sphere{Point3{260, 150, 45}, 50, dielectric{1.5}})
	objects.Add(&sphere{Point3{0, 150, 145}, 50, metal{Color3{0.8, 0.8, 0.9}, 1.0}})

	boundary := sphere{Point3{360, 150, 145}, 70, dielectric{1.5}}
	objects.Add(&boundary)
	objects.Add(newConstantMedium(&boundary, 0.2, solidColor{Color3{0.2, 0.4, 0.9}}))
	boundary = sphere{Point3{0, 0, 0}, 5000, dielectric{1.5}}
	objects.Add(newConstantMedium(&boundary, 1000, solidColor{Color3{1, 1, 1}}))

	emat := lambertian{newImageTexture("earthmap.jpg")}
	objects.Add(&sphere{Point3{400, 200, 400}, 100, emat})
	pertext := noiseTexture{newPerlin(), 0.1}
	objects.Add(&sphere{Point3{220, 280, 300}, 80, lambertian{pertext}})

	var boxes2 hittableList
	white := lambertian{solidColor{Color3{0.73, 0.73, 0.73}}}
	ns := 1000
	for j := 0; j < ns; j++ {
		boxes2.Add(&sphere{Point3{rand.Float64() * 165, rand.Float64() * 165, rand.Float64() * 165}, 10, white})
	}

	objects.Add(&translate{
		newRotateY(
			newBvhNode(boxes2.objects, 0, 1), 15),
		Vec3{-100, 270, 395},
	})

	return newBvhNode(objects.objects, 0, 1)
}
