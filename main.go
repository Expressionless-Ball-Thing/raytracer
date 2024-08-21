package main

import "math/rand/v2"

func main() {

	var world hittable_list
	var ground_mat material = &lambertian{vec3{0.5, 0.5, 0.5}}
	world = append(world, &sphere{&vec3{0, -1000, 0}, 1000, &ground_mat})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			choose_mat := rand.Float64()
			center := vec3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if (center.Sub(vec3{4, 0.2, 0}).Magnitude() > 0.9) {
				var sphere_mat material

				if choose_mat < 0.8 {
					//diffuse
					albedo := Random_unit_vec3().Mult(*Random_unit_vec3())
					sphere_mat = &lambertian{*albedo}
					world = append(world, NewSphere(&center, 0.2, &sphere_mat))
				} else if choose_mat < 0.95 {
					//metal
					albedo := RandomVec3(0.5, 1)
					fuzz := random_float64_bounded(0, 0.5)
					sphere_mat = &metal{*albedo, fuzz}
					world = append(world, NewSphere(&center, 0.2, &sphere_mat))
				} else {
					// glass
					sphere_mat = &dielectric{1.5}
					world = append(world, NewSphere(&center, 0.2, &sphere_mat))
				}
			}
		}
	}

	var material1 material = &dielectric{1.5}
	world = append(world, NewSphere(&vec3{0, 1, 0}, 1.0, &material1))

	var material2 material = &lambertian{vec3{0.4, 0.2, 0.1}}
	world = append(world, NewSphere(&vec3{-4, 1, 0}, 1.0, &material2))

	var material3 material = &metal{vec3{0.7, 0.6, 0.5}, 0.0}
	world = append(world, NewSphere(&vec3{4, 1, 0}, 1.0, &material3))

	camera := NewCamera(16.0/9.0, 1200)
	camera.sample_per_pixel = 100
	camera.max_depth = 10

	camera.vfov = 20
	camera.lookfrom = vec3{13, 2, 3}
	camera.lookat = vec3{0, 0, 0}
	camera.vup = vec3{0, 1, 0}

	camera.defocus_angle = 0.6
	camera.focus_distance = 10

	camera.render(world)
}
