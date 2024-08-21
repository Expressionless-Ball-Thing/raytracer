package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand/v2"
	"os"
)

type camera struct {
	aspect_ratio                   float64
	image_width                    int
	image_height                   int     // Rendered image height
	center                         *vec3   // Camera center
	pixel00_loc                    *vec3   // Location of pixel 0, 0
	pixel_delta_u                  *vec3   // Offset to pixel to the right
	pixel_delta_v                  *vec3   // Offset to pixel below
	sample_per_pixel               int     // Count of random samples for each pixel
	pixel_samples_scale            float64 // Color scale factor for a sum of pixel samples
	max_depth                      int     // Maximum number of ray bounces into scene
	vfov                           float64 // Vertical view angle (field of view) in degrees
	lookfrom                       vec3    // Point camera is looking from
	lookat                         vec3    // Point camera is looking at
	vup                            vec3    // Camera-relative "up" direction
	u, v, w                        vec3    // Camera frame basis vectors (u is camera right, v is camera up, w is opposite of view direction)
	defocus_angle                  float64 // Variation angle of rays through each pixel
	focus_distance                 float64 // Distance from camera lookfrom point to plane of perfect focus
	defocus_disk_u, defocus_disk_v vec3    // Defocus disk horizontal/vertical radius
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

			pixel_color := vec3{0, 0, 0}
			// Loop for antialiasing
			for sample := 0; sample < cam.sample_per_pixel; sample++ {
				ray := cam.get_ray(float64(i), float64(j))
				pixel_color.IAdd(*ray_color(&ray, cam.max_depth, world))
			}
			pixel_color.IScale(cam.pixel_samples_scale)

			// The original way of doing things.
			// pixel_center := cam.pixel00_loc.Add(*cam.pixel_delta_u.Scale(float64(i)).Add(*cam.pixel_delta_v.Scale(float64(j))))
			// ray_direction := pixel_center.Sub(*cam.center)
			// ray := ray{*cam.center, *ray_direction}
			// pixel_color := ray_color(&ray, world)

			// Gamma correction by taking the square root of it.
			pixel_color.X = linear_to_gamma(pixel_color.X)
			pixel_color.Y = linear_to_gamma(pixel_color.Y)
			pixel_color.Z = linear_to_gamma(pixel_color.Z)

			img.Set(i, j, color.NRGBA{
				uint8(255 * clamp(pixel_color.X)),
				uint8(255 * clamp(pixel_color.Y)),
				uint8(255 * clamp(pixel_color.Z)),
				0xff})
		}
	}

	file, err := os.Create("reflectanceTest.png")
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

	cam.pixel_samples_scale = 1.0 / float64(cam.sample_per_pixel)

	cam.center = &cam.lookfrom

	// Determine viewport dimensions.
	theta := cam.vfov * math.Pi / 180
	h := math.Tan(theta / 2) // We want viewport to the centered at the camera.
	viewport_height := 2 * h * cam.focus_distance
	viewport_width := viewport_height * (float64(cam.image_width) / float64(cam.image_height))

	// Calculate the u,v,w unit basis vectors for the camera coordinate frame.
	cam.w = *Unit(cam.lookfrom.Sub(cam.lookat))
	cam.u = *Unit(Cross(&cam.vup, &cam.w))
	cam.v = *Cross(&cam.w, &cam.u)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewport_u := (&cam.u).Scale(viewport_width)   // Vector across viewport horizontal edge
	viewport_v := (&cam.v).Scale(-viewport_height) // Vector down viewport vertical edge

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	cam.pixel_delta_u = viewport_u.Scale(1.0 / float64(cam.image_width))
	cam.pixel_delta_v = viewport_v.Scale(1.0 / float64(cam.image_height))

	// Calculate the location of the upper left pixel.
	viewport_upper_left := (cam.center).Sub(*cam.w.Scale(cam.focus_distance)).Sub(*viewport_u.Scale(0.5)).Sub(*viewport_v.Scale(0.5))
	// reminder that the edge of the viewport is offset by half a pixel cell.
	cam.pixel00_loc = viewport_upper_left.Add(*(*cam.pixel_delta_u.Add(*cam.pixel_delta_v)).Scale(0.5))

	// Calculate the camera defocus disk basis vectors.
	defocus_radius := cam.focus_distance * math.Tan(cam.defocus_angle*math.Pi/180/2)
	cam.defocus_disk_u = *cam.u.Scale(defocus_radius)
	cam.defocus_disk_v = *cam.v.Scale(defocus_radius)

}

// Figure out the color of the ray
func ray_color(r *ray, depth int, world hittable) *vec3 {

	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return &vec3{0, 0, 0}
	}

	var rec hit_record
	// tmin is 0.001 to get rid of shadow acne problem.
	if world.hit(r, 0.001, math.Inf(1), &rec) {

		var scattered ray
		var attenuation vec3
		if (*rec.material).scatter(r, &rec, &attenuation, &scattered) {
			return attenuation.Mult(*ray_color(&scattered, depth-1, world))
		}
		return &vec3{0, 0, 0}

		// direction := Random_on_hemisphere(&rec.normal)                           // Diffuse Material thing
		// direction := rec.normal.Add(*Random_unit_vec3())                         // True Lambertian diffusion
		// return ray_color(&ray{rec.point, *direction}, depth-1, world).Scale(0.7) // return 70% of the color from a bounce
	}

	unit_direction := Unit(&r.direction)
	a := (unit_direction.Y + 1.0) * 0.5
	return (&vec3{1.0, 1.0, 1.0}).Scale(1 - a).Add(*(&vec3{0.5, 0.7, 1.0}).Scale(a))
}

// Construct a camera ray originating from the origin and directed at randomly sampled
// point around the pixel location i, j.
func (cam *camera) get_ray(i float64, j float64) ray {
	offset := sample_square()
	pixel_sample := cam.pixel00_loc.Add(*cam.pixel_delta_u.Scale(i + offset.X)).Add(*cam.pixel_delta_v.Scale(j + offset.Y))
	ray_origin := cam.center
	if cam.defocus_angle > 0 {
		ray_origin = cam.defocus_disk_sample()
	}
	ray_direction := pixel_sample.Sub(*ray_origin)

	return ray{*ray_origin, *ray_direction}
}

func (cam *camera) defocus_disk_sample() *vec3 {
	p := random_in_unit_disk()
	return cam.center.Add(*cam.defocus_disk_u.Scale(p.X)).Add(*cam.defocus_disk_v.Scale(p.Y))
}

// Returns the vector to a random point in the [-.5,-.5]-[+.5,+.5] unit square.
func sample_square() vec3 {
	return vec3{rand.Float64() - 0.5, rand.Float64() - 0.5, 0}
}

// Helpers

// Clamp number between 0 and 1
func clamp(num float64) float64 {
	if num < 0.0 {
		return 0.0
	} else if num > 1 {
		return 1
	}
	return num
}

func linear_to_gamma(linear_component float64) float64 {
	if linear_component > 0 {
		return math.Sqrt(linear_component)
	}
	return 0
}
