package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/pkg/profile"
)

func bouncing_spheres() {

	// World
	var world Hit_List

	checker := NewCheckerFromColor(0.32, *NewVec3(0.2, 0.3, 0.1), *NewVec3(0.9, 0.9, 0.9))

	world.Add(
		NewSphere(*NewVec3(0, -1000, 0), 1000, NewLambertTex(checker)),
	)

	for a := -5; a < 5; a++ {
		for b := -5; b < 5; b++ {
			choose_mat := rand.Float64()
			center := NewVec3(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Sub(NewVec3(4, 0.2, 0)).Magnitude() > 0.9 {
				if choose_mat < 0.8 {
					albedo := NewVec3Random(0, 1).Mult(NewVec3Random(0, 1))
					mat := NewLambert(*albedo)
					// center2 := center.Add(NewVec3(0, Random_float64_bounded(0, 0.5), 0))
					world.Add(NewSphere(*center, 0.2, mat))
				} else if choose_mat < 0.95 {
					albedo := NewVec3Random(0.5, 1)
					fuzz := Random_float64_bounded(0, 0.5)
					mat := NewMetal(*albedo, fuzz)
					world.Add(NewSphere(*center, 0.2, mat))
				} else {
					mat := NewDielectric(1.5)
					world.Add(NewSphere(*center, 0.2, mat))
				}
			}
		}
	}

	mat1 := NewDielectric(1.5)
	world.Add(NewSphere(*NewVec3(0, 1, 0), 1.0, mat1))

	mat2 := NewLambert(*NewVec3(0.4, 0.2, 0.1))
	world.Add(NewSphere(*NewVec3(-4, 1, 0), 1.0, mat2))

	mat3 := NewMetal(*NewVec3(0.7, 0.6, 0.5), 0.0)
	world.Add(NewSphere(*NewVec3(4, 1, 0), 1.0, mat3))

	cam := NewCamera(1200, *NewVec3(13, 2, 3), *NewVec3(0, 0, 0), *NewVec3(0, 1, 0), 20, 16.0/9.0, 10.0, 0.6, *NewVec3(0.7, 0.8, 1.0))
	cam.render(&world, 100, 10)

}

func checkered_spheres() {
	var world Hit_List

	checker := NewCheckerFromColor(0.32, *NewVec3(0.2, 0.3, 0.1), *NewVec3(0.9, 0.9, 0.9))

	world.Add(
		NewSphere(*NewVec3(0, -10, 0), 10, NewLambertTex(checker)),
		NewSphere(*NewVec3(0, 10, 0), 10, NewLambertTex(checker)),
	)

	cam := NewCamera(1000, *NewVec3(13, 2, 3), *NewVec3(0, 0, 0), *NewVec3(0, 1, 0), 20, 16.0/9.0, 10.0, 0.6, *NewVec3(0.7, 0.8, 1.0))
	cam.render(&world, 10, 20)

}

func earth() {
	file, _ := os.Open("earth.png")

	earth_texture := NewImageTexture(file)
	earth_surface := NewLambertTex(earth_texture)
	globe := NewSphere(*NewVec3(0, 0, 0), 2, earth_surface)

	cam := NewCamera(1000, *NewVec3(0, 0, 12), *NewVec3(0, 0, 0), *NewVec3(0, 1, 0), 20, 16.0/9.0, 10, 0, *NewVec3(0.7, 0.8, 1.0))
	cam.render(NewList(globe), 10, 20)
}

func perlin_spheres() {
	var world Hit_List
	pertext := NewNoise(4)
	world.Add(
		NewSphere(*NewVec3(0, -1000, 0), 1000, NewLambertTex(pertext)),
		NewSphere(*NewVec3(0, 2, 0), 2, NewLambertTex(pertext)),
	)

	cam := NewCamera(1000, *NewVec3(13, 2, 3), *NewVec3(0, 0, 0), *NewVec3(0, 1, 0), 20, 16.0/9.0, 1, 0, *NewVec3(0.7, 0.8, 1.0))
	cam.render(&world, 10, 50)

}

func quads() {
	var world Hit_List

	// Materials
	left_red := NewLambert(*NewVec3(1.0, 0.2, 0.2))
	back_green := NewLambert(*NewVec3(0.2, 1.0, 0.2))
	right_blue := NewLambert(*NewVec3(0.2, 0.2, 1.0))
	upper_orange := NewLambert(*NewVec3(1.0, 0.5, 0.0))
	lower_teal := NewLambert(*NewVec3(0.2, 0.8, 0.8))

	//Quads
	world.Add(
		NewQuad(NewVec3(-3, -2, 5), NewVec3(0, 0, -4), NewVec3(0, 4, 0), left_red),
		NewQuad(NewVec3(-2, -2, 0), NewVec3(4, 0, 0), NewVec3(0, 4, 0), back_green),
		NewQuad(NewVec3(3, -2, 1), NewVec3(0, 0, 4), NewVec3(0, 4, 0), right_blue),
		NewQuad(NewVec3(-2, 3, 1), NewVec3(4, 0, 0), NewVec3(0, 0, 4), upper_orange),
		NewQuad(NewVec3(-2, -3, 5), NewVec3(4, 0, 0), NewVec3(0, 0, -4), lower_teal),
	)

	cam := NewCamera(1000, *NewVec3(0, 0, 9), *NewVec3(0, 0, 0), *NewVec3(0, 1, 0), 80, 16.0/9.0, 1, 0, *NewVec3(0.7, 0.8, 1.0))
	cam.render(&world, 100, 50)
}

func simple_light() {
	var world Hit_List
	pertext := NewNoise(4)
	world.Add(
		NewSphere(*NewVec3(0, -1000, 0), 1000, NewLambertTex(pertext)),
		NewSphere(*NewVec3(0, 2, 0), 2, NewLambertTex(pertext)),
	)

	difflight := NewDiffuseLightColor(*NewVec3(4, 4, 4))
	world.Add(
		NewQuad(NewVec3(3, 1, -2), NewVec3(2, 0, 0), NewVec3(0, 2, 0), difflight),
		NewSphere(*NewVec3(0, 7, 0), 2, difflight),
	)

	cam := NewCamera(1000, *NewVec3(26, 3, 6), *NewVec3(0, 2, 0), *NewVec3(0, 1, 0), 20, 16.0/9.0, 1, 0, *NewVec3(0, 0, 0))
	cam.render(&world, 100, 50)

}

func cornell_box() {
	red := NewLambert(*NewVec3(.65, .05, .05))
	white := NewLambert(*NewVec3(.73, .73, .73))
	green := NewLambert(*NewVec3(.12, .45, .15))
	light := NewDiffuseLightColor(*NewVec3(15, 15, 15))

	var world Hit_List

	world.Add(NewQuad(NewVec3(555, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), green))
	world.Add(NewQuad(NewVec3(0, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), red))
	world.Add(NewQuad(NewVec3(343, 554, 332), NewVec3(-130, 0, 0), NewVec3(0, 0, -105), light))
	world.Add(NewQuad(NewVec3(0, 0, 0), NewVec3(555, 0, 0), NewVec3(0, 0, 555), white))
	world.Add(NewQuad(NewVec3(555, 555, 555), NewVec3(-555, 0, 0), NewVec3(0, 0, -555), white))
	world.Add(NewQuad(NewVec3(0, 0, 555), NewVec3(555, 0, 0), NewVec3(0, 555, 0), white))
	world.Add(NewTranslate(NewRotate(NewBox(*NewVec3(0, 0, 0), *NewVec3(165, 330, 165), white), 0, 15, 0), *NewVec3(265, 0, 295)))
	world.Add(NewTranslate(NewRotate(NewBox(*NewVec3(0, 0, 0), *NewVec3(165, 165, 165), white), 0, -18, 0), *NewVec3(130, 0, 65)))

	cam := NewCamera(600, *NewVec3(278, 278, -800), *NewVec3(278, 278, 0), *NewVec3(0, 1, 0), 40, 1, 1, 0, *NewVec3(0, 0, 0))
	cam.render(&world, 200, 50)
}

func cornell_smoke() {
	red := NewLambert(*NewVec3(.65, .05, .05))
	white := NewLambert(*NewVec3(.73, .73, .73))
	green := NewLambert(*NewVec3(.12, .45, .15))
	light := NewDiffuseLightColor(*NewVec3(7, 7, 7))

	var world Hit_List

	world.Add(NewQuad(NewVec3(555, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), green))
	world.Add(NewQuad(NewVec3(0, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), red))
	world.Add(NewQuad(NewVec3(113, 554, 127), NewVec3(330, 0, 0), NewVec3(0, 0, 305), light))
	world.Add(NewQuad(NewVec3(0, 555, 0), NewVec3(555, 0, 0), NewVec3(0, 0, 555), white))
	world.Add(NewQuad(NewVec3(0, 0, 0), NewVec3(555, 0, 0), NewVec3(0, 0, 555), white))
	world.Add(NewQuad(NewVec3(0, 0, 555), NewVec3(555, 0, 0), NewVec3(0, 555, 0), white))

	var box1 Hittable = NewBox(*NewVec3(0, 0, 0), *NewVec3(165, 330, 165), white)
	box1 = NewRotate(box1, 0, 15, 0)
	box1 = NewTranslate(box1, *NewVec3(265, 0, 295))
	world.Add(NewConstantMediumAlbedo(&box1, 0.01, *NewVec3(0, 0, 0)))
	world.Add(box1)

	var box2 Hittable = NewBox(*NewVec3(0, 0, 0), *NewVec3(165, 165, 165), white)
	box2 = NewRotate(box2, 0, -18, 0)
	box2 = NewTranslate(box2, *NewVec3(130, 0, 65))
	world.Add(NewConstantMediumAlbedo(&box2, 0.01, *NewVec3(1, 1, 1)))
	world.Add(box2)

	cam := NewCamera(600, *NewVec3(278, 278, -800), *NewVec3(278, 278, 0), *NewVec3(0, 1, 0), 40, 1, 1, 0, *NewVec3(0, 0, 0))
	cam.render(&world, 200, 50)
}

func final_scene(image_width, samples_per_pixel, max_depth int) {
	var boxes1 Hit_List
	ground := NewLambert(*NewVec3(0.48, 0.83, 0.53))

	boxes_per_side := 20
	for i := 0; i < boxes_per_side; i++ {
		for j := 0; j < boxes_per_side; j++ {
			w := 100.0
			x0 := -1000.0 + float64(i)*w
			z0 := -1000.0 + float64(j)*w
			y0 := 0.0
			x1 := x0 + w
			y1 := Random_float64_bounded(1, 101)
			z1 := z0 + w

			boxes1.Add(NewBox(*NewVec3(x0, y0, z0), *NewVec3(x1, y1, z1), ground))
		}
	}

	var world Hit_List

	world.Add(NewBVHNode(boxes1.list))

	light := NewDiffuseLightColor(*NewVec3(7, 7, 7))
	world.Add(NewQuad(NewVec3(123, 554, 147), NewVec3(300, 0, 0), NewVec3(0, 0, 265), light))

	center1 := NewVec3(400, 400, 200)
	center2 := center1.Add(NewVec3(30, 0, 0))
	sphere_mat := NewLambert(*NewVec3(0.7, 0.3, 0.1))
	world.Add(NewMovingSphere(*center1, *center2, 50, sphere_mat))

	world.Add(NewSphere(*NewVec3(260, 150, 45), 50, NewDielectric(1.5)))
	world.Add(NewSphere(*NewVec3(0, 150, 145), 50, NewMetal(*NewVec3(0.8, 0.8, 0.9), 1.0)))

	var boundary Hittable = NewSphere(*NewVec3(360, 150, 145), 70, NewDielectric(1.5))
	world.Add(boundary)
	world.Add(NewConstantMediumAlbedo(&boundary, 0.2, Vec3{0.2, 0.4, 0.9}))
	boundary = NewSphere(*NewVec3(0, 0, 0), 5000, NewDielectric(1.5))
	world.Add(NewConstantMediumAlbedo(&boundary, 0.0001, *NewVec3(1, 1, 1)))

	f, _ := os.Open("earth.png")
	emat := NewLambertTex(NewImageTexture(f))
	world.Add(NewSphere(*NewVec3(400, 200, 400), 100, emat))
	pertext := NewNoise(0.2)
	world.Add(NewSphere(*NewVec3(220, 280, 300), 80, NewLambertTex(pertext)))

	var boxes2 Hit_List
	white := NewLambert(*NewVec3(.73, .73, .73))
	ns := 1000
	for j := 0; j < ns; j++ {
		boxes2.Add(NewSphere(*NewVec3Random(0, 165), 10, white))
	}

	world.Add(NewTranslate(
		NewRotate(
			NewBVHNode(boxes2.list), 0, 15, 0), *NewVec3(-100, 270, 395),
	),
	)

	cam := NewCamera(
		image_width,
		*NewVec3(478, 278, -600),
		*NewVec3(278, 278, 0),
		*NewVec3(0, 1, 0),
		40,
		1.0,
		1,
		0,
		*NewVec3(0, 0, 0))
	cam.render(&world, samples_per_pixel, max_depth)

}

func main() {
	wd, _ := os.Getwd()
	defer profile.Start(profile.ProfilePath(wd)).Stop()
	start := time.Now()
	switch 7 {
	case 1:
		bouncing_spheres()
	case 2:
		checkered_spheres()
	case 3:
		earth()
	case 4:
		perlin_spheres()
	case 5:
		quads()
	case 6:
		simple_light()
	case 7:
		cornell_box()
	case 8:
		cornell_smoke()
	case 9:
		final_scene(800, 10000, 40)
	default:
		final_scene(400, 250, 4)
	}
	elapsed := time.Since(start)
	fmt.Println("Took", elapsed)
}
