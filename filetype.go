package main

import (
	"bufio"
	"bytes"
	"debug/elf"
	"debug/macho"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// FileTypeInfo contains file type detection information
type FileTypeInfo struct {
	MIMEType    string
	FileFormat  string
	IsText      bool
	IsBinary    bool
	IsScript    bool
	Interpreter string
	Encoding    string
}

// DetectFileType detects the type and format of a file
func DetectFileType(path string) (*FileTypeInfo, error) {
	info := &FileTypeInfo{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	// Read first 512 bytes for magic number detection
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && n == 0 {
		return nil, err
	}
	header = header[:n]

	// Check for binary formats
	if detectMachO(path, info) {
		return info, nil
	}
	if detectELF(path, info) {
		return info, nil
	}

	// Check for script
	if detectScript(header, info) {
		return info, nil
	}

	// Check if text or binary
	if isText(header) {
		info.IsText = true
		info.FileFormat = "Text file"
		info.MIMEType = "text/plain"
		info.Encoding = detectEncoding(header)
	} else {
		info.IsBinary = true
		info.FileFormat = "Binary file"
		info.MIMEType = "application/octet-stream"
	}

	// Detect by extension
	detectByExtension(path, info)

	return info, nil
}

// detectMachO checks if file is a Mach-O binary (macOS)
func detectMachO(path string, info *FileTypeInfo) bool {
	file, err := macho.Open(path)
	if err != nil {
		return false
	}
	defer func() { _ = file.Close() }()

	info.IsBinary = true
	info.FileFormat = "Mach-O executable"
	info.MIMEType = "application/x-mach-binary"

	switch file.Type {
	case macho.TypeExec:
		info.FileFormat = "Mach-O executable"
	case macho.TypeDylib:
		info.FileFormat = "Mach-O dynamic library"
	case macho.TypeBundle:
		info.FileFormat = "Mach-O bundle"
	case macho.TypeObj:
		info.FileFormat = "Mach-O object file"
	}

	return true
}

// detectELF checks if file is an ELF binary (Linux)
func detectELF(path string, info *FileTypeInfo) bool {
	file, err := elf.Open(path)
	if err != nil {
		return false
	}
	defer func() { _ = file.Close() }()

	info.IsBinary = true
	info.FileFormat = "ELF executable"
	info.MIMEType = "application/x-executable"

	switch file.Type {
	case elf.ET_EXEC:
		info.FileFormat = "ELF executable"
	case elf.ET_DYN:
		info.FileFormat = "ELF shared object"
	case elf.ET_REL:
		info.FileFormat = "ELF relocatable"
	case elf.ET_CORE:
		info.FileFormat = "ELF core dump"
	}

	return true
}

// detectScript checks if file is a script with shebang
func detectScript(header []byte, info *FileTypeInfo) bool {
	if len(header) < 2 || header[0] != '#' || header[1] != '!' {
		return false
	}

	info.IsScript = true
	info.IsText = true

	// Extract shebang line
	scanner := bufio.NewScanner(bytes.NewReader(header))
	if scanner.Scan() {
		shebang := scanner.Text()
		info.Interpreter = strings.TrimSpace(strings.TrimPrefix(shebang, "#!"))

		// Determine script type
		if strings.Contains(info.Interpreter, "bash") {
			info.FileFormat = "Bash script"
			info.MIMEType = "text/x-shellscript"
		} else if strings.Contains(info.Interpreter, "sh") {
			info.FileFormat = "Shell script"
			info.MIMEType = "text/x-shellscript"
		} else if strings.Contains(info.Interpreter, "python") {
			info.FileFormat = "Python script"
			info.MIMEType = "text/x-python"
		} else if strings.Contains(info.Interpreter, "perl") {
			info.FileFormat = "Perl script"
			info.MIMEType = "text/x-perl"
		} else if strings.Contains(info.Interpreter, "ruby") {
			info.FileFormat = "Ruby script"
			info.MIMEType = "text/x-ruby"
		} else if strings.Contains(info.Interpreter, "node") || strings.Contains(info.Interpreter, "nodejs") {
			info.FileFormat = "Node.js script"
			info.MIMEType = "text/javascript"
		} else {
			info.FileFormat = "Script"
			info.MIMEType = "text/plain"
		}
	}

	return true
}

// isText checks if content appears to be text
func isText(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	// Check for null bytes (strong indicator of binary)
	if bytes.IndexByte(data, 0) != -1 {
		return false
	}

	// Check if valid UTF-8
	if !utf8.Valid(data) {
		return false
	}

	// Count printable characters
	printable := 0
	for _, b := range data {
		if b >= 32 && b < 127 || b == '\n' || b == '\r' || b == '\t' {
			printable++
		}
	}

	// If more than 95% printable, consider it text
	return float64(printable)/float64(len(data)) > 0.95
}

// detectEncoding detects text encoding
func detectEncoding(data []byte) string {
	if utf8.Valid(data) {
		return "UTF-8"
	}
	return "Unknown"
}

// detectByExtension provides additional hints based on file extension
func detectByExtension(path string, info *FileTypeInfo) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".go":
		info.FileFormat = "Go source code"
		info.MIMEType = "text/x-go"
	case ".c":
		info.FileFormat = "C source code"
		info.MIMEType = "text/x-c"
	case ".cpp", ".cc", ".cxx":
		info.FileFormat = "C++ source code"
		info.MIMEType = "text/x-c++"
	case ".h", ".hpp":
		info.FileFormat = "C/C++ header"
		info.MIMEType = "text/x-c"
	case ".js":
		info.FileFormat = "JavaScript"
		info.MIMEType = "text/javascript"
	case ".json":
		info.FileFormat = "JSON"
		info.MIMEType = "application/json"
	case ".xml":
		info.FileFormat = "XML"
		info.MIMEType = "application/xml"
	case ".html", ".htm":
		info.FileFormat = "HTML"
		info.MIMEType = "text/html"
	case ".css":
		info.FileFormat = "CSS"
		info.MIMEType = "text/css"
	case ".md", ".markdown":
		info.FileFormat = "Markdown"
		info.MIMEType = "text/markdown"
	case ".yaml", ".yml":
		info.FileFormat = "YAML"
		info.MIMEType = "text/yaml"
	case ".toml":
		info.FileFormat = "TOML"
		info.MIMEType = "text/toml"
	case ".zip":
		info.FileFormat = "ZIP archive"
		info.MIMEType = "application/zip"
	case ".tar":
		info.FileFormat = "TAR archive"
		info.MIMEType = "application/x-tar"
	case ".gz":
		info.FileFormat = "GZIP compressed"
		info.MIMEType = "application/gzip"
	case ".pdf":
		info.FileFormat = "PDF document"
		info.MIMEType = "application/pdf"
	case ".jpg", ".jpeg":
		info.FileFormat = "JPEG image"
		info.MIMEType = "image/jpeg"
	case ".png":
		info.FileFormat = "PNG image"
		info.MIMEType = "image/png"
	case ".gif":
		info.FileFormat = "GIF image"
		info.MIMEType = "image/gif"
	}
}
