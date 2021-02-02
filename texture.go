package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

type texture interface {
	value(u float64, v float64, p Vec3) Color3
}

type solidColor struct {
	colorValue Color3
}

func (s solidColor) value(u float64, v float64, p Vec3) Color3 {
	return s.colorValue
}

type checkerTexture struct {
	odd  texture
	even texture
}

func (s checkerTexture) value(u float64, v float64, p Vec3) Color3 {

	sines := math.Sin(10.0*p.X()) * math.Sin(10.0*p.Y()) * math.Sin(10.0*p.Z())
	if sines < 0 {
		return s.odd.value(u, v, p)
	}
	return s.even.value(u, v, p)
}

type noiseTexture struct {
	noise perlin
	scale float64
}

func (s noiseTexture) value(u float64, v float64, p Vec3) Color3 {
	return Color3{1, 1, 1}.Mult(0.5).Mult(1 + math.Sin(p.Z()*s.scale+10*s.noise.turb(p, 7)))
}

type imageTexture struct {
	im            image.Image
	width, height int
	ok            bool
}

func newImageTexture(filename string) imageTexture {
	reader, err := os.Open("textures/" + filename)
	if err != nil {
		fmt.Print("Could not open file: " + filename)
		return imageTexture{nil, 0, 0, false}
	}
	defer reader.Close()

	im, _, err := image.Decode(reader)
	if err != nil {
		fmt.Print("Error decoding image: " + filename)
		return imageTexture{nil, 0, 0, false}
	}
	b := im.Bounds()

	return imageTexture{im, b.Max.X, b.Max.Y, true}
}

func (s imageTexture) value(u float64, v float64, p Vec3) Color3 {
	if !s.ok {
		return Vec3{0.5, 0.5, 0.5}
	}

	// Clamp input texture coordinates to [0,1] x [1,0]
	u = Clamp(u, 0.0, 1.0)
	v = 1.0 - Clamp(v, 0.0, 1.0) // Flip V to image coordinates

	i := int(u * float64(s.width))
	j := int(v * float64(s.height))

	if i >= s.width {
		i = s.width - 1
	}
	if j >= s.height {
		j = s.height - 1
	}
	b := s.im.Bounds()

	// fmt.Printf("i: %d j: %d u: %f v: %f\n", i, j, u, v)
	// r, g, b1, _ := s.im.At(b.Min.X+i, b.Min.Y+j).RGBA()
	// fmt.Printf("r: %d g: %d b: %d\n", r>>8, g>>8, b1>>8)

	return RGBAToColor3(s.im.At(b.Min.X+i, b.Min.Y+j))

}
