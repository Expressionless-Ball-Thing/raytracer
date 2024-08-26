package main

import (
	"math"
)

// A Sphere
type Sphere struct {
	center     Vec3
	radius     float64
	material   Material
	is_moving  bool // Boolean to see if the Sphere is Moving
	center_vec Vec3 // The vector that indicates how the sphere's center is moving
	bbox       AABB // Bounding Box of the sphere
}

// Creates a New Sphere
func NewSphere(center Vec3, radius float64, material Material) *Sphere {

	rvec := NewVec3(radius, radius, radius)

	return &Sphere{
		center,
		radius,
		material,
		false,
		NewVec3(0, 0, 0),
		*NewAABB(center.Sub(rvec), center.Add(rvec)),
	}
}

// Creates a Moving Sphere given the two centers that it is moving in between.
func NewMovingSphere(center1, center2 Vec3, radius float64, material Material) *Sphere {

	rvec := NewVec3(radius, radius, radius)

	// Creating a new AABB by merging the two AABBs from the moving sphere's start and end points.
	box1, box2 := NewAABB(center1.Sub(rvec), center1.Add(rvec)), NewAABB(center2.Sub(rvec), center2.Add(rvec))
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

	return &Sphere{center1, radius, material, true, center2.Sub(center1), box3}
}

func (sphere *Sphere) hit(ray *Ray, ray_tmin float64, ray_tmax float64) (record Hit, ok bool) {

	var center Vec3
	if sphere.is_moving {
		center = sphere.sphere_center(ray.time)
	} else {
		center = sphere.center
	}

	oc := center.Sub(ray.origin)
	a := ray.direction.Length_Squared()
	h := Dot(ray.direction, oc)
	c := oc.Length_Squared() - (sphere.radius * sphere.radius)

	discriminant := h*h - (a * c)
	if discriminant < 0 {
		return record, false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range.
	root := (h - sqrtd) / a
	if root <= ray_tmin || ray_tmax <= root {
		root = (h + sqrtd) / a
		if root <= ray_tmin || ray_tmax <= root {
			return record, false
		}
	}

	record.t = root
	record.point = ray.At(record.t)
	record.normal = (record.point.Sub(sphere.center)).Scale(1.0 / sphere.radius)
	record.set_face_normal(ray, record.normal)
	record.material = sphere.material
	return record, true
}

// Return the Location of the sphere center at a given time if it's moving.
func (sphere *Sphere) sphere_center(time float64) Vec3 {
	return sphere.center.Add(sphere.center_vec.Scale(time))
}

// Returns the bounding box of the sphere
func (sphere *Sphere) bounding_box() (bounds *AABB) {
	return &sphere.bbox
}

func (sphere *Sphere) count() int {
	return 1
}
