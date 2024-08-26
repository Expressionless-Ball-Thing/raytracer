package main

import (
	"math"
	"math/rand/v2"
)

var (
	point_count int    = 256
	randVec3    []Vec3 = make([]Vec3, point_count)
	permX       []int  = make([]int, point_count)
	permY       []int  = make([]int, point_count)
	permZ       []int  = make([]int, point_count)
)

func perlin() {
	for i := 0; i < point_count; i++ {
		randVec3[i] = NewVec3Random(-1, 1).Unit()
	}

	perlin_generate_perm(permX)
	perlin_generate_perm(permY)
	perlin_generate_perm(permZ)

}

func noise(point Vec3) float64 {

	// Interpolate
	u := point[0] - math.Floor(point[0])
	v := point[1] - math.Floor(point[1])
	w := point[2] - math.Floor(point[2])
	// Use a Hermite cubic to round off the interpolation

	i := int(math.Floor(point[0]))
	j := int(math.Floor(point[1]))
	k := int(math.Floor(point[2]))

	var c [2][2][2]Vec3
	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = randVec3[permX[(i+di)&255]^permY[(j+dj)&255]^permZ[(k+dk)&255]]
			}
		}
	}

	return trillinear_interp(&c, u, v, w)
}

func perlin_generate_perm(array []int) {
	for i := 0; i < point_count; i++ {
		array[i] = i
	}

	permute(array, point_count)
}

func permute(array []int, n int) {
	for i := n - 1; i > 0; i-- {
		target := rand.IntN(i)
		tmp := array[i]
		array[i] = array[target]
		array[target] = tmp
	}
}

func trillinear_interp(c *[2][2][2]Vec3, u, v, w float64) float64 {
	uu := u * u * (3 - 2*u)
	vv := v * v * (3 - 2*v)
	ww := w * w * (3 - 2*w)
	accum := 0.0

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weight_v := NewVec3(u-float64(i), v-float64(j), w-float64(k))
				accum += (float64(i)*uu + (1-float64(i))*(1-uu)) * (float64(j)*vv + (1-float64(j))*(1-vv)) * (float64(k)*ww + (1-float64(k))*(1-ww)) * Dot(c[i][j][k], weight_v)
			}
		}
	}

	return accum
}

func turbulence(point Vec3, depth int) float64 {
	accum := 0.0
	temp_p := point
	weight := 1.0

	for i := 0; i < depth; i++ {
		accum += (weight * noise(temp_p))
		weight *= 0.5
		temp_p = temp_p.Scale(2)
	}

	return math.Abs(accum)
}
