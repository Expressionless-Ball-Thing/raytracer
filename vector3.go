package main

import (
	"math"
	"math/rand/v2"
)

type vec3 struct {
	X float64
	Y float64
	Z float64
}

// Negate a vec3
func (v1 *vec3) Negate() *vec3 {
	return &vec3{-v1.X, -v1.Y, -v1.Z}
}

// In-place vec3 addition
func (v1 *vec3) IAdd(v2 vec3) {
	v1.X += v2.X
	v1.Y += v2.Y
	v1.Z += v2.Z
}

// In-place vec3 scalar multiplication
func (v1 *vec3) IScale(t float64) {
	v1.X *= t
	v1.Y *= t
	v1.Z *= t
}

// Return the magnitude of the vec3.
func (v1 *vec3) Magnitude() float64 {
	return math.Sqrt(v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z)
}

func (v1 *vec3) Length_Squared() float64 {
	return (v1.X * v1.X) + (v1.Y * v1.Y) + (v1.Z * v1.Z)
}

func (v1 *vec3) near_zero() bool {
	s := 1e-8
	return (math.Abs(v1.X) < s) && (math.Abs(v1.Y) < s) && (math.Abs(v1.Z) < s)
}

// Aux functions

// Add two vec3
func (v1 *vec3) Add(v2 vec3) *vec3 {
	return &vec3{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

// Sub two vec3
func (v1 *vec3) Sub(v2 vec3) *vec3 {
	return &vec3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

// Multiply two vec3s element by element
func (v1 *vec3) Mult(v2 vec3) *vec3 {
	return &vec3{v1.X * v2.X, v1.Y * v2.Y, v1.Z * v2.Z}
}

// Scalar Multiply a vec3
func (v1 *vec3) Scale(t float64) *vec3 {
	return &vec3{
		v1.X * t,
		v1.Y * t,
		v1.Z * t,
	}
}

// Dot product of 2 vec3
func Dot(v1 *vec3, v2 *vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// Cross product of 2 vec3
func Cross(v1 *vec3, v2 *vec3) *vec3 {
	return &vec3{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
	}
}

// Unit vector of the vec3
func Unit(v1 *vec3) *vec3 {

	return &vec3{
		v1.X * 1 / v1.Magnitude(),
		v1.Y * 1 / v1.Magnitude(),
		v1.Z * 1 / v1.Magnitude(),
	}
}

// Generate a random vector3 whose elements are within bounds
func RandomVec3(min float64, max float64) *vec3 {
	return &vec3{min + (max-min)*rand.Float64(), min + (max-min)*rand.Float64(), min + (max-min)*rand.Float64()}
}

func random_float64_bounded(min float64, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

// Generates a random unit vec3 with length of 1.
func Random_unit_vec3() *vec3 {
	return Unit(RandomVec3(-1, 1))
}

func Random_on_hemisphere(normal *vec3) *vec3 {
	on_unit_sphere := Random_unit_vec3()
	if Dot(on_unit_sphere, normal) > 0 { // In the same hemisphere as the normal
		return on_unit_sphere
	} else {
		return on_unit_sphere.Negate()
	}
}

// Compute the reflected vector3 if incident vector v hits a reflective surface with normal vector n
func Reflect(v vec3, n vec3) vec3 {
	return *(v.Sub(*n.Scale(2 * Dot(&v, &n))))
}

// Compute the refractred vector3 if incident vector uv hits a reflective surface with normal vector n and the refractive indexes have a ratio of etai_over_etat
func Refract(uv vec3, n vec3, etai_over_etat float64) *vec3 {
	cos_theta := math.Min(Dot(uv.Negate(), &n), 1.0)
	r_out_perp := (uv.Add(*n.Scale(cos_theta))).Scale(etai_over_etat)
	r_out_parallel := n.Scale(-math.Sqrt(math.Abs(1.0 - r_out_perp.Length_Squared())))
	return r_out_perp.Add(*r_out_parallel)
}

func random_in_unit_disk() vec3 {
	for {
		p := vec3{random_float64_bounded(-1, 1), random_float64_bounded(-1, 1), 0}
		if p.Length_Squared() < 1 {
			return p
		}
	}
}
