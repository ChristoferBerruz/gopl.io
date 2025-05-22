// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.
//!+main

// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
)

//!-main
// Packages not needed by version in book.
import (
	"log"
	"net/http"
	"time"
)

//!+main

// Exercise 1.5: Change palette from black on white to green on black
var greenColor = color.RGBA{0x00, 0x80, 0x00, 0xff}
var redColor = color.RGBA{0x80, 0x00, 0x00, 0xff}
var blueColor = color.RGBA{0x00, 0x00, 0x80, 0xff}
var yellowColor = color.RGBA{0xff, 0xff, 0x00, 0xff}
var palette = []color.Color{color.Black, greenColor, redColor, blueColor, yellowColor}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	//!-main
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())
	cycles := 5.0

	if len(os.Args) > 1 && os.Args[1] == "web" {
		//!+http
		handler := func(w http.ResponseWriter, r *http.Request) {
			// Exercise 1.12
			// Handle a user defined cycles value
			err := r.ParseForm()
			if err == nil {
				requestedValue, getCycleErr := strconv.Atoi(r.FormValue("cycles"))
				if getCycleErr == nil {
					cycles = float64(requestedValue)
				}
			}
			lissajous(w, cycles)

		}
		http.HandleFunc("/", handler)
		//!-http
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	//!+main
	lissajous(os.Stdout, cycles)
}

func lissajous(out io.Writer, cycles float64) {
	const ( // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	var colorIndex uint8
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			// Pick a color from the pallet. Exclude the 0th color (background), tho
			colorIndex = uint8(rand.Intn(len(palette)-1) + 1)
			// Exercise 1.6: modify the program to produce images of multiple color
			// This can be done by picking a color at random.
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				colorIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main
