package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

func bouncing_spheres() {

	// World
	var world Hit_List

	var checker Texture = NewCheckerFromColor(0.32, NewVec3(0.2, 0.3, 0.1), NewVec3(0.9, 0.9, 0.9))

	world.Add(
		NewSphere(NewVec3(0, -1000, 0), 1000, NewLambertTex(checker)),
	)

	for a := -5; a < 5; a++ {
		for b := -5; b < 5; b++ {
			choose_mat := rand.Float64()
			center := NewVec3(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Sub(NewVec3(4, 0.2, 0)).Magnitude() > 0.9 {
				if choose_mat < 0.8 {
					albedo := NewVec3Random(0, 1).Mult(NewVec3Random(0, 1))
					mat := NewLambert(albedo)
					center2 := center.Add(NewVec3(0, Random_float64_bounded(0, 0.5), 0))
					world.Add(NewMovingSphere(center, center2, 0.2, mat))
				} else if choose_mat < 0.95 {
					albedo := NewVec3Random(0.5, 1)
					fuzz := Random_float64_bounded(0, 0.5)
					mat := NewMetal(albedo, fuzz)
					world.Add(NewSphere(center, 0.2, mat))
				} else {
					mat := NewDielectric(1.5)
					world.Add(NewSphere(center, 0.2, mat))
				}
			}
		}
	}

	var mat1 Material = NewDielectric(1.5)
	world.Add(NewSphere(NewVec3(0, 1, 0), 1.0, mat1))

	var mat2 Material = NewLambert(NewVec3(0.4, 0.2, 0.1))
	world.Add(NewSphere(NewVec3(-4, 1, 0), 1.0, mat2))

	var mat3 Material = NewMetal(NewVec3(0.7, 0.6, 0.5), 0.0)
	world.Add(NewSphere(NewVec3(4, 1, 0), 1.0, mat3))

	var node Hittable = NewBVHNode(world.list)
	// world.list
	// world.aabb = *thing.bounding_box()

	cam := NewCamera(1000, NewVec3(13, 2, 3), NewVec3(0, 0, 0), NewVec3(0, 1, 0), 20, 16.0/9.0, 10.0, 0.6)
	cam.render(node, 10, 20)

}

func checked_spheres() {
	var world Hit_List

	checker := NewCheckerFromColor(0.32, NewVec3(0.2, 0.3, 0.1), NewVec3(0.9, 0.9, 0.9))

	world.Add(
		NewSphere(NewVec3(0, -10, 0), 10, NewLambertTex(checker)),
		NewSphere(NewVec3(0, 10, 0), 10, NewLambertTex(checker)),
	)

	cam := NewCamera(1000, NewVec3(13, 2, 3), NewVec3(0, 0, 0), NewVec3(0, 1, 0), 20, 16.0/9.0, 10.0, 0.6)
	cam.render(&world, 10, 20)

}

func earth() {
	file, _ := os.Open("earth.png")

	earth_texture, err := NewImageTexture(file)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	earth_surface := NewLambertTex(earth_texture)
	globe := NewSphere(NewVec3(0, 0, 0), 2, earth_surface)

	cam := NewCamera(1000, NewVec3(0, 0, 12), NewVec3(0, 0, 0), NewVec3(0, 1, 0), 20, 16.0/9.0, 10, 0)
	cam.render(NewList(globe), 10, 20)
}

func perlin_spheres() {
	var world Hit_List
	var pertext Texture = NewNoise(4)
	world.Add(
		NewSphere(NewVec3(0, -1000, 0), 1000, NewLambertTex(pertext)),
		NewSphere(NewVec3(0, 2, 0), 2, NewLambertTex(pertext)),
	)

	cam := NewCamera(1000, NewVec3(13, 2, 3), NewVec3(0, 0, 0), NewVec3(0, 1, 0), 20, 16.0/9.0, 1, 0)
	cam.render(&world, 10, 50)

}

func main() {
	start := time.Now()
	perlin_spheres()
	elapsed := time.Since(start)
	fmt.Println("Took", elapsed)
}
