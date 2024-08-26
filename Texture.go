package main

import (
	"image"
	"image/png"
	"io"
	"math"
)

type Texture interface {
	value(u, v float64, point Vec3) Vec3
}

// Constant Color Texture
type Solid struct {
	albedo Vec3
}

// Create a Constant Color Texture from a Vec3 of the RGB intensities
func NewSolidTexture(albedo Vec3) *Solid {
	return &Solid{
		albedo,
	}
}

func (solid *Solid) value(u, v float64, point Vec3) Vec3 {
	return solid.albedo
}

// Checker Texture
type Checker struct {
	inv_scale float64
	even, odd *Texture
}

// Create a checker texture from two other texture material
func NewCheckerTexture(scale float64, even, odd Texture) *Checker {
	return &Checker{
		1.0 / scale,
		&even,
		&odd,
	}
}

// Create a checker texture from two Color Vectors
func NewCheckerFromColor(scale float64, c1, c2 Vec3) *Checker {
	return NewCheckerTexture(scale, NewSolidTexture(c1), NewSolidTexture(c2))
}

func (checker *Checker) value(u, v float64, point Vec3) Vec3 {

	isEven := int(math.Floor(checker.inv_scale*point[0])+math.Floor(checker.inv_scale*point[1])+math.Floor(checker.inv_scale*point[2]))%2 == 0
	if isEven {
		return (*checker.even).value(u, v, point)
	} else {
		return (*checker.odd).value(u, v, point)
	}
}

// Image Texture
type Image struct {
	data          image.Image
	width, height int
}

func NewImageTexture(rc io.Reader) (*Image, error) {
	im, err := png.Decode(rc)
	if err != nil {
		return nil, err
	}

	im.Bounds()

	return &Image{
		im,
		im.Bounds().Max.X,
		im.Bounds().Max.Y,
	}, err

}

func (im *Image) value(u, v float64, point Vec3) Vec3 {
	if im.height <= 0 {
		return NewVec3(0, 1, 1)
	}

	// Clamp input texture coordinates to [0,1] x [1,0]
	temp_u := clamp(u)
	temp_v := 1 - clamp(v) // Flip V to image coordinates
	i := int(temp_u * float64(im.width))
	j := int(temp_v * float64(im.height))
	pixel := im.data.At(i, j)

	// Divide by 0xffff because that's the max value RGBA can be in golang implementation.
	r, g, b, _ := pixel.RGBA()
	return NewVec3(float64(r)/65535, float64(g)/65535, float64(b)/65535)
}

// Perlin Noise Texture
type Noise struct {
	scale float64
}

// Make a new perlin noise texture
func NewNoise(scale float64) *Noise {
	perlin()
	return &Noise{scale}
}

func (noi *Noise) value(u, v float64, point Vec3) Vec3 {
	// thing := NewVec3(1, 1, 1).Scale(1 + noise(point.Scale(noi.scale))).Scale(0.5)
	// thing := NewVec3(1, 1, 1).Scale(turbulence(point, 7))
	thing := NewVec3(0.5, 0.5, 0.5).Scale(math.Sin(turbulence(point, 7)*10+noi.scale*point[2]) + 1)
	return thing
}
