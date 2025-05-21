/*
This program prints the name of the program and the command line arguments
one per line. It also includes a basic benchmarking between
parsing the cmd args one at a time and using strings.join instead.
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func parseAndPrint() {
	var s, sep string
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

func joinAndPrint() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

func parseAndPrintIncludingName() {
	// Implements exercises 1.1 and 1.2
	for i, arg := range os.Args[:] {
		fmt.Println(i, arg)
	}
}

func timeIt(name string, f func()) {
	start := time.Now()
	f()
	elapsedSeconds := time.Since(start).Seconds()
	msg := name + " took " + fmt.Sprintf("%f", elapsedSeconds) + " seconds"
	fmt.Println(msg)
}

func main() {
	fmt.Println("Running exercise")
	parseAndPrintIncludingName()
	fmt.Println("Running basic benchmarking")
	timeIt("CustomParsing", parseAndPrint)
	timeIt("UsingJoin", joinAndPrint)
}
