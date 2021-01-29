package main

import "math"

type aabb struct {
	minimum Point3
	maximum Point3
}

func (b aabb) hit(r *ray, tMin float64, tMax float64) bool {
	for a := 0; a < 3; a++ {
		// t0 := math.Min((b.minimum[a]-r.origin[a])/r.direction[a],
		// 	(b.maximum[a]-r.origin[a])/r.direction[a])
		// t1 := math.Max((b.minimum[a]-r.origin[a])/r.direction[a],
		// 	(b.maximum[a]-r.origin[a])/r.direction[a])

		// tMin = math.Max(t0, tMin)
		// tMax = math.Min(t1, tMax)

		invD := 1.0 / r.direction[a]
		t0 := (b.minimum[a] - r.origin[a]) * invD
		t1 := (b.maximum[a] - r.origin[a]) * invD

		if invD < 0.0 {
			t0, t1 = t1, t0
		}

		if t0 > tMin {
			tMin = t0
		}

		if t1 < tMax {
			tMax = t1
		}

		if tMax <= tMin {
			return false
		}
	}
	return true
}

func surroundingBox(box0 aabb, box1 aabb) aabb {
	small := Point3{
		math.Min(box0.minimum.X(), box1.minimum.X()),
		math.Min(box0.minimum.Y(), box1.minimum.Y()),
		math.Min(box0.minimum.Z(), box1.minimum.Z()),
	}

	big := Point3{
		math.Max(box0.maximum.X(), box1.maximum.X()),
		math.Max(box0.maximum.Y(), box1.maximum.Y()),
		math.Max(box0.maximum.Z(), box1.maximum.Z()),
	}

	return aabb{small, big}
}
