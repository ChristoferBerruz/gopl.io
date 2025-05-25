// Package genconv is a general purpose conversion library
package genconv

import "fmt"

type Celsius float64
type Fahrenheit float64
type Kelvin float64

type Meter float64
type Feet float64

type Pound float64
type Kilogram float64

const (
	KelvinToCShift        = 273.15
	MeterToFeetScaler     = 3.28084
	KilogramToPoundScaler = 2.20462
)

func (c Celsius) String() string    { return fmt.Sprintf("%g°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%g°F", f) }
func (k Kelvin) String() string     { return fmt.Sprintf("%gK", k) }
func (m Meter) String() string      { return fmt.Sprintf("%gm", m) }
func (f Feet) String() string       { return fmt.Sprintf("%gft", f) }
func (p Pound) String() string      { return fmt.Sprintf("%glb", p) }
func (k Kilogram) String() string   { return fmt.Sprintf("%gkg", k) }

// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

// CToK converts a Celsius temperature to a Kelvin
func CToK(c Celsius) Kelvin { return Kelvin(c + KelvinToCShift) }

// KToC converts a Kelvin into Celsius
func KToC(k Kelvin) Celsius { return Celsius(k - KelvinToCShift) }

// FToK converts a Fahrenheit into Kelvin
func FToK(f Fahrenheit) Kelvin { return CToK(FToC(f)) }

// KToF converts a Kelvin into a Fahrenheit
func KToF(k Kelvin) Fahrenheit { return CToF(KToC(k)) }

// MeterToFeet converts a Meter into a Foot
func MeterToFeet(m Meter) Feet { return Feet(m * MeterToFeetScaler) }

// FeetToMeter converts a Foot into a Meter
func FeetToMeter(f Feet) Meter { return Meter(f / MeterToFeetScaler) }

// KiloToPound converts a Kilogram into a Pound
func KiloToPound(k Kilogram) Pound { return Pound(k * KilogramToPoundScaler) }

// PoundToKilo converts a Pound into a Kilogram
func PoundToKilo(p Pound) Kilogram { return Kilogram(p / KilogramToPoundScaler) }
