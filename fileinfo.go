package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var CalculateHashesFlag bool

// FileInfo contains detailed information about a file
type FileInfo struct {
	Path            string
	Size            int64
	Permissions     string
	Owner           string
	Group           string
	IsWritableByAll bool
	RequiresSudo    bool
	Arch            string
	OS              string
	SymlinkChain    []string
	FileType        *FileTypeInfo
	BinaryInfo      *BinaryInfo
	HashInfo        *HashInfo
}

// GetFileInfo retrieves comprehensive file information
func GetFileInfo(path string) (*FileInfo, error) {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Get file info
	info, err := os.Lstat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	fi := &FileInfo{
		Path:        absPath,
		Size:        info.Size(),
		Permissions: info.Mode().String(),
	}

	// Resolve symlink chain
	fi.SymlinkChain, err = resolveSymlinkChain(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve symlink chain: %w", err)
	}

	// Get platform-specific information
	if err := getPlatformSpecificInfo(fi, info); err != nil {
		return nil, fmt.Errorf("failed to get platform-specific info: %w", err)
	}

	// Detect file type
	fileType, err := DetectFileType(absPath)
	if err == nil {
		fi.FileType = fileType
	}

	// Analyze binary if applicable
	if fi.FileType != nil && fi.FileType.IsBinary {
		binaryInfo, err := AnalyzeBinary(absPath, fi.FileType)
		if err == nil {
			fi.BinaryInfo = binaryInfo
		}
	}

	// Calculate hashes if requested
	if CalculateHashesFlag {
		hashInfo, err := CalculateHashes(absPath)
		if err == nil {
			fi.HashInfo = hashInfo
		}
	}

	return fi, nil
}

// resolveSymlinkChain follows symlinks and returns the chain
func resolveSymlinkChain(path string) ([]string, error) {
	chain := []string{}
	current := path

	for {
		info, err := os.Lstat(current)
		if err != nil {
			return nil, err
		}

		if info.Mode()&os.ModeSymlink == 0 {
			// Not a symlink, we're done
			break
		}

		// Read the symlink
		target, err := os.Readlink(current)
		if err != nil {
			return nil, err
		}

		// If target is relative, make it absolute
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(current), target)
		}

		chain = append(chain, fmt.Sprintf("%s → %s", current, target))
		current = target
	}

	return chain, nil
}
