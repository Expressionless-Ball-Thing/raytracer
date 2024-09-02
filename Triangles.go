package main

import "math"

type Triangle struct {
	Q        Vec3 // A corner on the Triangle
	u, v     Vec3 // Vector of the two sides extend from the corner Q
	w        Vec3 // Constant Value for calculating intersections
	material *Material
	normal   Vec3    // The normal
	D        float64 // Constant D
	bbox     AABB
}

// Create a new Quadrilaterial plane Given a point and two direction vectors
func NewTriangle(Q *Vec3, u, v *Vec3, material *Material) *Triangle {
	n := Cross(u, v)
	normal := n.Unit()
	return &Triangle{
		Q:        *Q,
		u:        *u,
		v:        *v,
		w:        *n.Scale(1 / (Dot(n, n))),
		material: material,
		normal:   *normal,
		D:        Dot(normal, Q),
		// Compute the bounding box of all four vertices.
		bbox: *MergeAABB(*NewAABB(*Q, *Q.Add(u).Add(v)), *NewAABB(*Q.Add(u), *Q.Add(v))),
	}
}

// If you do the path for a ray intersecting with a plane (tip, represent the plane in point normal form)
// You will find that the intersections t is equal to t = (D - n.P)/(n.d)
// where n is the normal, P and d are the origin point and direction of the Ray, D is n.v, where v is the point of intersection.
func (tri *Triangle) hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) (ok bool) {
	denom := Dot(&tri.normal, &ray.direction)

	// No hit if the ray is parallel to the plane.
	if math.Abs(denom) < 1e-8 {
		return false
	}

	// Return false if the hit point parameter t is outside the ray interval.
	t := (tri.D - Dot(&tri.normal, &ray.origin)) / denom
	if t < ray_tmin || t > ray_tmax {
		return false
	}

	// Determine if the hit point lies within the planar shape using its plane coordinates.
	intersection := ray.At(t)
	planar_hitpt_vector := intersection.Sub(&tri.Q) // P - Q basically
	alpha := Dot(&tri.w, Cross(planar_hitpt_vector, &tri.v))
	beta := Dot(&tri.w, Cross(&tri.u, planar_hitpt_vector))

	//return false if it is outside the primitive, otherwise set the hit record UV coordinates and return true.
	if !(alpha > 0 && beta > 0 && alpha+beta < 1) {
		return false
	}

	record.u = alpha
	record.v = beta
	record.t = t
	record.point = intersection
	record.material = tri.material
	record.set_face_normal(ray, tri.normal)
	return true
}

func (tri *Triangle) bounding_box() (bounds *AABB) {
	return &tri.bbox
}
