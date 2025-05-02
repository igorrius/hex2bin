package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/igorrius/hex2bin/converter"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

// getOutputFileName determines the output file name based on input file and mode
func getOutputFileName(inputFile, mode string) (string, error) {
	ext := ""
	switch mode {
	case "bin2hex":
		ext = ".hex"
	case "hex2bin":
		ext = ".bin"
	default:
		return "", fmt.Errorf("cannot determine output file extension for mode '%s' (should be 'bin2hex' or 'hex2bin')", mode)
	}
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	return base + ext, nil
}

// getModeFromInputExt tries to determine the mode from the input file extension
func getModeFromInputExt(inputFile string) (string, error) {
	ext := strings.ToLower(filepath.Ext(inputFile))
	switch ext {
	case ".hex":
		return "hex2bin", nil
	case ".bin":
		return "bin2hex", nil
	default:
		return "", fmt.Errorf("cannot determine mode from input file extension '%s' (should be .hex or .bin)", ext)
	}
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 4 {
		fmt.Printf("hex2bin v%s (built %s)\n", version, buildTime)
		fmt.Println("Usage: hex2bin <input_file> [output_file] [mode]")
		fmt.Println("Mode: 'bin2hex' for binary to Intel HEX conversion")
		fmt.Println("      'hex2bin' for Intel HEX to binary conversion")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := ""
	mode := ""

	if len(os.Args) == 4 {
		outputFile = os.Args[2]
		mode = os.Args[3]
	} else if len(os.Args) == 3 {
		// Could be output file or mode
		arg := os.Args[2]
		if arg == "bin2hex" || arg == "hex2bin" {
			mode = arg
			var err error
			outputFile, err = getOutputFileName(inputFile, mode)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			outputFile = arg
			var err error
			mode, err = getModeFromInputExt(inputFile)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		// Only input file provided
		var err error
		mode, err = getModeFromInputExt(inputFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		outputFile, err = getOutputFileName(inputFile, mode)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

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
