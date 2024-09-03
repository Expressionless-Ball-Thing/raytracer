package main

import (
	"math"
)

// Translation
type Translate struct {
	object *Hittable
	offset Vec3
	bbox   AABB
}

func NewTranslate(object Hittable, offset *Vec3) *Hittable {
	var thing Hittable = &Translate{
		object: &object,
		offset: *offset,
		bbox:   *(object).bounding_box().AddOffset(offset),
	}

	return &thing
}

func (tran *Translate) hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) (ok bool) {
	// Move the ray backwards by the offset
	offset_ray := NewRay(*ray.origin.Sub(&tran.offset), ray.direction, ray.time)

	// Determine whether an intersection exists along the offset ray (and if so, where)

	if !(*tran.object).hit(&offset_ray, ray_tmin, ray_tmax, record) {
		return false
	}

	record.point.IAdd(&tran.offset)
	return true

}

func (tran *Translate) bounding_box() (bounds *AABB) {
	return &tran.bbox
}

// Rotation, made to be general.
type Rotate struct {
	object                                                         *Hittable
	alpha_cos, beta_cos, gamma_cos, alpha_sin, beta_sin, gamma_sin float64 // Cos and Sine in Radians for Yaw, Pitch and roll
	bbox                                                           AABB
}

func NewRotate(object Hittable, alpha, beta, gamma float64) *Rotate {
	var rotate Rotate

	temp_alpha := alpha * math.Pi / 180
	temp_beta := beta * math.Pi / 180
	temp_gamma := gamma * math.Pi / 180

	rotate.alpha_cos = math.Cos(temp_alpha)
	rotate.alpha_sin = math.Sin(temp_alpha)
	rotate.beta_cos = math.Cos(temp_beta)
	rotate.beta_sin = math.Sin(temp_beta)
	rotate.gamma_cos = math.Cos(temp_gamma)
	rotate.gamma_sin = math.Sin(temp_gamma)

	rotate.object = &object
	rotate.bbox = *(object).bounding_box()

	min := NewVec3(math.Inf(1), math.Inf(1), math.Inf(1))
	max := NewVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1))

	// bbox calculations
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*rotate.bbox.maxVec[0] + (1-float64(i))*rotate.bbox.minVec[0]
				y := float64(j)*rotate.bbox.maxVec[1] + (1-float64(j))*rotate.bbox.minVec[1]
				z := float64(k)*rotate.bbox.maxVec[2] + (1-float64(k))*rotate.bbox.minVec[2]

				tester := rotate.RotateClockwise(NewVec3(x, y, z))

				for c := 0; c < 3; c++ {
					min[c] = math.Min(min[c], tester[c])
					max[c] = math.Max(max[c], tester[c])
				}
			}
		}
	}

	rotate.bbox = *NewAABB(*min, *max)
	return &rotate
}

func (rot *Rotate) hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) (ok bool) {

	origin, direction := ray.origin, ray.direction

	origin = *rot.RotateClockwise(&origin)
	direction = *rot.RotateClockwise(&direction)

	rotated := NewRay(origin, direction, ray.time)

	// Determine whether an intersection exists in object space (and if so, where)
	if !(*rot.object).hit(&rotated, ray_tmin, ray_tmax, record) {
		return false
	}

	// Change the intersection point from object space to world space
	point, normal := record.point, record.normal

	record.point = *rot.RotateAntiClockWise(&point)
	record.normal = *rot.RotateAntiClockWise(&normal)

	return true
}

func (rot *Rotate) bounding_box() (bounds *AABB) {
	return &rot.bbox
}

// Scaling
type Scale struct {
	object    *Hittable
	scale_fac *Vec3 // Scale factor in each axis
	bbox      AABB
}

func NewScale(object Hittable, x, y, z float64) *Hittable {
	var thing Hittable = &Scale{
		object:    &object,
		scale_fac: NewVec3(x, y, z),
		bbox:      *(object).bounding_box().Scale(NewVec3(x, y, z)),
	}

	return &thing
}

func (scale *Scale) hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) (ok bool) {

	new_ray := NewRay(*ray.origin.Div(scale.scale_fac), *ray.direction.Div(scale.scale_fac), ray.time)

	// Determine whether an intersection exists in object space (and if so, where)
	if !(*scale.object).hit(&new_ray, ray_tmin, ray_tmax, record) {
		return false
	}
	// Change the intersection point from object space to world space
	record.point = *record.point.Mult(scale.scale_fac)
	record.normal = *record.normal.Mult(scale.scale_fac)
	return true
}

func (scale *Scale) bounding_box() (bounds *AABB) {
	return &scale.bbox
}

// Shearing
type Shear struct {
	object                       *Hittable
	x_y, x_z, y_x, y_z, z_x, z_y float64 // (axis to move)_(axis to move in relation to)
	bbox                         AABB
}

func NewShear(object Hittable, x_y, x_z, y_x, y_z, z_x, z_y float64) *Shear {

	var thing *Shear = &Shear{
		&object,
		x_y, x_z, y_x, y_z, z_x, z_y,
		*(object).bounding_box(),
	}

	thing.bbox = *NewAABB(*thing.ApplyShear(&object.bounding_box().minVec), *thing.ApplyShear(&object.bounding_box().maxVec))

	return thing
}

func (shear *Shear) ApplyShear(v1 *Vec3) *Vec3 {
	return &Vec3{
		v1[0] + v1[1]*shear.x_y + v1[2]*shear.x_z,
		v1[0]*shear.y_x + v1[1] + v1[2]*shear.y_z,
		v1[0]*shear.z_x + v1[1]*shear.z_y + v1[2],
	}
}

func (shear *Shear) ReveseShear(v1 *Vec3) *Vec3 {
	a, b, c, d, e, f := shear.x_y, shear.x_z, shear.y_x, shear.y_z, shear.z_x, shear.z_y
	factor := (1 - a*c - b*e + a*d*e + b*c*f - d*f)
	return &Vec3{
		v1[0]*(1-d*f)/factor + v1[1]*(-a+b*f)/factor + v1[2]*(-b+a*d)/factor,
		v1[0]*(-c+d*e)/factor + v1[1]*(1-b*e)/factor + v1[2]*(b*c-d)/factor,
		v1[0]*(-e+c*f)/factor + v1[1]*(a*e-f)/factor + v1[2]*(1-a*c)/factor,
	}
}

func (shear *Shear) hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) (ok bool) {

	new_ray := NewRay(*shear.ReveseShear(&ray.origin), *shear.ReveseShear(&ray.direction), ray.time)

	// Determine whether an intersection exists in object space (and if so, where)
	if !(*shear.object).hit(&new_ray, ray_tmin, ray_tmax, record) {
		return false
	}
	// Change the intersection point from object space to world space
	record.point = *shear.ApplyShear(&record.point)
	record.normal = *shear.ApplyShear(&record.normal)
	return true
}

func (shear *Shear) bounding_box() (bounds *AABB) {
	return &shear.bbox
}
