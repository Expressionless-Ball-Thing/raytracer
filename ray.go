package main

// Represents a Ray
type Ray struct {
	origin    Vec3
	direction Vec3 // Assumed Unit Vector
	time      float64
}

// Creates a new ray given an origin and a direction
func NewRay(origin Vec3, direction Vec3, time float64) Ray {
	return Ray{origin, *direction.Unit(), time}
}

// At returns the ray at point t (lerp in a way)
func (r Ray) At(t float64) Vec3 {
	return *r.origin.Add(r.direction.Scale(t))
}
