// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 61.
//!+

// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/cmplx"
	"os"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)
	filename := "mandelbrot.png"
	var out io.Writer
	handle, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n. Using stdout instead.", filename, err)
		out = os.Stdout
	} else {
		fmt.Printf("Created file %s.\n", filename)
		defer handle.Close()
		out = handle
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(out, img) // NOTE: ignoring errors
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			hue := float64(n) / float64(iterations)
			return hsvToRGB(hue, 1, 1) // Convert to RGB
		}
	}
	return color.Black
}

// Simple HSV to RGB conversion
func hsvToRGB(h, s, v float64) color.Color {
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	switch i % 6 {
	case 0:
		return color.RGBA{uint8(v * 255), uint8(t * 255), uint8(p * 255), 255}
	case 1:
		return color.RGBA{uint8(q * 255), uint8(v * 255), uint8(p * 255), 255}
	case 2:
		return color.RGBA{uint8(p * 255), uint8(v * 255), uint8(t * 255), 255}
	case 3:
		return color.RGBA{uint8(p * 255), uint8(q * 255), uint8(v * 255), 255}
	case 4:
		return color.RGBA{uint8(t * 255), uint8(p * 255), uint8(v * 255), 255}
	default:
		return color.RGBA{uint8(v * 255), uint8(p * 255), uint8(q * 255), 255}
	}
}

//!-

// Some other interesting functions:

func acos(z complex128) color.Color {
	v := cmplx.Acos(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}

func sqrt(z complex128) color.Color {
	v := cmplx.Sqrt(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{128, blue, red}
}

// f(x) = x^4 - 1
//
// z' = z - f(z)/f'(z)
//
//	= z - (z^4 - 1) / (4 * z^3)
//	= z - (z - 1/z^3) / 4
func newton(z complex128) color.Color {
	const iterations = 37
	const contrast = 7
	for i := uint8(0); i < iterations; i++ {
		z -= (z - 1/(z*z*z)) / 4
		if cmplx.Abs(z*z*z*z-1) < 1e-6 {
			return color.Gray{255 - contrast*i}
		}
	}
	return color.Black
}
