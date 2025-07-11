// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 83.

// The sha256 command computes the SHA256 hash (an array) of a string.
package main

import "fmt"

//!+
import "crypto/sha256"

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)
	fmt.Printf("%d distinct bits\n", numberOfDistinctBits(c1, c2))
	// Output:
	// 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	// 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	// false
	// [32]uint8
}

// numberOfDistinctBits counts the number of bits that differ between two SHA256 hashes.
func numberOfDistinctBits(oneHash, otherHash [32]byte) int {
	distinctBits := 0
	for i := 0; i < len(oneHash); i++ {
		// XOR the bytes of the two hashes to find differing bits
		// then count the number of set bits in the result.
		xorResult := oneHash[i] ^ otherHash[i]
		for j := 0; j < 8; j++ {
			if (xorResult & (1 << j)) != 0 {
				distinctBits++
			}
		}
	}
	return distinctBits
}

//!-
