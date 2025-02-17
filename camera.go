package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
)

type camera struct {
	aspect_ratio                   float64
	image_width                    int
	image_height                   int     // Rendered image height
	camera_center                  Vec3    // Camera center
	pixel00_loc                    Vec3    // Location of pixel 0, 0
	pixel_delta_u                  Vec3    // Offset to pixel to the right
	pixel_delta_v                  Vec3    // Offset to pixel below
	sample_per_pixel               int     // Count of random samples for each pixel
	pixel_samples_scale            float64 // Color scale factor for a sum of pixel samples
	max_depth                      int     // Maximum number of ray bounces into scene
	background                     Vec3    // Scene background color
	vfov                           float64 // Vertical view angle (field of view) in degrees
	lookfrom                       Vec3    // Point camera is looking from
	lookat                         Vec3    // Point camera is looking at
	vup                            Vec3    // Camera-relative "up" direction
	u, v, w                        Vec3    // Camera frame basis vectors (u is camera right, v is camera up, w is opposite of view direction)
	defocus_angle                  float64 // Variation angle of rays through each pixel
	focus_distance                 float64 // Distance from camera lookfrom point to plane of perfect focus
	defocus_disk_u, defocus_disk_v Vec3    // Defocus disk horizontal/vertical radius
}

// Makes a new camera given the aspect ratio and image width

// image_width int, lookFrom, lookAt, vup Vec3, vfov, aspect_ratio, focus_distance, defocus_angle float64
func NewCamera(image_width int, lookFrom, lookAt, vup Vec3, vfov, aspect_ratio, focus_distance, defocus_angle float64, background Vec3) (cam *camera) {

	var camera camera

	camera.lookfrom = lookFrom
	camera.lookat = lookAt
	camera.vup = vup
	camera.aspect_ratio = 16.0 / 9.0
	camera.image_width = image_width
	camera.vfov = vfov
	camera.camera_center = lookFrom
	camera.defocus_angle = defocus_angle
	camera.focus_distance = focus_distance
	camera.background = background

	// Calculate the image height, and ensure that it's at least 1.
	camera.image_height = int(float64(image_width) / float64(aspect_ratio))
	if camera.image_height < 1 {
		camera.image_height = 1
	}

	// Viewport Dimensions
	theta := vfov * math.Pi / 180
	h := math.Tan(theta / 2)
	viewport_height := 2.0 * h * camera.focus_distance
	viewport_width := viewport_height * (float64(camera.image_width) / float64(camera.image_height))

	// Calculate the u,v,w unit basis vectors for the camera coordinate frame.
	camera.w = *lookFrom.Sub(&lookAt).Unit()
	camera.u = *Cross(&vup, &camera.w).Unit()
	camera.v = *Cross(&camera.w, &camera.u)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewport_u := camera.u.Scale(viewport_width)
	viewport_v := camera.v.Scale(-viewport_height)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	camera.pixel_delta_u = *viewport_u.Scale(1.0 / float64(camera.image_width))
	camera.pixel_delta_v = *viewport_v.Scale(1.0 / float64(camera.image_height))

	// Calculate the location of the upper left pixel.
	viewport_upper_left := camera.camera_center.Sub(camera.w.Scale(camera.focus_distance)).Sub(viewport_u.Scale(0.5)).Sub(viewport_v.Scale(0.5))
	camera.pixel00_loc = *viewport_upper_left.Add((camera.pixel_delta_u.Add(&camera.pixel_delta_v)).Scale(0.5))

	return &camera
}

type result struct {
	rownum     int
	pixel_list []Vec3
}

// Render the scene
// world Hittable, sample_per_pixel, max_depth int
// Parallized row by row as well using a worker pool, model is based on this https://gobyexample.com/worker-pools
func (cam *camera) render(world Hittable, sample_per_pixel, max_depth int) {
	cam.sample_per_pixel = sample_per_pixel
	cam.max_depth = max_depth
	cam.pixel_samples_scale = 1.0 / float64(cam.sample_per_pixel)

	jobs := make(chan int, cam.image_height)         // Job channel, indicates the row number
	results := make(chan result, cam.image_height*2) // Results channel

	upLeft := image.Point{0, 0}
	lowRight := image.Point{cam.image_width, cam.image_height}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	// The worker function.
	worker := func() {
		for row_num := range jobs {
			pixelRow := make([]Vec3, cam.image_width)
			for col_num := 0; col_num < cam.image_width; col_num++ {
				pixel_color := NewVec3(0, 0, 0)
				// Loop for antialiasing
				for sample := 0; sample < cam.sample_per_pixel; sample++ {
					ray := cam.get_ray(float64(col_num), float64(row_num))
					pixel_color.IAdd((*cam).ray_color(ray, cam.max_depth, world))
				}

				pixelRow[col_num] = *pixel_color.Scale(cam.pixel_samples_scale).Gamma(2)
			}
			results <- result{rownum: row_num, pixel_list: pixelRow}
		}
	}

	// number of workers in the pool
	worker_count := runtime.NumCPU()

	// Spin up the worker threads
	for i := 0; i < worker_count; i++ {
		go worker()
	}

	// Distribute work
	for i := 0; i < cam.image_height; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 0; i < cam.image_height; i++ {
		result := <-results
		for idx, pixel_color := range result.pixel_list {
			img.Set(idx, result.rownum, color.NRGBA{
				uint8(255 * clamp(pixel_color[0])),
				uint8(255 * clamp(pixel_color[1])),
				uint8(255 * clamp(pixel_color[2])),
				0xff})
		}
	}
	close(results)

	file, err := os.Create("main.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	png.Encode(file, img)
}

// Construct a camera ray originating from the defocus disk and directed at a randomly
// sampled point around the pixel location i, j.
func (cam *camera) get_ray(i float64, j float64) Ray {
	offset := sample_square()
	pixel_sample := cam.pixel00_loc.Add(cam.pixel_delta_u.Scale(i + offset[0])).Add(cam.pixel_delta_v.Scale(j + offset[1]))
	ray_origin := cam.camera_center
	if cam.defocus_angle > 0 {
		ray_origin = cam.defocus_disk_sample()
	}
	ray_direction := pixel_sample.Sub(&ray_origin)
	ray_time := rand.Float64()

	return Ray{ray_origin, *ray_direction, ray_time}
}

func (cam *camera) defocus_disk_sample() Vec3 {
	p := random_in_unit_disk()
	t := cam.camera_center.Add(cam.defocus_disk_u.Scale(p[0]))
	return *t.Add(cam.defocus_disk_v.Scale(p[1]))
}

func (camera *camera) ray_color(ray Ray, depth int, world Hittable) *Vec3 {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return NewVec3(0, 0, 0)
	}

	var rec Hit
	if !world.hit(&ray, 0.001, math.MaxFloat64, &rec) {
		return &camera.background
	}

	var attenuation Vec3
	var scattered Ray
	color_from_emission := (*rec.material).emitted(rec.u, rec.v, &rec.point)

	if !(*rec.material).scatter(&ray, &rec, &attenuation, &scattered) {
		return &color_from_emission
	}
	color_from_scatter := (camera.ray_color(scattered, depth-1, world)).Mult(&attenuation)
	return color_from_emission.Add(color_from_scatter)
}

func sample_square() Vec3 {
	return Vec3{rand.Float64() - 0.5, rand.Float64() - 0.5, 0}
}

// Clamp number between 0 and 1
func clamp(num float64) float64 {
	if num < 0.0 {
		return 0.0
	} else if num > 1 {
		return 1
	}
	return num
}
