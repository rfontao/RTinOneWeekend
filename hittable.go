package main

import (
	"math"
	"math/rand"
	"sort"
)

type hitRecord struct {
	p         Point3
	normal    Vec3
	t         float64
	frontFace bool
	mat       material
}

type hittable interface {
	hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool)
	boundingBox(time0 float64, time1 float64) (aabb, bool)
}

func (rec *hitRecord) setFaceNormal(r *ray, outwardNormal Vec3) {
	rec.frontFace = r.direction.Dot(outwardNormal) < 0
	if rec.frontFace == true {
		rec.normal = outwardNormal
	} else {
		rec.normal = outwardNormal.Mult(-1.0)
	}
}

type sphere struct {
	center Point3
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

func (s *sphere) boundingBox(time0 float64, time1 float64) (aabb, bool) {
	return aabb{
		s.center.Sub(Vec3{s.radius, s.radius, s.radius}),
		s.center.Add(Vec3{s.radius, s.radius, s.radius}),
	}, true
}

type movingSphere struct {
	center0, center1 Point3
	time0, time1     float64
	radius           float64
	mat              material
}

func (s *movingSphere) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	// const tolerance float64 = 0.01

	rToCenter := r.origin.Sub(s.center(r.time)) //A - C
	a := r.direction.LengthSquared()            // r dir DOT r dir
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
	outwardNormal := hitPoint.Sub(s.center(r.time)).Div(s.radius)
	rec := hitRecord{
		t:   root,
		p:   hitPoint,
		mat: s.mat,
	}
	rec.setFaceNormal(r, outwardNormal)

	return &rec, true
}

func (s *movingSphere) center(time float64) Point3 {
	return Lerp(s.center0, s.center1, time)
	// return s.center0 + ((time-s.time0)/(s.time1-s.time0))*(s.center1.Sub(s.center0))
}

func (s *movingSphere) boundingBox(time0 float64, time1 float64) (aabb, bool) {
	box0 := aabb{
		s.center(time0).Sub(Vec3{s.radius, s.radius, s.radius}),
		s.center(time0).Add(Vec3{s.radius, s.radius, s.radius}),
	}

	box1 := aabb{
		s.center(time1).Sub(Vec3{s.radius, s.radius, s.radius}),
		s.center(time1).Add(Vec3{s.radius, s.radius, s.radius}),
	}

	return surroundingBox(box0, box1), true
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

func (list *hittableList) Add(h hittable) {
	list.objects = append(list.objects, h)
}

func (list *hittableList) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {

	if len(list.objects) == 0 {
		return outputBox, false
	}

	firstBox := true

	for _, object := range list.objects {
		tempBox, exists := object.boundingBox(time0, time1)
		if !exists {
			return outputBox, false
		}

		if firstBox {
			outputBox = tempBox
		} else {
			outputBox = surroundingBox(outputBox, tempBox)
		}
		firstBox = false
	}

	return outputBox, true
}

type By func(h1 hittable, h2 hittable) bool

func (by By) Sort(l []hittable) {
	hs := &hittableSorter{
		hittables: l,
		by:        by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(hs)
}

type hittableSorter struct {
	hittables []hittable
	by        func(h1 hittable, h2 hittable) bool
}

func (s *hittableSorter) Len() int { return len(s.hittables) }
func (s *hittableSorter) Swap(i, j int) {
	s.hittables[i], s.hittables[j] = s.hittables[j], s.hittables[i]
}
func (s *hittableSorter) Less(i, j int) bool {
	return s.by(s.hittables[i], s.hittables[j])
}

type bvhNode struct {
	left, right hittable
	box         aabb
}

func newBvhNode(list []hittable, start int, end int, time0 float64, time1 float64) *bvhNode {
	objs := list

	var bvh bvhNode

	axis := rand.Intn(2 + 1)

	comparator := []By{
		//X axis
		func(h1 hittable, h2 hittable) bool {
			return boxCompare(h1, h2, 0)
		},
		//Y axis
		func(h1 hittable, h2 hittable) bool {
			return boxCompare(h1, h2, 1)
		},
		//Z axis
		func(h1 hittable, h2 hittable) bool {
			return boxCompare(h1, h2, 2)
		},
	}[axis]

	objectSpan := end - start

	if objectSpan == 1 {
		bvh.right = objs[start]
		bvh.left = objs[start]
	} else if objectSpan == 2 {
		if comparator(objs[start], objs[start+1]) {
			bvh.left = objs[start]
			bvh.right = objs[start+1]
		} else {
			bvh.left = objs[start+1]
			bvh.right = objs[start]
		}
	} else {
		//only part of list
		comparator.Sort(objs)

		mid := start + objectSpan/2
		bvh.left = newBvhNode(objs, start, mid, time0, time1)
		bvh.right = newBvhNode(objs, mid, end, time0, time1)
	}

	boxLeft, existsLeft := bvh.left.boundingBox(time0, time1)
	boxRight, existsRight := bvh.right.boundingBox(time0, time1)

	if !existsLeft || !existsRight {
		panic("No bounding box in bvhnode constructor")
	}

	bvh.box = surroundingBox(boxLeft, boxRight)
	return &bvh
}

func (bvh *bvhNode) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {
	if !bvh.box.hit(r, tMin, tMax) {
		return nil, false
	}

	var result *hitRecord
	rec, hitLeft := bvh.left.hit(r, tMin, tMax)
	if hitLeft {
		result = rec
		tMax = rec.t
	}

	rec, hitRight := bvh.right.hit(r, tMin, tMax)
	if hitRight {
		result = rec
	}

	return result, hitLeft || hitRight

}

func (bvh bvhNode) boundingBox(time0 float64, time1 float64) (aabb, bool) {
	return bvh.box, true
}

func boxCompare(a hittable, b hittable, axis int) bool {

	boxA, existsA := a.boundingBox(0, 0)
	boxB, existsB := b.boundingBox(0, 0)

	if !existsA || !existsB {
		panic("No bounding box in bvh_node constructor\n")
	}

	return boxA.minimum[axis] < boxB.minimum[axis]
}
