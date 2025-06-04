// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 61.
//!+

// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"os"
	"strconv"
)

type RenderingOptions struct {
	Xmin, Xmax float64
	Ymin, Ymax float64
	Zoom       float64
	Width      int
	Height     int
}

var defaultOptions = RenderingOptions{
	Xmin:   -2,
	Xmax:   2,
	Ymin:   -2,
	Ymax:   2,
	Zoom:   1,
	Width:  1024,
	Height: 1024,
}

func main() {
	web := flag.Bool("web", false, "Run in web mode")
	flag.Parse()
	if *web {
		fmt.Println("Running in web mode")
		http.HandleFunc("/fractals", webHandler)
		fmt.Println("Serving fractals at http://localhost:8080/fractals")
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		drawFractal("mandelbrot.png", mandelbrot, defaultOptions)
		drawFractal("newton.png", newton, defaultOptions)
		drawFractal("acos.png", acos, defaultOptions)
		drawFractal("sqrt.png", sqrt, defaultOptions)
	}
}

// parseRenderingOptionsFromQuery extracts rendering options from the HTTP request query parameters.
// It uses default values if the parameters are not provided or invalid.
func parseRenderingOptionsFromQuery(r *http.Request, def RenderingOptions) RenderingOptions {
	getFloat := func(key string, defVal float64) float64 {
		if v := r.URL.Query().Get(key); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
		return defVal
	}
	getInt := func(key string, defVal int) int {
		if v := r.URL.Query().Get(key); v != "" {
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
		return defVal
	}
	return RenderingOptions{
		Xmin:   getFloat("xmin", def.Xmin),
		Xmax:   getFloat("xmax", def.Xmax),
		Ymin:   getFloat("ymin", def.Ymin),
		Ymax:   getFloat("ymax", def.Ymax),
		Zoom:   getFloat("zoom", def.Zoom),
		Width:  getInt("width", def.Width),
		Height: getInt("height", def.Height),
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the fractals as PNG images
	w.Header().Set("Content-Type", "image/png")
	// Get the fractal from the query parameter
	fractal := r.URL.Query().Get("fractal")
	var fractalFunc func(complex128) color.Color
	switch fractal {
	case "mandelbrot":
		fractalFunc = mandelbrot
	case "newton":
		fractalFunc = newton
	case "acos":
		fractalFunc = acos
	case "sqrt":
		fractalFunc = sqrt
	default:
		fractalFunc = mandelbrot
	}
	opts := parseRenderingOptionsFromQuery(r, defaultOptions)
	renderFractal(w, fractalFunc, opts)
}

func drawFractal(filename string, fractalFunc func(complex128) color.Color, opts RenderingOptions) {
	handle, _ := os.Create(filename)
	defer handle.Close() // Close the file when we're done
	renderFractal(handle, fractalFunc, opts)
	fmt.Printf("Successfully drew fractal to %s\n", filename)
}

func renderFractal(out io.Writer, fractalFunc func(complex128) color.Color, opts RenderingOptions) {
	xspan := (opts.Xmax - opts.Xmin) / opts.Zoom
	yspan := (opts.Ymax - opts.Ymin) / opts.Zoom
	xstart := (opts.Xmin+opts.Xmax)/2 - xspan/2
	ystart := (opts.Ymin+opts.Ymax)/2 - yspan/2

	img := image.NewRGBA(image.Rect(0, 0, opts.Width, opts.Height))
	for py := 0; py < opts.Height; py++ {
		for px := 0; px < opts.Width; px++ {
			// x, y represent a pixel. To get the color of the pixel
			// we can do a supersampling of diving the pixel into 4 subpixels
			// and averaging the color of the subpixels.
			var subpixelColors []color.Color
			for subpy := 0; subpy < 2; subpy++ {
				deltay := float64(subpy) / 2.0
				for subpx := 0; subpx < 2; subpx++ {
					deltax := float64(subpx) / 2.0
					y := (float64(py)+deltay)/float64(opts.Height)*yspan + ystart
					x := (float64(px)+deltax)/float64(opts.Width)*xspan + xstart
					// Get the color for the pixel at (x, y)
					z := complex(x, y)
					subpixelColors = append(subpixelColors, fractalFunc(z))
				}
			}
			// Average the colors of the subpixels, allocate uint32 to match RGBA
			var r, g, b, a uint32
			for _, c := range subpixelColors {
				r1, g1, b1, a1 := c.RGBA()
				r += r1
				g += g1
				b += b1
				a += a1
			}
			// Divide by the number of subpixels to get the average color
			numSubpixels := uint32(len(subpixelColors))
			r /= numSubpixels
			g /= numSubpixels
			b /= numSubpixels
			a /= numSubpixels
			avgColor := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}
			img.Set(px, py, avgColor) // Set the pixel color in the image
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
	// z is the guess for the root of the equation z^4 - 1 = 0.
	const iterations = 37
	const epsilon = 1e-6

	// The four roots of z^4 = 1
	roots := []complex128{
		1,
		-1,
		complex(0, 1),
		complex(0, -1),
	}
	var n int
	for n = 0; n < iterations; n++ {
		z = z - (z*z*z*z-1)/(4*z*z*z)
		for i, root := range roots {
			if cmplx.Abs(z-root) < epsilon {
				// Color by root and shade by iterations
				shade := uint8(255 - uint8(255*n/iterations))
				switch i {
				case 0:
					return color.RGBA{shade, 0, 0, 255} // Red
				case 1:
					return color.RGBA{0, shade, 0, 255} // Green
				case 2:
					return color.RGBA{0, 0, shade, 255} // Blue
				case 3:
					return color.RGBA{shade, shade, 0, 255} // Yellow
				}
			}
		}
	}
	return color.Black
}
