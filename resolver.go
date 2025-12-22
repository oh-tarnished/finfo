package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveCommand tries to find a command in PATH
func ResolveCommand(name string) (string, error) {
	// First check if it's already a valid path
	if _, err := os.Stat(name); err == nil {
		return filepath.Abs(name)
	}

	// Try to find it in PATH using 'which'
	cmd := exec.Command("which", name)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command '%s' not found in PATH", name)
	}

	path := strings.TrimSpace(string(output))
	return path, nil
}

// FindLibrary searches for library files (.so, .a, .dylib) in common library paths
func FindLibrary(name string) ([]string, error) {
	// Common library search paths
	searchPaths := []string{
		"/usr/lib",
		"/usr/local/lib",
		"/opt/homebrew/lib",
		"/lib",
		"/usr/lib64",
		"/usr/local/lib64",
	}

	// Add LD_LIBRARY_PATH if set
	if ldPath := os.Getenv("LD_LIBRARY_PATH"); ldPath != "" {
		paths := strings.Split(ldPath, ":")
		searchPaths = append(searchPaths, paths...)
	}

	// Add DYLD_LIBRARY_PATH for macOS
	if dyldPath := os.Getenv("DYLD_LIBRARY_PATH"); dyldPath != "" {
		paths := strings.Split(dyldPath, ":")
		searchPaths = append(searchPaths, paths...)
	}

	var found []string
	extensions := []string{".dylib", ".so", ".a"}

	// Normalize library name (remove lib prefix and extension if present)
	libName := strings.TrimPrefix(name, "lib")
	for _, ext := range extensions {
		libName = strings.TrimSuffix(libName, ext)
	}

	// Search for the library
	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); os.IsNotExist(err) {
			continue
		}

		// Try different naming patterns
		patterns := []string{
			fmt.Sprintf("lib%s.dylib", libName),
			fmt.Sprintf("lib%s.so", libName),
			fmt.Sprintf("lib%s.so.*", libName),
			fmt.Sprintf("lib%s.a", libName),
			fmt.Sprintf("%s.dylib", libName),
			fmt.Sprintf("%s.so", libName),
			fmt.Sprintf("%s.a", libName),
		}

		for _, pattern := range patterns {
			matches, err := filepath.Glob(filepath.Join(searchPath, pattern))
			if err != nil {
				continue
			}
			found = append(found, matches...)
		}
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("library '%s' not found in standard paths", name)
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, path := range found {
		if !seen[path] {
			seen[path] = true
			unique = append(unique, path)
		}
	}

	return unique, nil
}
