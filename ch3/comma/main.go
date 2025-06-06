// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 73.

// Comma prints its argument numbers with a comma at each power of 1000.
//
// Example:
//
//	$ go build gopl.io/ch3/comma
//	$ ./comma 1 12 123 1234 1234567890
//	1
//	12
//	123
//	1,234
//	1,234,567,890
package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		recursiveResult := addCommas(arg, comma)
		iterativeResult := addCommas(arg, iterativeComma)
		fmt.Printf("  %s = %s\n", recursiveResult, iterativeResult)
	}
}

// !+
// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func addCommas(s string, commaFunc func(string) string) string {
	// Handle optional signs and "."
	if len(s) == 0 {
		return s
	}
	var sign string
	if s[0] == '+' || s[0] == '-' {
		sign = string(s[0])
		s = s[1:]
	} else {
		sign = ""
	}
	dotLocation := -1
	for i, c := range s {
		if c == '.' {
			dotLocation = i
			break
		}
	}
	if dotLocation == -1 {
		return sign + commaFunc(s)
	} else {
		// Split the string at the dot and add commas to the integer part.
		intPart := s[:dotLocation]
		fracPart := s[dotLocation:]
		return sign + commaFunc(intPart) + fracPart
	}
}

func iterativeComma(s string) string {
	var buf bytes.Buffer
	n := len(s)
	if n <= 3 {
		return s
	}
	// We can calculate the number of characters before the first comma.
	firstComma := n % 3
	for i := 0; i < firstComma; i++ {
		buf.WriteByte(s[i])
		if i == firstComma-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(s[firstComma]) // Write the first digit after the first comma.
	for i := firstComma + 1; i < n; i++ {
		if (i-firstComma)%3 == 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(s[i])
	}
	return buf.String()
}

//!-
