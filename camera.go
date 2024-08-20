package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type camera struct {
	aspect_ratio  float64
	image_width   int
	image_height  int   // Rendered image height
	center        *vec3 // Camera center
	pixel00_loc   *vec3 // Location of pixel 0, 0
	pixel_delta_u *vec3 // Offset to pixel to the right
	pixel_delta_v *vec3 // Offset to pixel below
}

// Makes a new camera given the aspect ratio and image width
func NewCamera(aspect_ratio float64, image_width int) *camera {
	var cam camera
	cam.aspect_ratio = aspect_ratio
	cam.image_width = image_width
	return &cam
}

// Render the scene
func (cam *camera) render(world hittable) {

	cam.initalize()

	upLeft := image.Point{0, 0}
	lowRight := image.Point{cam.image_width, cam.image_height}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	for j := 0; j < cam.image_height; j++ {
		for i := 0; i < cam.image_width; i++ {
			pixel_center := cam.pixel00_loc.Add(*cam.pixel_delta_u.Scale(float64(i)).Add(*cam.pixel_delta_v.Scale(float64(j))))
			ray_direction := pixel_center.Sub(*cam.center)
			ray := ray{*cam.center, *ray_direction}

			pixel_color := ray_color(&ray, world)

			img.Set(i, j, color.NRGBA{uint8(255 * pixel_color.X), uint8(255 * pixel_color.Y), uint8(255 * pixel_color.Z), 0xff})
		}
	}

	file, err := os.Create("image2.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	png.Encode(file, img)
}

// Set up the camera
func (cam *camera) initalize() {

	// Calculate the image height, and ensure that it's at least 1.
	cam.image_height = int(float64(cam.image_width) / float64(cam.aspect_ratio))
	if cam.image_height < 1 {
		cam.image_height = 1
	}

	cam.center = &vec3{0, 0, 0}

	// Determine viewport dimensions.
	focal_length := 1.0
	viewport_height := 2.0
	viewport_width := viewport_height * (float64(cam.image_width) / float64(cam.image_height))

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewport_u := vec3{viewport_width, 0, 0}
	viewport_v := vec3{0, -viewport_height, 0}

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	cam.pixel_delta_u = viewport_u.Scale(1.0 / float64(cam.image_width))
	cam.pixel_delta_v = viewport_v.Scale(1.0 / float64(cam.image_height))

	// Calculate the location of the upper left pixel.
	viewport_upper_left := cam.center.Sub(vec3{0, 0, focal_length}).Sub(*viewport_u.Scale(0.5)).Sub(*viewport_v.Scale(0.5))
	// reminder that the edge of the viewport is offset by half a pixel cell.
	cam.pixel00_loc = viewport_upper_left.Add(*(*cam.pixel_delta_u.Add(*cam.pixel_delta_v)).Scale(0.5))

}

func ray_color(ray *ray, world hittable) vec3 {

	var rec hit_record
	if world.hit(ray, 0, math.Inf(1), &rec) {
		return *(rec.normal.Add(vec3{1, 1, 1})).Scale(0.5)
	}

	unit_direction := Unit(&ray.direction)
	a := (unit_direction.Y + 1.0) * 0.5
	return *(&vec3{1.0, 1.0, 1.0}).Scale(1 - a).Add(*(&vec3{0.5, 0.7, 1.0}).Scale(a))
}
