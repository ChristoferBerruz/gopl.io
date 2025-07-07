// Exercise 4.2: Implement a program that prints
// the (SHA256, SHA384, SHA512) hashes of its standard input.

package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
)

func main() {
	hashingAlgorithm := flag.String("algorithm", "sha256", "Hashing algorithm to use: sha256, sha384, sha512")
	flag.Parse()
	toHash := flag.Arg(0)
	fmt.Printf("User input: %q\n", toHash)
	switch *hashingAlgorithm {
	case "sha256":
		hash := sha256.Sum256([]byte(toHash))
		fmt.Printf("SHA256: %x\n", hash)
	case "sha384":
		hash := sha512.Sum384([]byte(toHash))
		fmt.Printf("SHA384: %x\n", hash)
	case "sha512":
		hash := sha512.Sum512([]byte(toHash))
		fmt.Printf("SHA512: %x\n", hash)
	default:
		fmt.Println("Unsupported hashing algorithm. Use sha256, sha384, or sha512.")
		return
	}
}
