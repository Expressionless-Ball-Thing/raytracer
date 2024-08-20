package main

type hit_record struct {
	point      vec3
	normal     vec3
	t          float64
	front_face bool
}

// Sets the hit record normal vector.
// NOTE: the parameter `outward_normal` is assumed to have unit length.
func (record *hit_record) set_face_normal(ray *ray, outward_normal *vec3) {
	record.front_face = Dot(&ray.direction, outward_normal) < 0
	if record.front_face {
		record.normal = *outward_normal
	} else {
		record.normal = *(outward_normal.Negate())
	}
}

type hittable interface {

	// Calculates whether a hit can be made with the object within the given bounds.
	hit(ray *ray, ray_tmin float64, ray_tmax float64, record *hit_record) bool
}

type hittable_list []hittable

// See if the ray hits anything in the list of hittable things, and update record with the object closest to the camera.
func (list hittable_list) hit(ray *ray, ray_tmin float64, ray_tmax float64, record *hit_record) bool {
	var temp hit_record
	var hit_bool bool
	closest_so_far := ray_tmax

	for _, object := range list {
		// We want the object closest to the camera
		if object.hit(ray, ray_tmin, closest_so_far, &temp) {
			hit_bool = true
			closest_so_far = temp.t
			*record = temp
		}
	}
	return hit_bool

}
