package main

import "math"

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
func (v1 *vec3) IAdd(v2 *vec3) {
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

// Aux functions

// Add two vec3
func (v1 *vec3) Add(v2 vec3) *vec3 {
	return &vec3{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

// Sub two vec3
func (v1 *vec3) Sub(v2 vec3) *vec3 {
	return &vec3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
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
