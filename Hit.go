package main

type Hittable interface {

	// Calculates whether a hit can be made with the object within the given bounds.
	// Returns a nil pointer if no hits can be made
	hit(ray *Ray, ray_tmin float64, ray_tmax float64) (record Hit, ok bool)

	bounding_box() (bounds *AABB)

	count() int
}

// Struct to store the details of a ray hitting a surface
type Hit struct {
	point      Vec3
	normal     Vec3
	t          float64
	front_face bool // Hack way to check front_face or not Dot(&in, &n) < 0
	material   Material
}

// Sets the hit record normal vector.
// NOTE: the parameter `outward_normal` is assumed to have unit length.
func (record *Hit) set_face_normal(ray *Ray, outward_normal Vec3) {
	record.front_face = Dot(ray.direction, outward_normal) < 0
	if record.front_face {
		record.normal = outward_normal
	} else {
		record.normal = (outward_normal.Negate())
	}
}

type Hit_List struct {
	list []Hittable
	aabb AABB
}

func NewList(objects ...Hittable) *Hit_List {
	var list Hit_List
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
func (list *Hit_List) hit(ray *Ray, ray_tmin float64, ray_tmax float64) (record Hit, ok bool) {
	closest_so_far := ray_tmax
	anything := false

	for _, object := range list.list {
		// We want the object closest to the camera
		if temp, v := object.hit(ray, ray_tmin, closest_so_far); v {
			anything = true
			closest_so_far = temp.t
			record = temp
		}
	}

	return record, anything

}

func (list *Hit_List) count() int {
	return len(list.list)
}

func (list *Hit_List) bounding_box() (bounds *AABB) {
	return &list.aabb
}
