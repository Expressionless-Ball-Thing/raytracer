package main

import (
	"math"
)

type Quad struct {
	Q        Vec3 // A point on the Quad plane
	u, v     Vec3 // Direction vectors
	w        Vec3 // Constant Value for calculating intersections
	material *Material
	normal   Vec3    // The normal
	D        float64 // Constant D
	bbox     AABB
}

// Create a new Quadrilaterial plane Given a point and two direction vectors
func NewQuad(Q *Vec3, u, v *Vec3, material *Material) *Quad {
	n := Cross(u, v)
	normal := n.Unit()
	return &Quad{
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
func (quad *Quad) hit(ray *Ray, ray_tmin float64, ray_tmax float64, record *Hit) (ok bool) {
	denom := Dot(&quad.normal, &ray.direction)

	// No hit if the ray is parallel to the plane.
	if math.Abs(denom) < 1e-8 {
		return false
	}

	// Return false if the hit point parameter t is outside the ray interval.
	t := (quad.D - Dot(&quad.normal, &ray.origin)) / denom
	if t < ray_tmin || t > ray_tmax {
		return false
	}

	// Determine if the hit point lies within the planar shape using its plane coordinates.
	intersection := ray.At(t)
	planar_hitpt_vector := intersection.Sub(&quad.Q) // P - Q basically
	alpha := Dot(&quad.w, Cross(planar_hitpt_vector, &quad.v))
	beta := Dot(&quad.w, Cross(&quad.u, planar_hitpt_vector))

	//return false if it is outside the primitive, otherwise set the hit record UV coordinates and return true.
	if alpha < 0 || alpha > 1 || beta < 0 || beta > 1 {
		return false
	}
	record.u = alpha
	record.v = beta
	record.t = t
	record.point = intersection
	record.material = quad.material
	record.set_face_normal(ray, quad.normal)
	return true
}

func (quad *Quad) bounding_box() (bounds *AABB) {
	return &quad.bbox
}

// Makes 3D box (six sides) that contains the two opposite vertices a & b.
func NewBox(a, b Vec3, material *Material) *Hit_List {
	var sides Hit_List

	min := NewVec3(math.Min(a[0], b[0]), math.Min(a[1], b[1]), math.Min(a[2], b[2]))
	max := NewVec3(math.Max(a[0], b[0]), math.Max(a[1], b[1]), math.Max(a[2], b[2]))

	dx := NewVec3(max[0]-min[0], 0, 0)
	dy := NewVec3(0, max[1]-min[1], 0)
	dz := NewVec3(0, 0, max[2]-min[2])

	sides.Add(
		NewQuad(NewVec3(min[0], min[1], max[2]), dx, dy, material),          // front
		NewQuad(NewVec3(max[0], min[1], max[2]), dz.Negate(), dy, material), // right
		NewQuad(NewVec3(max[0], min[1], min[2]), dx.Negate(), dy, material), // back
		NewQuad(NewVec3(min[0], min[1], min[2]), dz, dy, material),          // left
		NewQuad(NewVec3(min[0], max[1], max[2]), dx, dz.Negate(), material), // top
		NewQuad(NewVec3(min[0], min[1], min[2]), dx, dz, material),          // bottom
	)

	return &sides
}
