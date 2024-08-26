package main

import (
	"math"
	"sort"
)

type BVH struct {
	left, right *Hittable
	aabb        AABB
}

func NewBVHNode(objects []Hittable) *BVH {

	// Build the bounding box of the span of source objects.
	var node BVH
	node.aabb = AABB{NewVec3(math.Inf(1), math.Inf(1), math.Inf(1)), NewVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1))}
	for _, v := range objects {
		node.aabb = *MergeAABB(node.aabb, *v.bounding_box())
	}

	var longest int

	if node.aabb.maxVec[0]-node.aabb.minVec[0] > node.aabb.maxVec[1]-node.aabb.minVec[1] {
		if node.aabb.maxVec[0]-node.aabb.minVec[0] > node.aabb.maxVec[2]-node.aabb.minVec[2] {
			longest = 0
		} else {
			longest = 2
		}
	} else {
		if node.aabb.maxVec[1]-node.aabb.minVec[1] > node.aabb.maxVec[2]-node.aabb.minVec[2] {
			longest = 1
		} else {
			longest = 2
		}
	}

	if len(objects) == 1 {
		node.left = &objects[0]
		node.right = &objects[0]
	} else if len(objects) == 2 {
		node.left = &objects[0]
		node.right = &objects[1]
	} else {
		sort.Slice(objects, func(a, b int) bool {
			return objects[a].bounding_box().minVec[longest] < objects[b].bounding_box().minVec[longest]
		})

		mid := len(objects) / 2
		var thing Hittable = NewBVHNode(objects[0:mid])
		node.left = &thing
		var thing2 Hittable = NewBVHNode(objects[mid:])
		node.right = &thing2
	}

	return &node
}

func (node *BVH) hit(ray *Ray, ray_tmin float64, ray_tmax float64) (record Hit, ok bool) {
	if !(node.aabb.hit_test(ray, ray_tmin, ray_tmax)) {
		return record, false
	}

	var hit_left, hit_right bool
	record, hit_left = (*node.left).hit(ray, ray_tmin, ray_tmax)
	var right_rec Hit
	if hit_left {
		right_rec, hit_right = (*node.right).hit(ray, ray_tmin, record.t)
	} else {
		right_rec, hit_right = (*node.right).hit(ray, ray_tmin, ray_tmax)
	}

	if hit_right {
		record = right_rec
	}

	return record, (hit_left || hit_right)

}

func (node *BVH) count() int {
	return 1 + (*node.left).count() + (*node.right).count()
}

func (node *BVH) bounding_box() (bounds *AABB) {
	return &node.aabb
}
