package main

import (
	"math"
	"math/rand/v2"
)

type Material interface {

	// Calcuate the light color emitted from that point
	emitted(u, v float64, point *Vec3) Vec3

	// Given incident direction and the Normal of the surface, calculate the reflected direction and the attenuation
	scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool)
}

// Lambert describes a diffuse material.
type Lambert struct {
	texture Texture
}

// NewLambert creates a new Lambert material with the given color.
func NewLambert(albedo Vec3) *Lambert {
	return &Lambert{NewSolidTexture(albedo)}
}

func NewLambertTex(texture Texture) *Lambert {
	return &Lambert{texture}
}

// Scatter scatters incoming light rays in a hemisphere about the normal.
func (lambert *Lambert) scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool) {

	// TODO: scatter with some fixed probability p and have attenuation be albedo/p .

	scatter_direction := hit.normal.Add(Random_unit_Vec3()) // added a normal vector to make it closer to the surface normal
	if scatter_direction.near_zero() {
		scatter_direction = hit.normal
	}

	return scatter_direction, lambert.texture.value(hit.u, hit.v, hit.point), true
}

func (lambert *Lambert) emitted(u, v float64, point *Vec3) Vec3 {
	return NewVec3(0, 0, 0)
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

func (metal *Metal) scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool) {

	reflected := Reflect(in, hit.normal).Unit().Add(Random_unit_Vec3().Scale(metal.fuzz)) // Adding some fuzz to make the reflections look fuzzy
	return reflected, metal.albedo, (Dot(reflected, hit.normal) > 0)
}

func (metal *Metal) emitted(u, v float64, point *Vec3) Vec3 {
	return NewVec3(0, 0, 0)
}

type Dielectric struct {
	// Refractive index in vacuum or air, or the ratio of the material's refractive index over the refractive index of the enclosing media
	refraction_index float64
}

func NewDielectric(refraction_index float64) *Dielectric {
	return &Dielectric{refraction_index}
}

func (dielec *Dielectric) scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool) {

	var ri float64
	if hit.front_face {
		ri = 1.0 / dielec.refraction_index
	} else {
		ri = dielec.refraction_index
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

func (dielec *Dielectric) emitted(u, v float64, point *Vec3) Vec3 {
	return NewVec3(0, 0, 0)
}

// Use Schlick's approximation for reflectance.
func reflectance(cosine float64, refraction_index float64) float64 {
	r0 := (1 - refraction_index) / (1 + refraction_index)
	r0 *= r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}

/**
Lights
*/

// Diffuse Light
type DiffuseLight struct {
	texture *Texture
}

// Creates Diffuse Light given a texture
func NewDiffuseLight(texture Texture) *DiffuseLight {
	return &DiffuseLight{
		texture: &texture,
	}
}

// Creates Diffuse Light given the light color
func NewDiffuseLightColor(emit Vec3) *DiffuseLight {
	var tex Texture = NewSolidTexture(emit)
	return &DiffuseLight{
		texture: &tex,
	}
}

// Determines the light emitted at a certain point
func (diffuse *DiffuseLight) emitted(u, v float64, point *Vec3) Vec3 {
	return (*diffuse.texture).value(u, v, *point)
}

func (diffuse *DiffuseLight) scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool) {
	return out, attenuation, false
}

// Isotropic Material
type Isotropic struct {
	tex *Texture
}

func NewIsotropic(albedo *Vec3) *Isotropic {
	var temp Texture = NewSolidTexture(*albedo)
	return &Isotropic{
		tex: &temp,
	}
}

func NewIsotropicFromTexture(texture *Texture) *Isotropic {
	return &Isotropic{
		tex: texture,
	}
}

func (iso *Isotropic) emitted(u, v float64, point *Vec3) Vec3 {
	return NewVec3(0, 0, 0)
}

func (iso *Isotropic) scatter(in Vec3, hit Hit) (out Vec3, attenuation Vec3, ok bool) {
	return Random_unit_Vec3(), (*iso.tex).value(hit.u, hit.v, hit.point), true
}
