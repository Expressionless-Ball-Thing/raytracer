package main

import (
	"math"
	"math/rand/v2"
)

type material interface {
	scatter(ray_in *ray, rec *hit_record, attenuation *vec3, scattered *ray) bool
}

// Lambertian
type lambertian struct {
	albedo vec3
}

func (mat *lambertian) scatter(ray_in *ray, rec *hit_record, attenuation *vec3, scattered *ray) bool {
	// TODO: Make this scatter with some fixed probability p and have attenuation be albedo/p
	scatter_direction := *rec.normal.Add(*Random_unit_vec3())

	// Check to prevent the scenario where the scatter_direction is exactly the opposite of the normal vector.
	if scatter_direction.near_zero() {
		scatter_direction = rec.normal
	}

	scattered.origin = rec.point
	scattered.direction = scatter_direction
	*attenuation = mat.albedo
	return true
}

// Metal
type metal struct {
	albedo vec3
	fuzz   float64
}

func (mat *metal) scatter(ray_in *ray, rec *hit_record, attenuation *vec3, scattered *ray) bool {

	reflected := Reflect(ray_in.direction, rec.normal)
	reflected = *Unit(&reflected).Add(*Random_unit_vec3().Scale(mat.fuzz)) // Adding some fuzz to make the reflections look fuzzy
	scattered.origin = rec.point
	scattered.direction = reflected
	*attenuation = mat.albedo
	return (Dot(&scattered.direction, &rec.normal) > 0)
}

// Dielectric
type dielectric struct {
	// Refractive index in vacuum or air, or the ratio of the material's refractive index over the refractive index of the enclosing media
	refraction_index float64
}

func (mat *dielectric) scatter(ray_in *ray, rec *hit_record, attenuation *vec3, scattered *ray) bool {

	attenuation.X, attenuation.Y, attenuation.Z = 1, 1, 1
	var ri float64
	if rec.front_face {
		ri = 1.0 / mat.refraction_index
	} else {
		ri = mat.refraction_index
	}

	unit_direction := Unit(&ray_in.direction)
	cos_theta := math.Min(Dot(unit_direction.Negate(), &rec.normal), 1.0)
	sin_theta := math.Sqrt(1 - cos_theta*cos_theta)

	cannot_refract := ri*sin_theta > 1.0
	var direction vec3

	// Decide if the ray goes through total internal refraction.
	if cannot_refract || reflectance(cos_theta, ri) > rand.Float64() {
		direction = Reflect(*unit_direction, rec.normal)
	} else {
		direction = *Refract(*unit_direction, rec.normal, ri)
	}

	scattered.origin = rec.point
	scattered.direction = direction
	return true
}

func reflectance(cosine float64, refraction_index float64) float64 {
	r0 := (1 - refraction_index) / (1 + refraction_index)
	r0 *= r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
