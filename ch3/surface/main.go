// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
)

type Point struct {
	x, y, z float64
}

type Polygon struct {
	a, b, c, d Point
}

type RenderingOptions struct {
	width, height, cells     int
	xyrange, xyscale, zscale float64
	angle                    float64
}

type UserOptions struct {
	width, height, cells int
}

func getRenderingOptions(opts UserOptions) RenderingOptions {
	// Default rendering options
	xyrange := 30.0
	xyscale := float64(opts.width) / 2 / xyrange
	zscale := float64(opts.height) * 0.4
	return RenderingOptions{
		width:   opts.width,
		height:  opts.height,
		cells:   opts.cells,
		xyrange: xyrange,
		xyscale: xyscale,
		zscale:  zscale,
	}
}

var defaultOptions RenderingOptions

func init() {
	width, height := 600, 320 // canvas size in pixels
	cells := 100              // number of grid cells
	defaultOptions = getRenderingOptions(UserOptions{
		width:  width,
		height: height,
		cells:  cells,
	})
}                         // to protect currentOptions
const rad30 = math.Pi / 6 // 30 degrees in radians

var sin30, cos30 = math.Sin(rad30), math.Cos(rad30) // sin(30°), cos(30°)

func main() {
	web := flag.Bool("web", false, "run as web server")
	flag.Parse()
	if *web {
		fmt.Println("Running as web server. Output will be served at http://localhost:8080/surface")
		http.HandleFunc("/surface", surfaceHandler)
		http.ListenAndServe(":8080", nil)
		return
	}
	// If no web flag is set, run as a command line tool
	// check whether there is a command line argument for the function to use
	f := defaultSurface // default function
	filename := "surface.svg"
	if flag.NArg() > 0 {
		f = getSurfaceFunction(flag.Arg(0))
	}
	var out io.Writer
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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
	renderSVG(out, f, defaultOptions)
}

func getSurfaceFunction(name string) func(x, y float64) float64 {
	switch name {
	case "eggbox":
		return eggBox
	case "saddle":
		return saddle
	case "moguls":
		return moguls
	default:
		return defaultSurface // default function if no match
	}
}

func surfaceHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to SVG
	w.Header().Set("Content-Type", "image/svg+xml")

	// Get the surface type from the query parameter
	surfaceType := r.URL.Query().Get("type")
	userHeight := r.URL.Query().Get("height")
	userWidth := r.URL.Query().Get("width")
	userCells := r.URL.Query().Get("cells")
	// some defaults
	userOpts := UserOptions{
		width:  defaultOptions.width,
		height: defaultOptions.height,
		cells:  defaultOptions.cells,
	}
	if userHeight != "" {
		if height, err := strconv.Atoi(userHeight); err == nil {
			userOpts.height = height
		}
	}
	if userWidth != "" {
		if width, err := strconv.Atoi(userWidth); err == nil {
			userOpts.width = width
		}
	}
	if userCells != "" {
		if cells, err := strconv.Atoi(userCells); err == nil {
			userOpts.cells = cells
		}
	}
	userSurfaceOptions := getRenderingOptions(userOpts)
	f := getSurfaceFunction(surfaceType)
	// Render the SVG directly to the response writer
	renderSVG(w, f, userSurfaceOptions)
}

func renderSVG(out io.Writer, f func(x, y float64) float64, opts RenderingOptions) {
	// Collect all valid polygons first
	cells := opts.cells
	width := opts.width
	height := opts.height
	var polygons []Polygon
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			a := corner(i+1, j, f, opts)
			b := corner(i, j, f, opts)
			c := corner(i, j+1, f, opts)
			d := corner(i+1, j+1, f, opts)
			if anyIsNan(a.x, a.y, b.x, b.y, c.x, c.y, d.x, d.y) {
				continue
			}
			polygons = append(polygons, Polygon{a, b, c, d})
		}
	}
	// Find the minimum and maximum z values
	minZ, maxZ := math.Inf(1), math.Inf(-1)
	for _, poly := range polygons {
		for _, p := range []Point{poly.a, poly.b, poly.c, poly.d} {
			if p.z < minZ {
				minZ = p.z
			}
			if p.z > maxZ {
				maxZ = p.z
			}
		}
	}
	// Dump the SVG header
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	// Output polygons with color
	for _, poly := range polygons {
		avgZ := (poly.a.z + poly.b.z + poly.c.z + poly.d.z) / 4
		color := colorForZ(avgZ, minZ, maxZ)
		fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s'/>\n",
			poly.a.x, poly.a.y, poly.b.x, poly.b.y, poly.c.x, poly.c.y, poly.d.x, poly.d.y, color)
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

func corner(i, j int, f func(x, y float64) float64, opts RenderingOptions) Point {
	// Find point (x,y) at corner of cell (i,j).
	xyrange := opts.xyrange
	cells := float64(opts.cells)
	xyscale := opts.xyscale
	zscale := opts.zscale
	width := float64(opts.width)
	height := float64(opts.height)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)
	// Handle the case where z can be a non-finite
	// value that will produce invalid polygons
	if math.IsNaN(z) || math.IsInf(z, 0) {
		return Point{math.NaN(), math.NaN(), math.NaN()}
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return Point{sx, sy, z}
}

func colorForZ(z, minZ, maxZ float64) string {
	if maxZ == minZ {
		return "#888888"
	}
	t := (z - minZ) / (maxZ - minZ)
	r := int(255 * t)
	g := 0
	b := int(255 * (1 - t))
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
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
