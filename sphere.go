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

	var rec Hit
	rec.t = root
	rec.point = ray.At(rec.t)
	rec.normal = (rec.point.Sub(sphere.center)).Scale(1.0 / sphere.radius)
	rec.set_face_normal(ray, rec.normal)
	rec.material = sphere.material
	rec.u, rec.v = get_sphere_uv(rec.normal)
	return rec, true
}

// Return the Location of the sphere center at a given time if it's moving.
func (sphere *Sphere) sphere_center(time float64) Vec3 {
	return sphere.center.Add(sphere.center_vec.Scale(time))
}

// Returns the bounding box of the sphere
func (sphere *Sphere) bounding_box() (bounds *AABB) {
	return &sphere.bbox
}

// Get the UV coordinates relative to a sphere given a Vec3.
// p: a given point on the sphere of radius one, centered at the origin.
// u: returned value [0,1] of angle around the Y axis from X=-1.
// v: returned value [0,1] of angle from Y=-1 to Y=+1.
//
//	<1 0 0> yields <0.50 0.50>       <-1  0  0> yields <0.00 0.50>
//	<0 1 0> yields <0.50 1.00>       < 0 -1  0> yields <0.50 0.00>
//	<0 0 1> yields <0.25 0.50>       < 0  0 -1> yields <0.75 0.50>
func get_sphere_uv(point Vec3) (u, v float64) {
	theta := math.Acos(-point[1])
	phi := math.Atan2(-point[2], point[0]) + math.Pi

	return phi / (2 * math.Pi), theta / math.Pi
}

func (sphere *Sphere) count() int {
	return 1
}
