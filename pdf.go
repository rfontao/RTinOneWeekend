package main

import (
	"math"
	"math/rand"
)

type pdf interface {
	value(direction Vec3) float64
	generate(rnd *rand.Rand) Vec3
}

type cosinePdf struct {
	uvw onb
}

func newCosinePdf(w Vec3) cosinePdf {
	return cosinePdf{buildFromW(w)}
}

func (pdf cosinePdf) value(direction Vec3) float64 {
	cosine := direction.Normalize().Dot(pdf.uvw.w())
	if cosine <= 0 {
		return 0
	}
	return cosine / math.Pi
}

func (pdf cosinePdf) generate(rnd *rand.Rand) Vec3 {
	return pdf.uvw.local(RandomCosineDirection(rnd))
}

type hittablePdf struct {
	obj hittable
	o   Point3
}

func (pdf hittablePdf) value(direction Vec3) float64 {
	return pdf.obj.pdfValue(pdf.o, direction)
}

func (pdf hittablePdf) generate(rnd *rand.Rand) Vec3 {
	return pdf.obj.random(pdf.o, rnd)
}

type mixturePdf struct {
	p [2]pdf
}

func (pdf mixturePdf) value(direction Vec3) float64 {
	return 0.5*pdf.p[0].value(direction) + 0.5*pdf.p[1].value(direction)
}

func (pdf mixturePdf) generate(rnd *rand.Rand) Vec3 {
	if RandomDouble(rnd) < 0.5 {
		return pdf.p[0].generate(rnd)
	}
	return pdf.p[1].generate(rnd)
}
