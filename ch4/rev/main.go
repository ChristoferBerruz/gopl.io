// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 86.

// Rev reverses a slice.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	//!+array
	a := [...]int{0, 1, 2, 3, 4, 5}
	reverseArrayPointer(&a)
	fmt.Println(a) // "[5 4 3 2 1 0]"
	//!-array

	//!+slice
	s := []int{0, 1, 2, 3, 4, 5}
	// Rotate s left by two positions.
	rotate(s, 2)
	fmt.Println(s) // "[2 3 4 5 0 1]"
	//!-slice

	//!+eliminate
	// Eliminate adjacent duplicates from a slice of strings.
	eliminateStrings := []string{"a", "b", "b", "c", "c", "c", "d"}
	eliminateStrings = eliminateAdjacentDuplicates(eliminateStrings)
	fmt.Println(eliminateStrings) // "[a b c d]"

	// Squash spaces in a byte slice.
	squashSpacesBytes := []byte("  a   b   c  ")
	squashed := squashSpaces(squashSpacesBytes)
	fmt.Printf("%q\n", squashed) // " a b c "

	// Reverse UTF-8 encoded byte slice.
	someUTF8String := []byte("Hello, 世界")
	someUTF8String = reverseUTF8(someUTF8String)
	fmt.Printf("%s\n", string(someUTF8String)) // "界世 ,olleH"

	// Interactive test of reverse.
	input := bufio.NewScanner(os.Stdin)
outer:
	for input.Scan() {
		var ints []int
		for _, s := range strings.Fields(input.Text()) {
			x, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue outer
			}
			ints = append(ints, int(x))
		}
		reverse(ints)
		fmt.Printf("%v\n", ints)
	}
	// NOTE: ignoring potential errors from input.Err()
}

// !+rev
// reverse reverses a slice of ints in place.
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// reverseArrayPointer reverses an array of ints in place using a pointer.
func reverseArrayPointer(s *[6]int) {
	// Note that arrays are of type [N]T, which means that the length
	// is part of the type. Using *[]int would mean a pointer to a slice
	for i, j := 0, len(*s)-1; i < j; i, j = i+1, j-1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}

// rotate rotates a slice of ints left by n positions in a single pass.
func rotate(s []int, n int) {
	n = n % len(s) // Ensure n is within bounds
	if n <= 0 || len(s) == 0 {
		// No rotation needed if n is 0 or slice is empty.
		return
	}
	tmp := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		tmp[i] = s[(i+n)%len(s)]
	}
	copy(s, tmp)
}

// eliminateAdjacentDuplicates removes adjacent duplicates from a slice of ints in-place.
func eliminateAdjacentDuplicates(s []string) []string {
	if len(s) == 0 {
		return s
	}
	j := 0 // Index for the last unique element
	for i := 1; i < len(s); i++ {
		if s[i] != s[j] { // Compare current with last unique
			j++         // Move to the next position for unique element
			s[j] = s[i] // Update the slice with the new unique element
		}
	}
	return s[:j+1] // Return the slice up to the last unique element
}

// squashSpaces squashes each run of adjacent Unicode spaces in b into a single ASCII space.
// It returns the new length of the slice.
func squashSpaces(b []byte) []byte {
	write := 0
	read := 0
	spacePending := false

	for read < len(b) {
		r, size := utf8.DecodeRune(b[read:])
		if unicode.IsSpace(r) {
			if !spacePending {
				b[write] = ' '
				write++
				spacePending = true
			}
			// skip this space rune
		} else {
			n := utf8.EncodeRune(b[write:], r)
			write += n
			spacePending = false
		}
		read += size
	}
	return b[:write]
}

// reverseUTF8 reverses the characters of a UTF-8 encoded []byte slice in place.
func reverseUTF8(b []byte) []byte {
	var runes []rune
	originalSlice := b // Keep the original slice for reference
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		runes = append(runes, r)
		b = b[size:]
	}
	// Reverse runes
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	copy(originalSlice, []byte(string(runes)))
	return originalSlice[:]
}

//!-rev
