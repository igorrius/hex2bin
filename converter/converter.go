package converter

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/marcinbor85/gohex"
)

// Intel HEX record types
const (
	DataRecord           = 0x00
	EndOfFileRecord      = 0x01
	ExtSegmentAddrRecord = 0x02
	StartSegmentRecord   = 0x03
	ExtLinearAddrRecord  = 0x04
	StartLinearRecord    = 0x05
)

// Calculate checksum for Intel HEX record
func calculateChecksum(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	return -sum
}

// Write Intel HEX record
func writeHexRecord(w io.Writer, recordType byte, address uint16, data []byte) error {
	record := make([]byte, 0, len(data)+4)
	record = append(record, byte(len(data)))
	record = append(record, byte(address>>8), byte(address))
	record = append(record, recordType)
	record = append(record, data...)
	checksum := calculateChecksum(record)

	hexData := strings.ToUpper(hex.EncodeToString(data))
	_, err := fmt.Fprintf(w, ":%02X%04X%02X%s%02X\n",
		len(data), address, recordType, hexData, checksum)
	return err
}

func BinToIntelHex(inputFile string, outputFile string) error {
	// Read the binary input file
	binData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	// Create output file
	ofw, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer ofw.Close()

	w := bufio.NewWriter(ofw)

	// Write extended linear address record for the first segment
	err = writeHexRecord(w, ExtLinearAddrRecord, 0x0000, []byte{0x00, 0x00})
	if err != nil {
		return fmt.Errorf("error writing extended address: %v", err)
	}

	// Write data records in 32-byte chunks
	for i := 0; i < len(binData); i += 32 {
		end := i + 32
		if end > len(binData) {
			end = len(binData)
		}
		chunk := binData[i:end]
		err = writeHexRecord(w, DataRecord, uint16(i), chunk)
		if err != nil {
			return fmt.Errorf("error writing data chunk: %v", err)
		}
	}

	// Write end of file record
	err = writeHexRecord(w, EndOfFileRecord, 0x0000, nil)
	if err != nil {
		return fmt.Errorf("error writing EOF: %v", err)
	}

	return w.Flush()
}

func IntelHexToBin(inputFile string, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(file)

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(file)
	if err != nil {
		return fmt.Errorf("error parsing Intel HEX file: %v", err)
	}

	// Get all data segments
	segments := mem.GetDataSegments()
	if len(segments) == 0 {
		return fmt.Errorf("no data segments found in Intel HEX file")
	}

	// Create output file
	ofw, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer ofw.Close()

	// Write binary data in order of addresses
	for _, segment := range segments {
		_, err = ofw.Write(segment.Data)
		if err != nil {
			return fmt.Errorf("error writing binary data: %v", err)
		}
	}

	return nil
}
