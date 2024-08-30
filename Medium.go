package main

import (
	"math"
	"math/rand/v2"
)

// Constant medium
type Constant struct {
	boundary        *Hittable
	neg_inv_density float64
	phase_function  *Material
}

func NewConstantMedium(boundary *Hittable, density float64, tex *Texture) *Constant {
	var temp Material = NewIsotropicFromTexture(tex)
	return &Constant{
		boundary:        boundary,
		neg_inv_density: -1.0 / density,
		phase_function:  &temp,
	}
}

func NewConstantMediumAlbedo(boundary *Hittable, density float64, albedo Vec3) *Constant {
	var temp Material = NewIsotropic(&albedo)
	return &Constant{
		boundary:        boundary,
		neg_inv_density: -1.0 / density,
		phase_function:  &temp,
	}
}

func (constant *Constant) hit(ray *Ray, ray_tmin float64, ray_tmax float64) (record Hit, ok bool) {

	rec1, ok1 := (*constant.boundary).hit(ray, math.Inf(-1), math.Inf(1))
	if !ok1 {
		return record, false
	}

	rec2, ok2 := (*constant.boundary).hit(ray, rec1.t+0.0001, math.Inf(1))
	if !ok2 {
		return record, false
	}

	if rec1.t < ray_tmin {
		rec1.t = ray_tmin
	}
	if rec2.t > ray_tmax {
		rec2.t = ray_tmax
	}

	if rec1.t >= rec2.t {
		return record, false
	}

	if rec1.t < 0 {
		rec1.t = 0
	}

	ray_length := ray.direction.Magnitude()
	distance_inside_boundary := (rec2.t - rec1.t) * ray_length
	hit_distance := constant.neg_inv_density * math.Log(rand.Float64())

	if hit_distance > distance_inside_boundary {
		return record, false
	}

	var rec Hit
	rec.t = rec1.t + (hit_distance / ray_length)
	rec.point = ray.At(rec.t)

	rec.normal = NewVec3(1, 0, 0) // arbitrary
	rec.front_face = true         // arbitrary
	rec.material = *constant.phase_function

	return rec, true
}

func (constant *Constant) bounding_box() (bounds *AABB) {
	return (*constant.boundary).bounding_box()
}
