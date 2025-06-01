// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	filename := "surface.svg"
	// check whether there is a command line argument for the function to use
	var f func(x, y float64) float64
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "eggbox":
			f = eggBox
			filename = "eggbox.svg"
		case "saddle":
			f = saddle
			filename = "saddle.svg"
		case "moguls":
			f = moguls
			filename = "moguls.svg"
		default:
			fmt.Printf("Unknown function %q, using default surface function.\n", os.Args[1])
			f = defaultSurface
		}
	} else {
		fmt.Println("No function argument provided, using default surface function.")
		f = defaultSurface // default function if no argument is provided
	}
	var out io.Writer
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Error opening file %q: %v", filename, err)
		fmt.Println("Defaulting to stdout")
		out = os.Stdout
	} else {
		fmt.Printf("Successfully opened file %q. Output will be saved there.", filename)
		out = file
		// Defer closing the file before main returns
		defer file.Close()
	}
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j, f)
			bx, by := corner(i, j, f)
			cx, cy := corner(i, j+1, f)
			dx, dy := corner(i+1, j+1, f)
			if anyIsNan(ax, ay, bx, by, cx, cy, dx, dy) {
				continue
			}
			fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintf(out, "</svg>")
}

func anyIsNan(values ...float64) bool {
	for _, value := range values {
		if math.IsNaN(value) {
			return true
		}
	}
	return false
}

func corner(i, j int, f func(x, y float64) float64) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)
	// Handle the case where z can be a non-finite
	// value that will produce invalid polygons
	if math.IsNaN(z) || math.IsInf(z, 0) {
		return math.NaN(), math.NaN()
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func defaultSurface(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

func eggBox(x, y float64) float64 {
	// Eggbox function: f(x,y) = sin(x) * cos(y)
	return math.Sin(x) * math.Cos(y)
}

func saddle(x, y float64) float64 {
	// Saddle function: f(x,y) = x^2 - y^2
	return (x*x - y*y) / 100.0
}

func moguls(x, y float64) float64 {
	// Moguls: wavy bumps using sine and cosine
	return math.Sin(x) * math.Cos(y) / 2
}

//!-
