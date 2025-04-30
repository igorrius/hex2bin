package converter

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/marcinbor85/gohex"
)

func BinToIntelHex(inputFile string, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(file)

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(file)
	if err != nil {
		panic(err)
	}
	for _, segment := range mem.GetDataSegments() {
		fmt.Printf("%+v\n", segment)
	}
	bytes := mem.ToBinary(0xFFF0, 128, 0x00)
	fmt.Printf("%v\n", bytes)

	return nil
}

func IntelHexToBin(inputFile string, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	mem := gohex.NewMemory()
	mem.SetStartAddress(0x80008000)
	mem.AddBinary(0x10008000, []byte{0x01, 0x02, 0x03, 0x04})
	mem.AddBinary(0x20000000, make([]byte, 256))

	mem.DumpIntelHex(file, 16)

	return nil
}
