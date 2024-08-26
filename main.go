package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func main() {

	start := time.Now()
	// World
	var world Hit_List

	var material_ground Material = NewLambert(NewVec3(0.5, 0.5, 0.5))

	world.Add(
		NewSphere(NewVec3(0, -1000, -1), 1000, material_ground),
	)

	for a := -5; a < 5; a++ {
		for b := -5; b < 5; b++ {
			choose_mat := rand.Float64()
			center := NewVec3(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Sub(NewVec3(4, 0.2, 0)).Magnitude() > 0.9 {
				if choose_mat < 0.8 {
					albedo := NewVec3Random(0, 1).Mult(NewVec3Random(0, 1))
					mat := NewLambert(albedo)
					// center2 := center.Add(NewVec3(0, Random_float64_bounded(0, 0.5), 0))
					world.Add(NewSphere(center, 0.2, mat))
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

	cam := NewCamera(400, NewVec3(13, 0, 3), NewVec3(0, 0, 0), NewVec3(0, 1, 0), 20, 16.0/9.0, 10.0, 0.6)
	cam.render(node, 10, 20)

	elapsed := time.Since(start)

	fmt.Println("Took", elapsed)
}
