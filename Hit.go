package main

import (
	"math"
)

type Hittable interface {

	// Calculates whether a hit can be made with the object within the given bounds and alters the record.
	hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) (ok bool)

	bounding_box() (bounds *AABB)
}

// Struct to store the details of a ray hitting a surface
type Hit struct {
	point      Vec3
	normal     Vec3
	t          float64
	u, v       float64 // surface coordinates of the ray-object hit point.
	front_face bool    // Hack way to check front_face or not Dot(&in, &n) < 0
	material   *Material
}

// Sets the hit record normal vector.
// NOTE: the parameter `outward_normal` is assumed to have unit length.
func (record *Hit) set_face_normal(ray *Ray, outward_normal Vec3) {
	record.front_face = Dot(&ray.direction, &outward_normal) < 0
	if record.front_face {
		record.normal = outward_normal
	} else {
		record.normal = (*outward_normal.Negate())
	}
}

type Hit_List struct {
	list []Hittable
	aabb AABB
}

func NewList(objects ...Hittable) *Hit_List {

	var list Hit_List
	list.aabb = AABB{*NewVec3(math.Inf(1), math.Inf(1), math.Inf(1)), *NewVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1))}
	list.list = objects
	for _, v := range objects {
		list.aabb = *MergeAABB(list.aabb, *v.bounding_box())
	}
	return &list
}

func (list *Hit_List) Add(hittable ...Hittable) {
	list.list = append(list.list, hittable...)
	for _, v := range hittable {
		list.aabb = *MergeAABB(list.aabb, *v.bounding_box())
	}
}

// See if the ray hits anything in the list of hittable things, and update record with the object closest to the camera.
func (list *Hit_List) hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) bool {
	var temp Hit
	closest_so_far := ray_tmax
	anything := false
	for _, object := range list.list {
		// We want the object closest to the camera
		if object.hit(ray, ray_tmin, closest_so_far, &temp) {
			anything = true
			closest_so_far = temp.t
			*record = temp
		}
	}

	return anything

}

func (list *Hit_List) bounding_box() (bounds *AABB) {
	return &list.aabb
}

type Translate struct {
	object *Hittable
	offset Vec3
	bbox   AABB
}

func NewTranslate(object Hittable, offset Vec3) *Translate {
	return &Translate{
		object: &object,
		offset: offset,
		bbox:   *object.bounding_box().AddOffset(offset),
	}
}

func (tran *Translate) hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) (ok bool) {
	// Move the ray backwards by the offset
	offset_ray := NewRay(*ray.origin.Sub(&tran.offset), ray.direction, ray.time)

	// Determine whether an intersection exists along the offset ray (and if so, where)

	if ok2 := (*tran.object).hit(&offset_ray, ray_tmin, ray_tmax, record); ok2 {
		record.point.IAdd(&tran.offset)
		return true
	}

	return false

}

func (tran *Translate) bounding_box() (bounds *AABB) {
	return &tran.bbox
}

type Rotate struct {
	object             *Hittable
	alpha, beta, gamma float64 // Angles in Radians for Yaw, Pitch and roll
	bbox               AABB
}

func NewRotate(object Hittable, alpha, beta, gamma float64) *Rotate {
	var rotate Rotate

	rotate.alpha = alpha * math.Pi / 180
	rotate.beta = beta * math.Pi / 180
	rotate.gamma = gamma * math.Pi / 180

	rotate.object = &object
	rotate.bbox = *object.bounding_box()

	min := NewVec3(math.Inf(1), math.Inf(1), math.Inf(1))
	max := NewVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1))

	// bbox calculations
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*rotate.bbox.maxVec[0] + (1-float64(i))*rotate.bbox.minVec[0]
				y := float64(j)*rotate.bbox.maxVec[1] + (1-float64(j))*rotate.bbox.minVec[1]
				z := float64(k)*rotate.bbox.maxVec[2] + (1-float64(k))*rotate.bbox.minVec[2]

				tester := NewVec3(x, y, z).RotateGen(rotate.alpha, rotate.beta, rotate.gamma)

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

	// origin[0] = rot.cos_theta*ray.origin[0] - rot.sin_theta*ray.origin[2]
	// origin[2] = rot.sin_theta*ray.origin[0] + rot.cos_theta*ray.origin[2]

	origin = *origin.RotateGen(-rot.alpha, -rot.beta, -rot.gamma)
	direction = *direction.RotateGen(-rot.alpha, -rot.beta, -rot.gamma)

	rotated := NewRay(origin, direction, ray.time)

	// Determine whether an intersection exists in object space (and if so, where)
	if !(*rot.object).hit(&rotated, ray_tmin, ray_tmax, record) {
		return false
	}

	// Change the intersection point from object space to world space
	point, normal := record.point, record.normal

	point = *point.RotateGen(-rot.alpha, -rot.beta, -rot.gamma)
	normal = *normal.RotateGen(-rot.alpha, -rot.beta, -rot.gamma)

	record.point = point
	record.normal = normal

	return true
}

func (rot *Rotate) bounding_box() (bounds *AABB) {
	return &rot.bbox
}
