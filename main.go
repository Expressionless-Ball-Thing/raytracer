package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {

	image_width, image_height := 256, 256

	upLeft := image.Point{0, 0}
	lowRight := image.Point{image_width, image_height}

	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	for j := 0; j < image_height; j++ {
		for i := 0; i < image_width; i++ {
			r := uint8(float64(i) / (float64(image_width - 1)) * 256)
			g := uint8(float64(j) / (float64(image_height - 1)) * 256)
			b := uint8(0)
			img.Set(i, j, color.NRGBA{r, g, b, 0xff})
		}
	}

	file, err := os.Create("image.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	png.Encode(file, img)

}
