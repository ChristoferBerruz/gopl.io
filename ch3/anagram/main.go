// Anagram Checker
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Extract two strings from command line arguments
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run main.go <string1> <string2>")
	}
	str1 := os.Args[1]
	str2 := os.Args[2]
	if isAnagram(str1, str2) {
		fmt.Printf("'%s' and '%s' are anagrams.\n", str1, str2)
	} else {
		fmt.Printf("'%s' and '%s' are not anagrams.\n", str1, str2)
	}
}

// preProcessString processes the input string to make it suitable for anagram comparison.
func preProcessString(s string) string {
	// call lowercase on both string to ensure case insensitivity for ASCII values
	s = strings.ToLower(s)
	// remove whitespaces to ensure that we consider strings, not only single words.
	s = strings.ReplaceAll(s, " ", "")
	return s
}

// isAnagram checks if two strings are anagrams of each other.
func isAnagram(str1, str2 string) bool {
	// Assume that both strings are UTF-8 encoded.
	str1 = preProcessString(str1)
	str2 = preProcessString(str2)
	if len(str1) != len(str2) {
		return false
	}
	counts := make(map[rune]int) // Rune because we want to handle Unicode characters
	for _, r := range str1 {
		counts[r]++
	}
	for _, r := range str2 {
		counts[r]--
		if counts[r] < 0 {
			// To many occurrences of r in str2 compared to str1
			return false
		}
	}
	return true
}
