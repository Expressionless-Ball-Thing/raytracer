package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Color Takes values forom 0 to 1 for each value in RGB

func ray_color(ray *ray) vec3 {

	unit_direction := Unit(&ray.direction)
	a := (unit_direction.Y + 1.0) * 0.5
	return *(&vec3{1.0, 1.0, 1.0}).Scale(1 - a).Add(*(&vec3{0.1, 0.7, 1.0}).Scale(a))
}

func main() {

	// Image
	aspect_ratio := 16.0 / 9.0
	image_width := 400

	// Calculate the image height, and ensure that it's at least 1.
	image_height := int(float64(image_width) / float64(aspect_ratio))
	if image_height < 1 {
		image_height = 1
	}

	// Camera
	focal_length := 1.0
	viewport_height := 2.0
	viewport_width := viewport_height * (float64(image_width) / float64(image_height))
	camera_center := vec3{0, 0, 0}

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewport_u := vec3{viewport_width, 0, 0}
	viewport_v := vec3{0, -viewport_height, 0}

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	pixel_delta_u := viewport_u.Scale(1.0 / float64(image_width))
	pixel_delta_v := viewport_v.Scale(1.0 / float64(image_height))

	// Calculate the location of the upper left pixel.
	viewport_upper_left := camera_center.Sub(vec3{0, 0, focal_length}).Sub(*viewport_u.Scale(0.5)).Sub(*viewport_v.Scale(0.5))
	// reminder that the edge of the viewport is offset by half a pixel cell.
	pixel00_loc := viewport_upper_left.Add(*(*pixel_delta_u.Add(*pixel_delta_v)).Scale(0.5))

	// Render

	upLeft := image.Point{0, 0}
	lowRight := image.Point{image_width, image_height}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	for j := 0; j < image_height; j++ {
		for i := 0; i < image_width; i++ {
			pixel_center := pixel00_loc.Add(*pixel_delta_u.Scale(float64(i)).Add(*pixel_delta_v.Scale(float64(j))))
			ray_direction := pixel_center.Sub(camera_center)
			ray := ray{camera_center, *ray_direction}

			pixel_color := ray_color(&ray)

			img.Set(i, j, color.NRGBA{uint8(255 * pixel_color.X), uint8(255 * pixel_color.Y), uint8(255 * pixel_color.Z), 0xff})
		}
	}

	file, err := os.Create("image.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	png.Encode(file, img)

}
