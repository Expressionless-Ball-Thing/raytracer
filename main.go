package main

func main() {

	var world hittable_list
	world = append(world, &sphere{vec3{0, 0, -1}, 0.5})
	world = append(world, &sphere{vec3{0, -100.5, -1}, 100})

	camera := NewCamera(16.0/9.0, 5000)

	camera.render(world)
}
