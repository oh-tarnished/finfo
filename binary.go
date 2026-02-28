package main

import (
	"debug/elf"
	"debug/macho"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// BinaryInfo contains binary analysis information
type BinaryInfo struct {
	IsExecutable    bool
	LinkedLibraries []string
	IsStripped      bool
	HasSignature    bool
	SignatureInfo   string
}

// AnalyzeBinary analyzes binary files for dependencies and properties
func AnalyzeBinary(path string, fileType *FileTypeInfo) (*BinaryInfo, error) {
	info := &BinaryInfo{}

	if !fileType.IsBinary {
		return info, nil
	}

	// Detect if executable
	info.IsExecutable = strings.Contains(fileType.FileFormat, "executable")

	// Platform-specific analysis
	switch runtime.GOOS {
	case "darwin":
		_ = analyzeMachO(path, info)
	case "linux":
		_ = analyzeELF(path, info)
	}

	return info, nil
}

// analyzeMachO analyzes Mach-O binaries (macOS)
func analyzeMachO(path string, info *BinaryInfo) error {
	file, err := macho.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// Get linked libraries using otool
	cmd := exec.Command("otool", "-L", path)
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i == 0 {
				continue // Skip the first line (file path)
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Extract library path (before the version info in parentheses)
			if idx := strings.Index(line, "("); idx > 0 {
				lib := strings.TrimSpace(line[:idx])
				info.LinkedLibraries = append(info.LinkedLibraries, lib)
			}
		}
	}

	// Check if stripped using nm
	cmd = exec.Command("nm", path)
	if err := cmd.Run(); err != nil {
		// If nm fails, binary is likely stripped
		info.IsStripped = true
	}

	// Check code signature
	cmd = exec.Command("codesign", "-dv", path)
	output, err = cmd.CombinedOutput()
	if err == nil && len(output) > 0 {
		info.HasSignature = true
		// Parse signature info
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Authority=") {
				info.SignatureInfo = strings.TrimSpace(strings.TrimPrefix(line, "Authority="))
				break
			}
		}
	}

	return nil
}

// analyzeELF analyzes ELF binaries (Linux)
func analyzeELF(path string, info *BinaryInfo) error {
	file, err := elf.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// Get linked libraries
	libs, err := file.ImportedLibraries()
	if err == nil {
		info.LinkedLibraries = libs
	}

	// Check if stripped
	symbols, err := file.Symbols()
	if err != nil || len(symbols) == 0 {
		info.IsStripped = true
	}

	// Get dynamic libraries using ldd (if available)
	cmd := exec.Command("ldd", path)
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		var lddLibs []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.Contains(line, "not a dynamic executable") {
				continue
			}
			// Parse ldd output: "libname.so => /path/to/lib (address)"
			parts := strings.Fields(line)
			if len(parts) >= 3 && parts[1] == "=>" {
				lddLibs = append(lddLibs, parts[2])
			} else if len(parts) >= 1 {
				lddLibs = append(lddLibs, parts[0])
			}
		}
		if len(lddLibs) > 0 {
			info.LinkedLibraries = lddLibs
		}
	}

	return nil
}

// FormatBinaryInfo formats binary information for display with colors
func FormatBinaryInfo(info *BinaryInfo, showFullLinkedLibs bool, labelFn, valueFn, treeFn, execFn func(a ...interface{}) string) string {
	if !info.IsExecutable && len(info.LinkedLibraries) == 0 {
		return ""
	}

	var sb strings.Builder

	if info.IsExecutable {
		sb.WriteString(labelFn("Binary Type : "))
		sb.WriteString(execFn("Executable"))
		if info.IsStripped {
			sb.WriteString(valueFn(" (stripped)"))
		}
		sb.WriteString("\n")
	}

	// Only show code signature if it exists and has info
	if info.HasSignature && info.SignatureInfo != "" {
		sb.WriteString(labelFn("Code Signature: "))
		sb.WriteString(valueFn(info.SignatureInfo))
		sb.WriteString("\n")
	}

	if len(info.LinkedLibraries) > 0 {
		sb.WriteString(labelFn("Linked Libraries:\n"))
		limit := 10
		if showFullLinkedLibs || len(info.LinkedLibraries) < limit {
			limit = len(info.LinkedLibraries)
		}
		for i := 0; i < limit; i++ {
			if i == limit-1 && len(info.LinkedLibraries) <= limit {
				sb.WriteString(fmt.Sprintf("  %s %s\n",
					treeFn("╰──"),
					valueFn(info.LinkedLibraries[i])))
			} else {
				sb.WriteString(fmt.Sprintf("  %s %s\n",
					treeFn("├──"),
					valueFn(info.LinkedLibraries[i])))
			}
		}
		if !showFullLinkedLibs && len(info.LinkedLibraries) > limit {
			sb.WriteString(fmt.Sprintf("  %s %s\n",
				treeFn("╰──"),
				valueFn(fmt.Sprintf("... and %d more (use --ll for full list)", len(info.LinkedLibraries)-limit))))
		}
	}

	return sb.String()
}

// FormatLinkedLibrariesOnlySection outputs only the linked libraries section (full list, no truncation)
func FormatLinkedLibrariesOnlySection(info *BinaryInfo, labelFn, valueFn, treeFn func(a ...interface{}) string) string {
	if len(info.LinkedLibraries) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(labelFn("Linked Libraries:\n"))
	for i, lib := range info.LinkedLibraries {
		if i == len(info.LinkedLibraries)-1 {
			sb.WriteString(fmt.Sprintf("  %s %s\n", treeFn("╰──"), valueFn(lib)))
		} else {
			sb.WriteString(fmt.Sprintf("  %s %s\n", treeFn("├──"), valueFn(lib)))
		}
	}
	return sb.String()
}
