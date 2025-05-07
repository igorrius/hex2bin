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

// GetOutputFileName determines the output file name based on input file and mode
func GetOutputFileName(inputFile, mode string) (string, error) {
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
	output := base + ext

	// Add suffix if file exists
	if _, err := os.Stat(output); err == nil {
		suffix := 1
		for {
			candidate := fmt.Sprintf("%s_%d%s", base, suffix, ext)
			if _, err := os.Stat(candidate); os.IsNotExist(err) {
				output = candidate
				break
			}
			suffix++
		}
	}

	return output, nil
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
	writeAll := false
	recordSize := 16
	args := os.Args[1:]
	// Check for --all and --record-bytes flags
	for i := 0; i < len(args); {
		if args[i] == "--all" {
			writeAll = true
			args = append(args[:i], args[i+1:]...)
			continue
		}
		if strings.HasPrefix(args[i], "--record-bytes=") {
			var n int
			_, err := fmt.Sscanf(args[i], "--record-bytes=%d", &n)
			if err == nil && n > 0 {
				recordSize = n
			}
			args = append(args[:i], args[i+1:]...)
			continue
		}
		i++
	}

	if len(args) < 1 || len(args) > 3 {
		fmt.Printf("hex2bin v%s (built %s)\n", version, buildTime)
		fmt.Println("Usage: hex2bin <input_file> [output_file] [mode] [--all] [--record-bytes=N]")
		fmt.Println("Mode: 'bin2hex' for binary to Intel HEX conversion")
		fmt.Println("      'hex2bin' for Intel HEX to binary conversion")
		fmt.Println("      '--all' to write all data (no sparse HEX, for bin2hex mode)")
		fmt.Println("      '--record-bytes=N' to set bytes per HEX record (default 16, e.g., 16 or 32)")
		os.Exit(1)
	}

	inputFile := args[0]
	outputFile := ""
	mode := ""

	if len(args) == 3 {
		outputFile = args[1]
		mode = args[2]
	} else if len(args) == 2 {
		arg := args[1]
		if arg == "bin2hex" || arg == "hex2bin" {
			mode = arg
			var err error
			outputFile, err = GetOutputFileName(inputFile, mode)
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
		var err error
		mode, err = getModeFromInputExt(inputFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		outputFile, err = GetOutputFileName(inputFile, mode)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	var err error
	switch mode {
	case "bin2hex":
		err = converter.BinToIntelHexWithMode(inputFile, outputFile, writeAll, recordSize)
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
