package main

import (
	"math"
	"math/rand/v2"
)

// 3 Element Vector
type Vec3 [3]float64

// Creates a Vec3 from 3 float values
func NewVec3(e0, e1, e2 float64) *Vec3 {
	return &Vec3{e0, e1, e2}
}

func NewVec3Random(min, max float64) *Vec3 {
	return &Vec3{Random_float64_bounded(min, max), Random_float64_bounded(min, max), Random_float64_bounded(min, max)}
}

// X returns the first element
func (v *Vec3) X() float64 {
	return v[0]
}

// Y returns the second element
func (v *Vec3) Y() float64 {
	return v[1]
}

// Z returns the third element
func (v *Vec3) Z() float64 {
	return v[2]
}

// Negate a vec3
func (v1 *Vec3) Negate() *Vec3 {
	return NewVec3(-v1[0], -v1[1], -v1[2])
}

// Returns the square of the Vec3's length
func (v1 *Vec3) Length_Squared() float64 {
	return v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2]
}

// Return the magnitude of the Vec3.
func (v1 *Vec3) Magnitude() float64 {
	return math.Sqrt(v1.Length_Squared())
}

// In-place vec3 addition
func (v1 *Vec3) IAdd(v2 *Vec3) {
	v1[0] += v2[0]
	v1[1] += v2[1]
	v1[2] += v2[2]
}

// In-place vec3 scalar multiplication
func (v1 *Vec3) IScale(t float64) {
	v1[0] *= t
	v1[1] *= t
	v1[2] *= t
}

func (v1 *Vec3) near_zero() bool {
	s := 1e-8
	return (math.Abs(v1[0]) < s) && (math.Abs(v1[1]) < s) && (math.Abs(v1[2]) < s)
}

// Unit vector of the vec3
func (v1 *Vec3) Unit() *Vec3 {
	return v1.Scale(1 / v1.Magnitude())
}

// Aux functions

// Add two vec3
func (v1 *Vec3) Add(v2 *Vec3) *Vec3 {
	return &Vec3{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

// Sub two Vec3
func (v1 *Vec3) Sub(v2 *Vec3) *Vec3 {
	return &Vec3{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

// Multiply two Vec3s element by element
func (v1 *Vec3) Mult(v2 *Vec3) *Vec3 {
	return &Vec3{v1[0] * v2[0], v1[1] * v2[1], v1[2] * v2[2]}
}

// Divide two vectors element by element
func (v1 *Vec3) Div(v2 *Vec3) *Vec3 {
	return &Vec3{v1[0] / v2[0], v1[1] / v2[1], v1[2] / v2[2]}
}

// Scalar Multiply a Vec3
func (v1 *Vec3) Scale(t float64) *Vec3 {
	return &Vec3{
		v1[0] * t,
		v1[1] * t,
		v1[2] * t,
	}
}

// Gamma raises each of R, G, and B to 1/n
func (v1 *Vec3) Gamma(n float64) *Vec3 {
	ni := 1 / n
	return NewVec3(
		math.Pow(v1[0], ni),
		math.Pow(v1[1], ni),
		math.Pow(v1[2], ni),
	)
}

// Dot product of 2 Vec3
func Dot(v1 *Vec3, v2 *Vec3) float64 {
	return (v1[0] * v2[0]) + (v1[1] * v2[1]) + (v1[2] * v2[2])
}

// Cross product of 2 Vec3
func Cross(v1 *Vec3, v2 *Vec3) *Vec3 {
	return &Vec3{
		v1[1]*v2[2] - v1[2]*v2[1],
		v1[2]*v2[0] - v1[0]*v2[2],
		v1[0]*v2[1] - v1[1]*v2[0],
	}
}

// Generate a random vector3 whose elements are within bounds
func RandomVec3(min float64, max float64) Vec3 {
	return Vec3{min + (max-min)*rand.Float64(), min + (max-min)*rand.Float64(), min + (max-min)*rand.Float64()}
}

func Random_float64_bounded(min float64, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

// Generates a random unit Vec3 with length of 1.
func Random_unit_Vec3() *Vec3 {
	for {
		t := RandomVec3(-1, 1)
		if t.Length_Squared() < 1 {
			return t.Unit()
		}
	}
}

func Random_on_hemisphere(normal *Vec3) *Vec3 {
	on_unit_sphere := Random_unit_Vec3()
	if Dot(on_unit_sphere, normal) > 0 { // In the same hemisphere as the normal
		return on_unit_sphere
	} else {
		return on_unit_sphere.Negate()
	}
}

// Compute the reflected vector3 if incident vector v hits a reflective surface with normal vector n
func Reflect(v *Vec3, n *Vec3) *Vec3 {
	return v.Sub(n.Scale(2 * Dot(v, n)))
}

// Compute the refractred vector3 if incident vector uv hits a reflective surface with normal vector n and the refractive indexes have a ratio of etai_over_etat
func Refract(uv *Vec3, n *Vec3, etai_over_etat float64) *Vec3 {
	cos_theta := math.Min(Dot(uv.Negate(), n), 1.0)
	r_out_perp := (uv.Add(n.Scale(cos_theta))).Scale(etai_over_etat)
	r_out_parallel := n.Scale(-math.Sqrt(math.Abs(1.0 - r_out_perp.Length_Squared())))
	return r_out_perp.Add(r_out_parallel)
}

func random_in_unit_disk() *Vec3 {
	for {
		p := &Vec3{Random_float64_bounded(-1, 1), Random_float64_bounded(-1, 1), 0}
		if p.Length_Squared() < 1 {
			return p
		}
	}
}

// Rotation Matrix 3D, all should be radians
func (v1 Vec3) RotateGen(alpha, beta, gamma float64) *Vec3 {
	return &Vec3{
		math.Cos(alpha)*math.Cos(beta)*v1[0] + (math.Cos(alpha)*math.Sin(beta)*math.Sin(gamma)-math.Sin(alpha)*math.Cos(gamma))*v1[1] + (math.Cos(alpha)*math.Sin(beta)*math.Cos(gamma)+math.Sin(alpha)*math.Sin(gamma))*v1[2],
		math.Sin(alpha)*math.Cos(beta)*v1[0] + (math.Sin(alpha)*math.Sin(beta)*math.Sin(gamma)+math.Cos(alpha)*math.Cos(gamma))*v1[1] + (math.Sin(alpha)*math.Sin(beta)*math.Cos(gamma)-math.Cos(alpha)*math.Sin(gamma))*v1[2],
		math.Sin(beta)*(-v1[0]) + (math.Cos(beta)*math.Sin(gamma))*v1[1] + (math.Cos(beta)*math.Cos(gamma))*v1[2],
	}
}

func (rotate *Rotate) RotateAntiClockWise(v1 *Vec3) *Vec3 {
	return &Vec3{
		rotate.alpha_cos*rotate.beta_cos*v1[0] + ((rotate.alpha_cos*rotate.beta_sin*rotate.gamma_sin)-(rotate.alpha_sin*rotate.gamma_cos))*v1[1] + ((rotate.alpha_cos*rotate.beta_sin*rotate.gamma_cos)+(rotate.alpha_sin*rotate.gamma_sin))*v1[2],
		rotate.alpha_sin*rotate.beta_cos*v1[0] + ((rotate.alpha_sin*rotate.beta_sin*rotate.gamma_sin)+(rotate.alpha_cos*rotate.gamma_cos))*v1[1] + ((rotate.alpha_sin*rotate.beta_sin*rotate.gamma_cos)-(rotate.alpha_cos*rotate.gamma_sin))*v1[2],
		-rotate.beta_sin*(v1[0]) + (rotate.beta_cos*rotate.gamma_sin)*v1[1] + (rotate.beta_cos * rotate.gamma_cos * v1[2]),
	}
}

func (rotate *Rotate) RotateClockwise(v1 *Vec3) *Vec3 {
	return &Vec3{
		rotate.alpha_cos*rotate.beta_cos*v1[0] + ((rotate.alpha_cos*rotate.beta_sin*rotate.gamma_sin)+(rotate.alpha_sin*rotate.gamma_cos))*v1[1] + (-(rotate.alpha_cos*rotate.beta_sin*rotate.gamma_cos)+(rotate.alpha_sin*rotate.gamma_sin))*v1[2],
		-rotate.alpha_sin*rotate.beta_cos*v1[0] + (-(rotate.alpha_sin*rotate.beta_sin*rotate.gamma_sin)+(rotate.alpha_cos*rotate.gamma_cos))*v1[1] + ((rotate.alpha_sin*rotate.beta_sin*rotate.gamma_cos)+(rotate.alpha_cos*rotate.gamma_sin))*v1[2],
		rotate.beta_sin*(v1[0]) - (rotate.beta_cos*rotate.gamma_sin)*v1[1] + (rotate.beta_cos * rotate.gamma_cos * v1[2]),
	}
}
