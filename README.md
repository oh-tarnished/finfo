# finfo - File Information CLI Tool

A powerful, colorful CLI tool for displaying comprehensive file information with advanced analysis capabilities.

## Features

- **Colorized Output** - Beautiful, readable terminal output with distinct colors per section
- **Smart Size Formatting** - Intelligent display of file sizes in relevant units
- **File Type Detection** - MIME type, format detection (ELF, Mach-O, scripts), encoding
- **Binary Analysis** - Linked libraries, stripped status, code signatures (macOS)
- **Hash Calculation** - MD5, SHA256, SHA512 checksums
- **File Comparison** - Git-like diff output comparing two files
- **Symlink Resolution** - Complete symlink chain visualization
- **Command Resolution** - Automatic PATH lookup for commands
- **Library Search** - Find and analyze `.so`, `.a`, `.dylib` files
- **Cross-Platform** - Works on macOS and Linux

## Installation

### Go install (all platforms)

Requires Go 1.21 or later:

```bash
go install github.com/oh-tarnished/finfo@latest
```

Ensure `$HOME/go/bin` is in your PATH.

### macOS

#### Homebrew
```bash
brew tap oh-tarnished/finfo
brew install finfo
```

#### Manual Install
```bash
# Download latest release
curl -LO https://github.com/oh-tarnished/finfo/releases/latest/download/finfo-darwin-arm64

# Make executable and move to PATH
chmod +x finfo-darwin-arm64
sudo mv finfo-darwin-arm64 /usr/local/bin/finfo
```

### Linux

#### Debian/Ubuntu (deb)
```bash
curl -LO https://github.com/oh-tarnished/finfo/releases/latest/download/finfo_<version>_linux_amd64.deb
sudo dpkg -i finfo_<version>_linux_amd64.deb
```

#### RHEL/Fedora (rpm)
```bash
curl -LO https://github.com/oh-tarnished/finfo/releases/latest/download/finfo_<version>_linux_amd64.rpm
sudo dnf install finfo_<version>_linux_amd64.rpm
```

#### Arch Linux
```bash
curl -LO https://github.com/oh-tarnished/finfo/releases/latest/download/finfo_<version>_linux_amd64.pkg.tar.zst
sudo pacman -U finfo_<version>_linux_amd64.pkg.tar.zst
```

### Build from source

```bash
git clone https://github.com/oh-tarnished/finfo.git
cd finfo
go build -o finfo
sudo mv finfo /usr/local/bin/
```

## Usage

### Basic Usage

```bash
# Show info for a specific file
finfo /usr/bin/python3

# Search for command in PATH
finfo python3

# Multiple files
finfo file1 file2 file3

# Use glob patterns
finfo *.so
```

### Advanced Features

```bash
# Show file with checksums
finfo --hash file.zip

# Compare two files (git-like diff)
finfo file1.txt file2.txt --diff

# Search for library files
finfo --lib ssl

# Disable colors
finfo --no-color /usr/bin/gcc

# Show only linked libraries (full list, no other info)
finfo --ll cmake
```

## Output Example

```
Path        : /usr/bin/python3
Size        : 31 KB
              ╰─ 31744 bytes
Arch        : arm64
OS          : Darwin
Permissions : -rwxr-xr-x
  ├─ Owner  : rwx (read, write, execute)
  ├─ Group  : rx (read, execute)
  ╰─ Others : rx (read, execute)
File Type   : Mach-O executable
MIME Type   : application/x-mach-binary
Binary Type : Executable
Linked Libraries:
  ├── /usr/lib/libSystem.B.dylib
  ├── /usr/lib/libc++.1.dylib
  ╰── /usr/lib/libz.1.dylib
Privileges:
  ├─ Owner        : root
  ├─ Writable by  : root only
  ╰─ Requires sudo: yes
Checksums: 
  ├─ MD5    : fe7a92c873699de5873ece5963176689
  ├─ SHA256 : 4970a8c688841ce9726a217c549b29e5f34bfdb53abfe221a06d78c903a05368
  ╰─ SHA512 : d14570aeaef900d09bf40d472c0d11813237863b8b98e863647fdd6eb35d7c9c...
```

## Flags

| Flag | Description |
|------|-------------|
| `--no-color` | Disable colored output |
| `--lib` | Search for library files (.so, .a, .dylib) |
| `--hash` | Calculate and show file checksums (MD5, SHA256, SHA512) |
| `--diff` | Compare two files and show differences |
| `--ll`, `--linked-libs` | Show only linked libraries (full list, no other info) |

## File Comparison

The `--diff` flag provides a git-like comparison of two files:

```bash
finfo file1 file2 --diff
```

Output includes:
- ✓ **Size comparison** - Byte-level difference
- ✓ **Permissions** - Mode comparison
- ✓ **Modification time** - Timestamp comparison
- ✓ **Checksums** - MD5, SHA256, SHA512 comparison
- ✓ **Final verdict** - IDENTICAL or DIFFERENT

## Color Scheme

- **Labels**: Cyan (bold)
- **Values**: White
- **Tree symbols**: Blue
- **Paths**: Green
- **Sizes**: Magenta
- **Permissions**: Yellow
- **Warnings**: Red
- **Executable type**: Green
- **Matches (✓)**: Green
- **Differences (✗)**: Red

## Development

### Building

```bash
# Build for current platform
just build

# Build for all platforms
just build-all

# Run tests
just test

# Install locally
just install
```

### Project Structure

```
finfo/
├── main.go              # Entry point
├── fileinfo.go          # Core file info gathering
├── fileinfo_darwin.go   # macOS-specific info
├── fileinfo_linux.go    # Linux-specific info
├── formatter.go         # Output formatting
├── filetype.go          # File type detection
├── binary.go            # Binary analysis
├── hash.go              # Hash calculation & comparison
├── resolver.go          # Command & library resolution
└── cmd/
    └── root.go          # CLI command definitions
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Author

Created by [@oh-tarnished](https://github.com/oh-tarnished)

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [fatih/color](https://github.com/fatih/color) for terminal colors
