package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Color definitions
var (
	labelColor   = color.New(color.FgCyan, color.Bold)
	valueColor   = color.New(color.FgWhite)
	treeColor    = color.New(color.FgBlue)
	pathColor    = color.New(color.FgGreen)
	sizeColor    = color.New(color.FgMagenta)
	permColor    = color.New(color.FgYellow)
	symlinkColor = color.New(color.FgCyan)
	warnColor    = color.New(color.FgRed)
	execColor    = color.New(color.FgGreen)
)

// DisableColors disables color output
var DisableColors = false

// FormatLinkedLibrariesOnly outputs only the linked libraries section with full list
func FormatLinkedLibrariesOnly(fi *FileInfo) string {
	if DisableColors {
		color.NoColor = true
	}
	if fi.BinaryInfo == nil || len(fi.BinaryInfo.LinkedLibraries) == 0 {
		return fmt.Sprintf("No linked libraries found for %s\n", fi.Path)
	}
	return FormatLinkedLibrariesOnlySection(fi.BinaryInfo, labelColor.Sprint, valueColor.Sprint, treeColor.Sprint)
}

// FormatFileInfo formats the file information into a tree-like display
func FormatFileInfo(fi *FileInfo) string {
	if DisableColors {
		color.NoColor = true
	}

	var sb strings.Builder

	// Path
	sb.WriteString(labelColor.Sprint("Path        : "))
	sb.WriteString(pathColor.Sprintf("%s\n", fi.Path))

	// Size with tree structure - intelligently show only relevant units
	sb.WriteString(labelColor.Sprint("Size        : "))

	sizeGB := float64(fi.Size) / 1024 / 1024 / 1024
	sizeMB := float64(fi.Size) / 1024 / 1024
	sizeKB := float64(fi.Size) / 1024

	// Determine which units to show based on size
	var sizeLines []string

	if sizeGB >= 0.1 {
		// Show GB if >= 100 MB
		sizeLines = append(sizeLines, sizeColor.Sprintf("%.2f GB", sizeGB))
	}
	if sizeMB >= 1 {
		// Show MB if >= 1 MB
		sizeLines = append(sizeLines, sizeColor.Sprintf("%.0f MB", sizeMB))
	}
	if sizeKB >= 1 && sizeMB < 1000 {
		// Show KB if >= 1 KB and < 1000 MB
		sizeLines = append(sizeLines, sizeColor.Sprintf("%.0f KB", sizeKB))
	}
	// Always show bytes
	sizeLines = append(sizeLines, sizeColor.Sprintf("%d bytes", fi.Size))

	// Print the first line (primary size)
	sb.WriteString(sizeLines[0])
	sb.WriteString("\n")

	// Print remaining lines with tree structure
	for i := 1; i < len(sizeLines); i++ {
		if i == len(sizeLines)-1 {
			fmt.Fprintf(&sb, "              %s ", treeColor.Sprint("╰─"))
		} else {
			fmt.Fprintf(&sb, "              %s ", treeColor.Sprint("├─"))
		}
		sb.WriteString(sizeLines[i])
		sb.WriteString("\n")
	}

	// Architecture
	sb.WriteString(labelColor.Sprint("Arch        : "))
	sb.WriteString(valueColor.Sprintf("%s\n", fi.Arch))

	// OS
	sb.WriteString(labelColor.Sprint("OS          : "))
	sb.WriteString(valueColor.Sprintf("%s\n", fi.OS))

	// Permissions
	sb.WriteString(labelColor.Sprint("Permissions : "))
	sb.WriteString(permColor.Sprintf("%s\n", fi.Permissions))
	perms := parsePermissions(fi.Permissions)
	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeColor.Sprint("├─"),
		treeColor.Sprint("Owner  :"),
		valueColor.Sprint(perms.Owner))
	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeColor.Sprint("├─"),
		treeColor.Sprint("Group  :"),
		valueColor.Sprint(perms.Group))
	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeColor.Sprint("╰─"),
		treeColor.Sprint("Others :"),
		valueColor.Sprint(perms.Others))

	// File Type Information
	if fi.FileType != nil {
		sb.WriteString(labelColor.Sprint("File Type   : "))
		sb.WriteString(valueColor.Sprintf("%s\n", fi.FileType.FileFormat))

		if fi.FileType.MIMEType != "" {
			sb.WriteString(labelColor.Sprint("MIME Type   : "))
			sb.WriteString(valueColor.Sprintf("%s\n", fi.FileType.MIMEType))
		}

		if fi.FileType.IsScript && fi.FileType.Interpreter != "" {
			sb.WriteString(labelColor.Sprint("Interpreter : "))
			sb.WriteString(valueColor.Sprintf("%s\n", fi.FileType.Interpreter))
		}

		if fi.FileType.IsText && fi.FileType.Encoding != "" {
			sb.WriteString(labelColor.Sprint("Encoding    : "))
			sb.WriteString(valueColor.Sprintf("%s\n", fi.FileType.Encoding))
		}
	}

	// Binary Information
	if fi.BinaryInfo != nil {
		binaryOutput := FormatBinaryInfo(fi.BinaryInfo, false, labelColor.Sprint, valueColor.Sprint, treeColor.Sprint, execColor.Sprint)
		if binaryOutput != "" {
			sb.WriteString(binaryOutput)
		}
	}

	// Privileges section - header in blue, labels in blue, values in white, warnings in red
	sb.WriteString(labelColor.Sprint("Privileges:\n"))
	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeColor.Sprint("├─"),
		treeColor.Sprint("Owner        :"),
		valueColor.Sprint(fi.Owner))

	writableBy := fi.Owner + " only"
	writableColor := valueColor
	if fi.IsWritableByAll {
		writableBy = "all users"
		writableColor = warnColor
	}
	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeColor.Sprint("├─"),
		treeColor.Sprint("Writable by  :"),
		writableColor.Sprint(writableBy))

	requiresSudo := "no"
	sudoColor := valueColor
	if fi.RequiresSudo {
		requiresSudo = "yes"
		sudoColor = warnColor
	}
	fmt.Fprintf(&sb, "  %s %s %s\n",
		treeColor.Sprint("╰─"),
		treeColor.Sprint("Requires sudo:"),
		sudoColor.Sprint(requiresSudo))

	// Hash Information
	if fi.HashInfo != nil {
		sb.WriteString(labelColor.Sprint("Checksums: \n"))
		sb.WriteString(FormatHashInfo(fi.HashInfo, treeColor.Sprint, treeColor.Sprint, sizeColor.Sprint))
	}

	// Symlink chain (if exists)
	if len(fi.SymlinkChain) > 0 {
		sb.WriteString(labelColor.Sprint("Symlink chain:\n"))
		for i, link := range fi.SymlinkChain {
			// Split the link into source and target
			parts := strings.Split(link, " → ")
			if len(parts) == 2 {
				if i == len(fi.SymlinkChain)-1 {
					fmt.Fprintf(&sb, "  %s ", treeColor.Sprint("╰──"))
					sb.WriteString(pathColor.Sprint(parts[0]))
					sb.WriteString(sizeColor.Sprint(" → "))
					sb.WriteString(sizeColor.Sprint(parts[1]))
				} else {
					fmt.Fprintf(&sb, "  %s ", treeColor.Sprint("├──"))
					sb.WriteString(pathColor.Sprint(parts[0]))
					sb.WriteString(sizeColor.Sprint(" → "))
					sb.WriteString(sizeColor.Sprintln(parts[1]))
				}
			} else {
				// Fallback if format is unexpected
				if i == len(fi.SymlinkChain)-1 {
					fmt.Fprintf(&sb, "  %s %s", treeColor.Sprint("╰──"), symlinkColor.Sprint(link))
				} else {
					fmt.Fprintf(&sb, "  %s %s\n", treeColor.Sprint("├──"), symlinkColor.Sprint(link))
				}
			}
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

// PermissionBreakdown represents the breakdown of permissions
type PermissionBreakdown struct {
	Owner  string
	Group  string
	Others string
}

// parsePermissions parses the permission string into a readable format
func parsePermissions(permStr string) PermissionBreakdown {
	// Permission string format: -rwxr-xr-x or similar
	// First character is file type, then 3 groups of 3 characters each

	pb := PermissionBreakdown{}

	if len(permStr) < 10 {
		return pb
	}

	// Owner permissions (chars 1-3)
	ownerPerms := permStr[1:4]
	pb.Owner = formatPermGroup(ownerPerms)

	// Group permissions (chars 4-6)
	groupPerms := permStr[4:7]
	pb.Group = formatPermGroup(groupPerms)

	// Others permissions (chars 7-9)
	othersPerms := permStr[7:10]
	pb.Others = formatPermGroup(othersPerms)

	return pb
}

// formatPermGroup formats a 3-character permission group into readable text
func formatPermGroup(perms string) string {
	var parts []string
	var short []string

	if perms[0] == 'r' {
		parts = append(parts, "read")
		short = append(short, "r")
	}
	if perms[1] == 'w' {
		parts = append(parts, "write")
		short = append(short, "w")
	}
	if perms[2] == 'x' {
		parts = append(parts, "execute")
		short = append(short, "x")
	}

	if len(parts) == 0 {
		return "--- (no permissions)"
	}

	shortStr := strings.Join(short, "")
	longStr := strings.Join(parts, ", ")
	return fmt.Sprintf("%s (%s)", shortStr, longStr)
}
