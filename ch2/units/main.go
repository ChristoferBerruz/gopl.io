// units takes multiple arguments and illustrates different conversions based on the possible units
// for the value.
package main

import (
	"fmt"
	"gopl.io/ch2/genconv"
	"os"
	"strconv"
)

func main() {
	for _, arg := range os.Args[1:] {
		fmt.Printf("Showing possible conversions for %s\n", arg)
		raw, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %s\n", err)
			os.Exit(1)
		}
		// Handling temperatures
		f := genconv.Fahrenheit(raw)
		c := genconv.Celsius(raw)
		k := genconv.Kelvin(raw)
		fmt.Printf("%s = %s = %s\n", f, genconv.FToC(f), genconv.FToK(f))
		fmt.Printf("%s = %s = %s\n", c, genconv.CToF(c), genconv.CToK(c))
		fmt.Printf("%s = %s = %s\n", k, genconv.KToF(k), genconv.KToC(k))
		// Handling distances
		meter := genconv.Meter(raw)
		feet := genconv.Feet(raw)
		fmt.Printf("%s = %s, %s = %s\n", meter, genconv.MeterToFeet(meter), feet, genconv.FeetToMeter(feet))
		// Handling weights
		kilos := genconv.Kilogram(raw)
		lbs := genconv.Pound(raw)
		fmt.Printf("%s = %s, %s = %s\n", kilos, genconv.KiloToPound(kilos), lbs, genconv.PoundToKilo(lbs))
	}
}
