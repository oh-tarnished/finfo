package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// HashInfo contains file hash information
type HashInfo struct {
	MD5    string
	SHA256 string
	SHA512 string
}

// CalculateHashes calculates MD5, SHA256, and SHA512 hashes for a file
func CalculateHashes(path string) (*HashInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	md5Hash := md5.New()
	sha256Hash := sha256.New()
	sha512Hash := sha512.New()

	// Use MultiWriter to calculate all hashes in one pass
	multiWriter := io.MultiWriter(md5Hash, sha256Hash, sha512Hash)

	if _, err := io.Copy(multiWriter, file); err != nil {
		return nil, err
	}

	return &HashInfo{
		MD5:    hex.EncodeToString(md5Hash.Sum(nil)),
		SHA256: hex.EncodeToString(sha256Hash.Sum(nil)),
		SHA512: hex.EncodeToString(sha512Hash.Sum(nil)),
	}, nil
}

// FormatHashInfo formats hash information for display with colors
func FormatHashInfo(info *HashInfo, labelFn, treeFn, valueFn func(a ...interface{}) string) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeFn("├─"),
		labelFn("MD5    :"),
		valueFn(info.MD5))
	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeFn("├─"),
		labelFn("SHA256 :"),
		valueFn(info.SHA256))
	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeFn("╰─"),
		labelFn("SHA512 :"),
		valueFn(info.SHA512))

	return sb.String()
}

// CompareFiles compares two files and returns differences with git-like formatting
func CompareFiles(path1, path2 string, labelFn, treeFn, matchFn, diffFn, valueFn func(a ...interface{}) string) (string, error) {
	// Get file info for both
	info1, err := os.Stat(path1)
	if err != nil {
		return "", err
	}

	info2, err := os.Stat(path2)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "\n%s\n", labelFn("diff "+filepath.Base(path1)+" "+filepath.Base(path2)))
	fmt.Fprintf(&sb, "%s %s\n", diffFn("---"), valueFn(path1))
	fmt.Fprintf(&sb, "%s %s\n\n", matchFn("+++"), valueFn(path2))

	// Compare sizes
	sb.WriteString(labelFn("Size:\n"))
	if info1.Size() == info2.Size() {
		fmt.Fprintf(&sb, "  %s Both files are %s\n",
			matchFn("✓"),
			valueFn(fmt.Sprintf("%d bytes", info1.Size())))
	} else {
		fmt.Fprintf(&sb, "  %s File 1: %s\n", diffFn("✗"), valueFn(fmt.Sprintf("%d bytes", info1.Size())))
		fmt.Fprintf(&sb, "  %s File 2: %s\n", diffFn("✗"), valueFn(fmt.Sprintf("%d bytes", info2.Size())))
		diff := info1.Size() - info2.Size()
		if diff > 0 {
			fmt.Fprintf(&sb, "  %s\n", diffFn(fmt.Sprintf("Δ File 1 is %d bytes larger", diff)))
		} else {
			fmt.Fprintf(&sb, "  %s\n", diffFn(fmt.Sprintf("Δ File 2 is %d bytes larger", -diff)))
		}
	}

	// Compare permissions
	fmt.Fprintf(&sb, "\n%s\n", labelFn("Permissions:"))
	if info1.Mode() == info2.Mode() {
		fmt.Fprintf(&sb, "  %s Both files have %s\n", matchFn("✓"), valueFn(info1.Mode().String()))
	} else {
		fmt.Fprintf(&sb, "  %s File 1: %s\n", diffFn("✗"), valueFn(info1.Mode().String()))
		fmt.Fprintf(&sb, "  %s File 2: %s\n", diffFn("✗"), valueFn(info2.Mode().String()))
	}

	// Compare modification times
	fmt.Fprintf(&sb, "\n%s\n", labelFn("Modified:"))
	if info1.ModTime().Equal(info2.ModTime()) {
		fmt.Fprintf(&sb, "  %s Both at %s\n",
			matchFn("✓"),
			valueFn(info1.ModTime().Format("2006-01-02 15:04:05")))
	} else {
		fmt.Fprintf(&sb, "  %s File 1: %s\n", diffFn("✗"), valueFn(info1.ModTime().Format("2006-01-02 15:04:05")))
		fmt.Fprintf(&sb, "  %s File 2: %s\n", diffFn("✗"), valueFn(info2.ModTime().Format("2006-01-02 15:04:05")))
	}

	// Calculate and compare hashes
	fmt.Fprintf(&sb, "\n%s\n", labelFn("Checksums:"))
	hash1, err1 := CalculateHashes(path1)
	hash2, err2 := CalculateHashes(path2)

	if err1 == nil && err2 == nil {
		if hash1.MD5 == hash2.MD5 {
			fmt.Fprintf(&sb, "  %s MD5: %s\n", matchFn("✓"), valueFn(hash1.MD5))
		} else {
			fmt.Fprintf(&sb, "  %s MD5 File 1: %s\n", diffFn("✗"), valueFn(hash1.MD5))
			fmt.Fprintf(&sb, "  %s MD5 File 2: %s\n", diffFn("✗"), valueFn(hash2.MD5))
		}

		if hash1.SHA256 == hash2.SHA256 {
			fmt.Fprintf(&sb, "  %s SHA256: %s\n", matchFn("✓"), valueFn(hash1.SHA256))
		} else {
			fmt.Fprintf(&sb, "  %s SHA256 File 1: %s\n", diffFn("✗"), valueFn(hash1.SHA256))
			fmt.Fprintf(&sb, "  %s SHA256 File 2: %s\n", diffFn("✗"), valueFn(hash2.SHA256))
		}

		if hash1.SHA512 == hash2.SHA512 {
			fmt.Fprintf(&sb, "  %s SHA512: %s\n", matchFn("✓"), valueFn(hash1.SHA512))
		} else {
			fmt.Fprintf(&sb, "  %s SHA512 File 1: %s\n", diffFn("✗"), valueFn(hash1.SHA512))
			fmt.Fprintf(&sb, "  %s SHA512 File 2: %s\n", diffFn("✗"), valueFn(hash2.SHA512))
		}

		// Overall verdict
		fmt.Fprintf(&sb, "\n%s\n", labelFn("Verdict:"))
		if hash1.SHA256 == hash2.SHA256 && info1.Size() == info2.Size() {
			fmt.Fprintf(&sb, "  %s\n", matchFn("✓ Files are IDENTICAL (same content)"))
		} else if hash1.SHA256 != hash2.SHA256 && info1.Size() == info2.Size() {
			fmt.Fprintf(&sb, "  %s\n", diffFn("✗ Files are DIFFERENT (same size, different content)"))
		} else {
			fmt.Fprintf(&sb, "  %s\n", diffFn("✗ Files are DIFFERENT"))
		}
	}

	sb.WriteString("\n")
	return sb.String(), nil
}
