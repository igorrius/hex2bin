package main

import (
	"fmt"
	"os"

	"github.com/igorrius/hex2bin/converter"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("hex2bin %s (built %s)\n", version, buildTime)
		fmt.Println("Usage: hex2bin <input_file> <output_file> <mode>")
		fmt.Println("Mode: 'bin2hex' for binary to Intel HEX conversion")
		fmt.Println("      'hex2bin' for Intel HEX to binary conversion")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	mode := os.Args[3]

	var err error
	switch mode {
	case "bin2hex":
		err = converter.BinToIntelHex(inputFile, outputFile)
	case "hex2bin":
		err = converter.IntelHexToBin(inputFile, outputFile)
	default:
		fmt.Println("Error: Invalid mode. Use 'bin2hex' or 'hex2bin'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}
