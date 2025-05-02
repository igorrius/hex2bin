# hex2bin

A command-line tool for converting between Intel HEX and binary file formats. This tool is particularly useful for embedded systems development, firmware updates, and memory programming.

## Description

hex2bin is a simple yet powerful tool that allows you to:
- Convert Intel HEX files to binary format
- Convert binary files to Intel HEX format
- Preserve memory addresses and data organization during conversion
- Handle extended linear addresses and data records properly

The tool is written in Go and provides a reliable way to work with Intel HEX files, which are commonly used in embedded systems and microcontroller programming.

## Installation

### Using Go Install

If you have Go installed on your system, you can install hex2bin using:

```bash
go install github.com/igorrius/hex2bin@latest
```

### Using Pre-built Binaries

#### Linux

1. Download the latest Linux binary from the [Releases](https://github.com/igorrius/hex2bin/releases) page
2. Make it executable:
   ```bash
   chmod +x hex2bin
   ```
3. Move it to your PATH:
   ```bash
   sudo mv hex2bin /usr/local/bin/
   ```

#### Windows

1. Download the latest Windows binary from the [Releases](https://github.com/igorrius/hex2bin/releases) page
2. Add the directory containing the binary to your PATH environment variable

## Usage

The tool supports two conversion modes:
- `hex2bin`: Intel HEX to binary
- `bin2hex`: binary to Intel HEX

### Arguments

- `input_file`: Path to the input file (**required**)
- `output_file`: Path to the output file (**optional**)
- `mode`: Conversion mode (`hex2bin` or `bin2hex`, **optional**)

If `output_file` is not provided, it will be automatically determined from the input file name and mode.
If `mode` is not provided, it will be inferred from the input file extension:
- `.hex` → `hex2bin`
- `.bin` → `bin2hex`

If both `output_file` and `mode` are omitted, both will be inferred from the input file extension.

### Examples

#### Minimal (auto mode and output):
```bash
hex2bin firmware.hex         # Converts to firmware.bin (mode: hex2bin)
hex2bin firmware.bin         # Converts to firmware.hex (mode: bin2hex)
```

#### Specify output file, auto mode:
```bash
hex2bin firmware.hex my.bin  # Converts to my.bin (mode: hex2bin)
hex2bin firmware.bin my.hex  # Converts to my.hex (mode: bin2hex)
```

#### Specify mode, auto output:
```bash
hex2bin firmware.hex hex2bin   # Converts to firmware.bin
hex2bin firmware.bin bin2hex   # Converts to firmware.hex
```

#### Full explicit:
```bash
hex2bin firmware.hex out.bin hex2bin
hex2bin firmware.bin out.hex bin2hex
```

## Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/igorrius/hex2bin.git
   cd hex2bin
   ```

2. Build the project:
   ```bash
   make
   ```

3. Install the binary:
   ```bash
   make install
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
