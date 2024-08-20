package main

import (
	"math"
)

type sphere struct {
	center vec3
	radius float64
}

/*
Calculate the closest Ray sphere intersection point and stores it in a hit record, if there is any.

base on the fact that ray-sphere intersection is effectively a quadratic problem of

(C - P) . (C - P) = r^2

where C is the center, P is the point along the ray, and r is the radius of the sphere. Which unfolds to

(-td + (C-P))^2 = r^2, or just

t^2 * d.d - 2t*d.(C-Q) + (C-Q)^2 - r^2 = 0

solve for discrimanant of t
*/
func (sphere *sphere) hit(ray *ray, ray_tmin float64, ray_tmax float64, record *hit_record) bool {

	oc := sphere.center.Sub(ray.origin)
	a := ray.direction.Length_Squared()
	h := Dot(&ray.direction, oc)
	c := oc.Length_Squared() - (sphere.radius * sphere.radius)

	discriminant := h*h - (a * c)
	if discriminant < 0 {
		return false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range.
	root := (h - sqrtd) / a
	if root <= ray_tmin || ray_tmax <= root {
		root = (h + sqrtd) / a
		if root <= ray_tmin || ray_tmax <= root {
			return false
		}
	}

	record.t = root
	record.point = ray.At(record.t)
	outward_normal := (record.point.Sub(sphere.center).Scale(1.0 / sphere.radius))
	record.set_face_normal(ray, outward_normal)
	return true
}
