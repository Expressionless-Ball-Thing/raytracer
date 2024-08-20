package main

type ray struct {
	origin    vec3
	direction vec3
}

func (r *ray) At(t float64) *vec3 {
	return r.origin.Add(*r.direction.Scale(t))
}
