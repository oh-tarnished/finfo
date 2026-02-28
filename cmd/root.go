/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// GetFileInfoFunc is a function type for getting file info
var GetFileInfoFunc func(string) (interface{}, error)

// FormatFileInfoFunc is a function type for formatting file info
var FormatFileInfoFunc func(interface{}) string

// FormatLinkedLibrariesOnlyFunc formats only the linked libraries section (full list)
var FormatLinkedLibrariesOnlyFunc func(interface{}) string

// SetDisableColorsFunc is a function type for setting color disable flag
var SetDisableColorsFunc func(bool)

// ResolveCommandFunc is a function type for resolving commands in PATH
var ResolveCommandFunc func(string) (string, error)

// FindLibraryFunc is a function type for finding libraries
var FindLibraryFunc func(string) ([]string, error)

// CompareFilesFunc is a function type for comparing files
var CompareFilesFunc func(string, string) (string, error)

// SetCalculateHashesFunc is a function type for setting hash calculation flag
var SetCalculateHashesFunc func(bool)

var noColor bool
var searchLib bool
var showHash bool
var diffMode bool
var showFullLinkedLibs bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "finfo [file path or command...]",
	Short: "File info at a glance — metadata, checksums, binary analysis",
	Long: `finfo is a CLI tool that displays comprehensive information about files including:
- File size in multiple units (GB, MB, KB, bytes)
- Architecture and OS
- Permissions and ownership
- Symlink chain (if applicable)

If the argument is not a valid path, finfo will search for it in PATH.

Examples:
  finfo /usr/bin/python3        # Show info for a specific file
  finfo python3                 # Search for 'python3' in PATH
  finfo file1 file2 file3       # Show info for multiple files
  finfo *.so                    # Use glob patterns
  finfo --lib ssl               # Search for SSL library files
  finfo --hash file.zip         # Show file with checksums
  finfo --diff file1 file2      # Compare two files
  finfo --ll cmake              # Show only linked libraries (full list)`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Set color preference
		if SetDisableColorsFunc != nil {
			SetDisableColorsFunc(noColor)
		}

		// Set hash calculation flag
		if SetCalculateHashesFunc != nil {
			SetCalculateHashesFunc(showHash)
		}

		// Handle diff mode
		if diffMode {
			if len(args) != 2 {
				fmt.Fprintf(os.Stderr, "Error: --diff requires exactly 2 file arguments\n")
				os.Exit(1)
			}
			if CompareFilesFunc != nil {
				output, err := CompareFilesFunc(args[0], args[1])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error comparing files: %v\n", err)
					os.Exit(1)
				}
				fmt.Print(output)
			}
			return
		}

		// Handle linked libraries only mode (--ll)
		if showFullLinkedLibs {
			for _, input := range args {
				filePath := input
				if _, err := os.Stat(input); os.IsNotExist(err) {
					if ResolveCommandFunc != nil {
						resolved, resolveErr := ResolveCommandFunc(input)
						if resolveErr == nil {
							filePath = resolved
						} else {
							fmt.Fprintf(os.Stderr, "Error: %v\n", resolveErr)
							continue
						}
					}
				}
				info, err := GetFileInfoFunc(filePath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", filePath, err)
					continue
				}
				if FormatLinkedLibrariesOnlyFunc != nil {
					output := FormatLinkedLibrariesOnlyFunc(info)
					fmt.Print(output)
				}
				if len(args) > 1 {
					fmt.Println()
				}
			}
			return
		}

		// Handle library search mode
		if searchLib {
			if FindLibraryFunc == nil {
				fmt.Fprintf(os.Stderr, "Error: Library search not available\n")
				os.Exit(1)
			}

			input := args[0]
			libraries, err := FindLibraryFunc(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Found %d library file(s) for '%s':\n\n", len(libraries), input)
			for i, lib := range libraries {
				info, err := GetFileInfoFunc(lib)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", lib, err)
					continue
				}
				output := FormatFileInfoFunc(info)
				fmt.Print(output)

				// Add separator between results (but not after the last one)
				if i < len(libraries)-1 {
					fmt.Println("\n" + strings.Repeat("─", 80))
				}
				fmt.Println()
			}
			return
		}

		// Handle multiple files
		var filePaths []string
		for _, input := range args {
			// Try to resolve as command if not a valid path
			filePath := input
			if _, err := os.Stat(input); os.IsNotExist(err) {
				if ResolveCommandFunc != nil {
					resolved, resolveErr := ResolveCommandFunc(input)
					if resolveErr == nil {
						filePath = resolved
					} else {
						fmt.Fprintf(os.Stderr, "Error: %v\n", resolveErr)
						continue
					}
				}
			}
			filePaths = append(filePaths, filePath)
		}

		if len(filePaths) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No valid files found\n")
			os.Exit(1)
		}

		// Display info for each file
		for i, filePath := range filePaths {
			info, err := GetFileInfoFunc(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", filePath, err)
				continue
			}

			output := FormatFileInfoFunc(info)
			fmt.Print(output)

			// Add separator between files (but not after the last one)
			if i < len(filePaths)-1 {
				fmt.Println(strings.Repeat("─", 80))
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.Flags().BoolVar(&searchLib, "lib", false, "Search for library files (.so, .a, .dylib)")
	rootCmd.Flags().BoolVar(&showHash, "hash", false, "Calculate and show file checksums (MD5, SHA256, SHA512)")
	rootCmd.Flags().BoolVar(&diffMode, "diff", false, "Compare two files and show differences")
	rootCmd.Flags().BoolVar(&showFullLinkedLibs, "ll", false, "Show only linked libraries (full list, no other info)")
	rootCmd.Flags().BoolVar(&showFullLinkedLibs, "linked-libs", false, "Alias for --ll")
}
