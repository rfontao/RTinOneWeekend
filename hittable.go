package main

import (
	"math"
)

type hitRecord struct {
	p         point3
	normal    vec3
	t         float64
	frontFace bool
	mat       material
}

type hittable interface {
	hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool)
}

func (rec *hitRecord) setFaceNormal(r *ray, outwardNormal vec3) {
	rec.frontFace = r.direction.Dot(outwardNormal) < 0
	if rec.frontFace == true {
		rec.normal = outwardNormal
	} else {
		rec.normal = outwardNormal.Mult(-1.0)
	}
}

type sphere struct {
	center point3
	radius float64
	mat    material
}

func (s *sphere) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	// const tolerance float64 = 0.01

	rToCenter := r.origin.Sub(s.center) //A - C
	a := r.direction.LengthSquared()    // r dir DOT r dir
	h := rToCenter.Dot(r.direction)
	c := rToCenter.LengthSquared() - math.Pow(s.radius, 2)

	discriminant := math.Pow(h, 2) - a*c

	if discriminant < 0 {
		return nil, false
	}
	discSqrt := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range
	root := (-h - discSqrt) / a

	if root < tMin || root > tMax {
		root = (-h + discSqrt) / a
		if root < tMin || root > tMax {
			return nil, false
		}
	}

	//Used to make the rays not colide with t = 0
	// if math.Abs(root) < tolerance {
	// 	return rec, false
	// }

	hitPoint := r.At(root)
	outwardNormal := hitPoint.Sub(s.center).Div(s.radius)
	rec := hitRecord{
		t:   root,
		p:   hitPoint,
		mat: s.mat,
	}
	rec.setFaceNormal(r, outwardNormal)

	return &rec, true
}

type hittableList struct {
	objects []hittable
}

func (list *hittableList) hit(r *ray, tMin float64, tMax float64) (rec *hitRecord, hit bool) {

	hitAnything := false
	closestSoFar := tMax

	for _, obj := range list.objects {
		if hitRec, hit := obj.hit(r, tMin, closestSoFar); hit {
			hitAnything = true
			closestSoFar = hitRec.t
			rec = hitRec
		}
	}

	return rec, hitAnything
}

//TODO
func (list *hittableList) Add(h hittable) {
	list.objects = append(list.objects, h)
}
