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
	u         float64
	v         float64
}

type hittable interface {
	hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool)
	boundingBox(time0 float64, time1 float64) (aabb, bool)
	pdfValue(o Point3, v Vec3) float64
	random(o Vec3, rnd *rand.Rand) Vec3
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
	u, v := s.getSphereUV(outwardNormal)
	rec := hitRecord{
		t:   root,
		p:   hitPoint,
		mat: s.mat,
		u:   u,
		v:   v,
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

func (s *sphere) pdfValue(o Point3, v Vec3) float64 {
	_, hit := s.hit(&ray{o, v, 0}, 0.001, infinity)
	if !hit {
		return 0
	}

	cosThetaMax := math.Sqrt(1.0 - s.radius*s.radius/(s.center.Sub(o).LengthSquared()))
	solidAngle := 2 * math.Pi * (1.0 - cosThetaMax)

	return 1.0 / solidAngle
}

func (s *sphere) random(o Vec3, rnd *rand.Rand) Vec3 {
	direction := s.center.Sub(o)
	distanceSquared := direction.LengthSquared()
	uvw := buildFromW(direction)
	return uvw.local(RandomToSphere(s.radius, distanceSquared, rnd))
}

func (s *sphere) getSphereUV(p Point3) (u float64, v float64) {
	// p: a given point on the sphere of radius one, centered at the origin.
	// u: returned value [0,1] of angle around the Y axis from X=-1.
	// v: returned value [0,1] of angle from Y=-1 to Y=+1.
	//     <1 0 0> yields <0.50 0.50>       <-1  0  0> yields <0.00 0.50>
	//     <0 1 0> yields <0.50 1.00>       < 0 -1  0> yields <0.50 0.00>
	//     <0 0 1> yields <0.25 0.50>       < 0  0 -1> yields <0.75 0.50

	minusP := p.Mult(-1)
	theta := math.Acos(minusP.Y())
	phi := math.Atan2(minusP.Z(), p.X()) + math.Pi

	u = phi / (2.0 * math.Pi)
	v = theta / math.Pi

	return u, v
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

func (s *movingSphere) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (s *movingSphere) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
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

func (list *hittableList) pdfValue(o Point3, v Vec3) float64 {

	weight := 1.0 / float64(len(list.objects))
	sum := 0.0

	for _, obj := range list.objects {
		sum += weight * obj.pdfValue(o, v)
	}

	return sum
}

func (list *hittableList) random(o Vec3, rnd *rand.Rand) Vec3 {
	return list.objects[rand.Intn(len(list.objects))].random(o, rnd)
}

type bvhNode struct {
	left, right hittable
	box         aabb
}

func newBvhNode(list []hittable, time0 float64, time1 float64) *bvhNode {
	objs := list

	var bvh bvhNode

	axis := rand.Intn(2 + 1)

	objectSpan := len(objs)

	if objectSpan == 1 {
		bvh.right = objs[0]
		bvh.left = objs[0]
	} else {
		//only part of list
		// comparator.Sort(objs)
		sort.Slice(objs, func(h1, h2 int) bool {
			return boxCompare(objs[h1], objs[h2], axis)
		})

		mid := objectSpan / 2
		bvh.left = newBvhNode(objs[:mid], time0, time1)
		bvh.right = newBvhNode(objs[mid:], time0, time1)
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

func (bvh *bvhNode) boundingBox(time0 float64, time1 float64) (aabb, bool) {
	return bvh.box, true
}

func (bvh *bvhNode) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (bvh *bvhNode) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
}

func boxCompare(a hittable, b hittable, axis int) bool {

	boxA, existsA := a.boundingBox(0, 0)
	boxB, existsB := b.boundingBox(0, 0)

	if !existsA || !existsB {
		panic("No bounding box in bvh_node constructor\n")
	}

	return boxA.minimum[axis] < boxB.minimum[axis]
}

type xyRect struct {
	mat               material
	x0, x1, y0, y1, k float64
}

func (rect *xyRect) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {
	outputBox = aabb{Point3{rect.x0, rect.y0, rect.k - 0.0001}, Point3{rect.x1, rect.y1, rect.k + 0.0001}}
	return outputBox, true
}

func (rect *xyRect) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	t := (rect.k - r.origin.Z()) / r.direction.Z()
	if t < tMin || t > tMax {
		return nil, false
	}

	x := r.origin.X() + t*r.direction.X()
	y := r.origin.Y() + t*r.direction.Y()

	if x < rect.x0 || x > rect.x1 || y < rect.y0 || y > rect.y1 {
		return nil, false
	}

	var rec hitRecord

	rec.u = (x - rect.x0) / (rect.x1 - rect.y0)
	rec.v = (y - rect.y0) / (rect.y1 - rect.y0)
	rec.t = t

	outwardNormal := Vec3{0, 0, 1}
	rec.setFaceNormal(r, outwardNormal)

	rec.mat = rect.mat
	rec.p = r.At(t)

	return &rec, true
}

func (rect *xyRect) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (rect *xyRect) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
}

type xzRect struct {
	mat               material
	x0, x1, z0, z1, k float64
}

func (rect *xzRect) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {
	outputBox = aabb{Point3{rect.x0, rect.k - 0.0001, rect.z0}, Point3{rect.x1, rect.k + 0.0001, rect.z1}}
	return outputBox, true
}

func (rect *xzRect) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	t := (rect.k - r.origin.Y()) / r.direction.Y()
	if t < tMin || t > tMax {
		return nil, false
	}

	x := r.origin.X() + t*r.direction.X()
	z := r.origin.Z() + t*r.direction.Z()

	if x < rect.x0 || x > rect.x1 || z < rect.z0 || z > rect.z1 {
		return nil, false
	}

	var rec hitRecord

	rec.u = (x - rect.x0) / (rect.x1 - rect.x0)
	rec.v = (z - rect.z0) / (rect.z1 - rect.z0)
	rec.t = t

	outwardNormal := Vec3{0, 1, 0}
	rec.setFaceNormal(r, outwardNormal)

	rec.mat = rect.mat
	rec.p = r.At(t)

	return &rec, true
}

func (rect *xzRect) pdfValue(o Point3, v Vec3) float64 {
	rec, hit := rect.hit(&ray{o, v, 0.0}, 0.001, infinity)
	if !hit {
		return 0
	}

	area := (rect.x1 - rect.x0) * (rect.z1 - rect.z0)
	distanceSquared := rec.t * rec.t * v.LengthSquared()
	cosine := math.Abs(v.Dot(rec.normal) / v.Length())

	return distanceSquared / (cosine * area)
}

func (rect *xzRect) random(o Vec3, rnd *rand.Rand) Vec3 {
	randomPoint := Point3{RandomDoubleRange(rect.x0, rect.x1, rnd), rect.k, RandomDoubleRange(rect.z0, rect.z1, rnd)}
	return randomPoint.Sub(o)
}

type yzRect struct {
	mat               material
	y0, y1, z0, z1, k float64
}

func (rect *yzRect) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {
	outputBox = aabb{Point3{rect.k - 0.0001, rect.y0, rect.z0}, Point3{rect.k + 0.0001, rect.y1, rect.z1}}
	return outputBox, true
}

func (rect *yzRect) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	t := (rect.k - r.origin.X()) / r.direction.X()
	if t < tMin || t > tMax {
		return nil, false
	}

	y := r.origin.Y() + t*r.direction.Y()
	z := r.origin.Z() + t*r.direction.Z()

	if y < rect.y0 || y > rect.y1 || z < rect.z0 || z > rect.z1 {
		return nil, false
	}

	var rec hitRecord

	rec.u = (y - rect.y0) / (rect.y1 - rect.y0)
	rec.v = (z - rect.z0) / (rect.z1 - rect.z0)
	rec.t = t

	outwardNormal := Vec3{1, 0, 0}
	rec.setFaceNormal(r, outwardNormal)

	rec.mat = rect.mat
	rec.p = r.At(t)

	return &rec, true
}

func (rect *yzRect) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (rect *yzRect) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
}

type box struct {
	boxMin, boxMax Point3
	sides          hittableList
}

func newBox(p0 Point3, p1 Point3, mat material) *box {
	var b box
	b.boxMin = p0
	b.boxMax = p1

	b.sides.Add(&xyRect{mat, p0.X(), p1.X(), p0.Y(), p1.Y(), p1.Z()})
	b.sides.Add(&xyRect{mat, p0.X(), p1.X(), p0.Y(), p1.Y(), p0.Z()})

	b.sides.Add(&xzRect{mat, p0.X(), p1.X(), p0.Z(), p1.Z(), p1.Y()})
	b.sides.Add(&xzRect{mat, p0.X(), p1.X(), p0.Z(), p1.Z(), p0.Y()})

	b.sides.Add(&yzRect{mat, p0.Y(), p1.Y(), p0.Z(), p1.Z(), p1.X()})
	b.sides.Add(&yzRect{mat, p0.Y(), p1.Y(), p0.Z(), p1.Z(), p0.X()})

	return &b
}

func (b *box) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {
	outputBox = aabb{b.boxMin, b.boxMax}
	return outputBox, true
}

func (b *box) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	return b.sides.hit(r, tMin, tMax)
}

func (b *box) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (b *box) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
}

type translate struct {
	obj    hittable
	offset Vec3
}

func (t *translate) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {
	outputBox, exists = t.obj.boundingBox(time0, time1)
	if !exists {
		return outputBox, false
	}

	outputBox = aabb{
		outputBox.minimum.Add(t.offset),
		outputBox.maximum.Add(t.offset),
	}

	return outputBox, true
}

func (t *translate) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	newRay := ray{r.origin.Sub(t.offset), r.direction, r.time}
	rec, hit := t.obj.hit(&newRay, tMin, tMax)
	if !hit {
		return nil, false
	}

	rec.p = rec.p.Add(t.offset)
	rec.setFaceNormal(&newRay, rec.normal)

	return rec, true
}

func (t *translate) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (t *translate) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
}

type rotateY struct {
	obj                hittable
	sinTheta, cosTheta float64
	hasBox             bool
	box                aabb
}

func newRotateY(obj hittable, angle float64) *rotateY {

	var rot rotateY
	rot.obj = obj
	radians := DegToRad(angle)
	rot.sinTheta = math.Sin(radians)
	rot.cosTheta = math.Cos(radians)
	rot.box, rot.hasBox = obj.boundingBox(0, 1)

	min := Point3{infinity, infinity, infinity}
	max := Point3{-infinity, -infinity, -infinity}

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*rot.box.maximum.X() + (1.0-float64(i))*rot.box.minimum.X()
				y := float64(j)*rot.box.maximum.Y() + (1.0-float64(j))*rot.box.minimum.Y()
				z := float64(k)*rot.box.maximum.Z() + (1.0-float64(k))*rot.box.minimum.Z()

				newX := rot.cosTheta*x + rot.sinTheta*z
				newZ := -rot.sinTheta*x + rot.cosTheta*z

				tester := Vec3{newX, y, newZ}

				for c := 0; c < 3; c++ {
					min[c] = math.Min(min[c], tester[c])
					max[c] = math.Max(max[c], tester[c])
				}
			}
		}
	}

	rot.box = aabb{min, max}

	return &rot
}

func (rot *rotateY) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {
	return rot.box, rot.hasBox
}

func (rot *rotateY) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	origin := r.origin.Copy()
	direction := r.direction.Copy()

	origin[0] = rot.cosTheta*r.origin[0] - rot.sinTheta*r.origin[2]
	origin[2] = rot.sinTheta*r.origin[0] + rot.cosTheta*r.origin[2]

	direction[0] = rot.cosTheta*r.direction[0] - rot.sinTheta*r.direction[2]
	direction[2] = rot.sinTheta*r.direction[0] + rot.cosTheta*r.direction[2]

	rotatedRay := ray{origin, direction, r.time}

	rec, hit := rot.obj.hit(&rotatedRay, tMin, tMax)
	if !hit {
		return nil, false
	}

	p := rec.p.Copy()
	normal := rec.normal.Copy()

	p[0] = rot.cosTheta*rec.p[0] + rot.sinTheta*rec.p[2]
	p[2] = -rot.sinTheta*rec.p[0] + rot.cosTheta*rec.p[2]

	normal[0] = rot.cosTheta*rec.normal[0] + rot.sinTheta*rec.normal[2]
	normal[2] = -rot.sinTheta*rec.normal[0] + rot.cosTheta*rec.normal[2]

	rec.p = p.Copy()
	rec.setFaceNormal(&rotatedRay, normal)

	return rec, true
}

func (rot *rotateY) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (rot *rotateY) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
}

type constantMedium struct {
	boundary      hittable
	phaseFunction material
	negInvDensity float64
}

func newConstantMedium(obj hittable, density float64, tex texture) *constantMedium {
	return &constantMedium{obj, isotropic{tex}, -1 / density}
}

func (m *constantMedium) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {
	return m.boundary.boundingBox(time0, time1)
}

func (m *constantMedium) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	rec1, hit1 := m.boundary.hit(r, -infinity, infinity)
	if !hit1 {
		return nil, false
	}

	rec2, hit2 := m.boundary.hit(r, rec1.t+0.0001, infinity)
	if !hit2 {
		return nil, false
	}

	if rec1.t < tMin {
		rec1.t = tMin
	}
	if rec2.t > tMax {
		rec2.t = tMax
	}

	if rec1.t >= rec2.t {
		return nil, false
	}

	if rec1.t < 0 {
		rec1.t = 0
	}

	rayLength := r.direction.Length()
	distanceInsideBoundary := (rec2.t - rec1.t) * rayLength
	hitDistance := m.negInvDensity * math.Log(rand.Float64())

	if hitDistance > distanceInsideBoundary {
		return nil, false
	}

	var rec hitRecord
	rec.t = rec1.t + hitDistance/rayLength
	rec.p = r.At(rec.t)

	// if rand.Float64() < 0.00001 {
	// fmt.Printf("hitDistance: %f\nrec.t: %f\n", hitDistance, rec.t)
	// rec.p.Print()
	// }

	rec.mat = m.phaseFunction

	return &rec, true
}

func (m *constantMedium) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (m *constantMedium) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
}

type flipFace struct {
	obj hittable
}

func (f *flipFace) boundingBox(time0 float64, time1 float64) (outputBox aabb, exists bool) {
	return f.obj.boundingBox(time0, time1)
}

func (f *flipFace) hit(r *ray, tMin float64, tMax float64) (*hitRecord, bool) {

	rec, hit := f.obj.hit(r, tMin, tMax)
	if !hit {
		return nil, false
	}

	rec.frontFace = !rec.frontFace
	return rec, true
}

func (f *flipFace) pdfValue(o Point3, v Vec3) float64 {
	return 0
}

func (f *flipFace) random(o Vec3, rnd *rand.Rand) Vec3 {
	return Vec3{1, 0, 0}
}
