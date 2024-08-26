package main

import (
	"math"
	"math/rand/v2"
)

type Material interface {

	// Given incident direction and the Normal of the surface, calculate the reflected direction and the attenuation
	scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool)
}

// Lambert describes a diffuse material.
type Lambert struct {
	Albedo Vec3
}

// NewLambert creates a new Lambert material with the given color.
func NewLambert(albedo Vec3) *Lambert {
	return &Lambert{Albedo: albedo}
}

// Scatter scatters incoming light rays in a hemisphere about the normal.
func (lambert *Lambert) scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool) {

	// TODO: scatter with some fixed probability p and have attenuation be albedo/p .

	scatter_direction := hit.normal.Add(Random_unit_Vec3()) // added a normal vector to make it closer to the surface normal
	if scatter_direction.near_zero() {
		scatter_direction = hit.normal
	}
	return scatter_direction, lambert.Albedo, true
}

// Metal
type Metal struct {
	albedo Vec3
	fuzz   float64
}

// NewMetal creates a new Metal material with the given color and fuzz factor
func NewMetal(albedo Vec3, fuzz float64) *Metal {
	return &Metal{albedo, fuzz}
}

func (mat *Metal) scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool) {

	reflected := Reflect(in, hit.normal).Unit().Add(Random_unit_Vec3().Scale(mat.fuzz)) // Adding some fuzz to make the reflections look fuzzy
	return reflected, mat.albedo, (Dot(reflected, hit.normal) > 0)
}

type Dielectric struct {
	// Refractive index in vacuum or air, or the ratio of the material's refractive index over the refractive index of the enclosing media
	refraction_index float64
}

func NewDielectric(refraction_index float64) *Dielectric {
	return &Dielectric{refraction_index}
}

func (mat *Dielectric) scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool) {

	var ri float64
	if hit.front_face {
		ri = 1.0 / mat.refraction_index
	} else {
		ri = mat.refraction_index
	}

	unit_direction := in.Unit()

	cos_theta := math.Min(Dot(unit_direction.Negate(), hit.normal), 1.0)
	sin_theta := math.Sqrt(1 - cos_theta*cos_theta)

	cannot_refract := ri*sin_theta > 1.0

	// Decide if the ray goes through total internal refraction.
	if cannot_refract || reflectance(cos_theta, ri) > rand.Float64() {
		out = Reflect(unit_direction, hit.normal)
	} else {
		out = Refract(unit_direction, hit.normal, ri)
	}

	return out, NewVec3(1, 1, 1), true
}

// Use Schlick's approximation for reflectance.
func reflectance(cosine float64, refraction_index float64) float64 {
	r0 := (1 - refraction_index) / (1 + refraction_index)
	r0 *= r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
