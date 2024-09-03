package main

import "math"

type AABB struct {
	minVec, maxVec Vec3 // Min and Max of the AABB, basically the intervals in the x, y and z intervals condensed together.
}

// Create a New AABB given two extrema points
func NewAABB(a, b Vec3) *AABB {

	var minVec Vec3
	var maxVec Vec3
	for i := 0; i < 3; i++ {
		if b[i] < a[i] {
			minVec[i] = b[i]
			maxVec[i] = a[i]
		} else {
			minVec[i] = a[i]
			maxVec[i] = b[i]
		}
	}

	delta := 0.0001

	for i := 0; i < 3; i++ {
		if maxVec[i]-minVec[i] < delta {
			minVec[i] -= delta / 2
			maxVec[i] += delta / 2
		}
	}

	return &AABB{minVec, maxVec}

}

// Merge two AABBs together
func MergeAABB(box1, box2 AABB) *AABB {

	var box3 AABB
	for axis := 0; axis < 3; axis++ {
		if box1.minVec[axis] <= box2.minVec[axis] {
			box3.minVec[axis] = box1.minVec[axis]
		} else {
			box3.minVec[axis] = box2.minVec[axis]
		}

		if box1.maxVec[axis] >= box2.maxVec[axis] {
			box3.maxVec[axis] = box1.maxVec[axis]
		} else {
			box3.maxVec[axis] = box2.maxVec[axis]
		}
	}

	return &box3
}

// Hit returns whether or not r hits the box between distances ray_tmin and ray_tmax
func (aabb *AABB) hit_test(ray *Ray, ray_tmin float64, ray_tmax float64) (ok bool) {

	for axis := 0; axis < 3; axis++ {

		adinv := 1.0 / ray.direction[axis]

		t0 := (aabb.minVec[axis] - ray.origin[axis]) * adinv
		t1 := (aabb.maxVec[axis] - ray.origin[axis]) * adinv

		if t0 < t1 {
			if t0 > ray_tmin {
				ray_tmin = t0
			}
			if t1 > ray_tmax {
				ray_tmax = t1
			}
		} else {
			if t1 > ray_tmin {
				ray_tmin = t1
			}
			if t0 > ray_tmax {
				ray_tmax = t0
			}
		}

		if ray_tmax <= ray_tmin {
			return false
		}

	}

	return true
}

// Modify AABB by an offset
func (bbox *AABB) AddOffset(offset *Vec3) *AABB {
	return NewAABB(*bbox.minVec.Add(offset), *bbox.maxVec.Add(offset))
}

// Scale an AABB
func (bbox *AABB) Scale(scale_fac *Vec3) *AABB {
	return NewAABB(*bbox.minVec.Mult(scale_fac), *bbox.maxVec.Mult(scale_fac))
}

func NewEmptyAABB() *AABB {
	return &AABB{*NewVec3(math.Inf(1), math.Inf(1), math.Inf(1)), *NewVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1))}
}

func NewUniverseAABB() *AABB {
	return &AABB{*NewVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1)), *NewVec3(math.Inf(1), math.Inf(1), math.Inf(1))}
}
