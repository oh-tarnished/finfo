package main

import (
	"github.com/oh-tarnished/finfo/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersion(version)
	// Set the function pointers for cmd package to use
	cmd.GetFileInfoFunc = func(path string) (interface{}, error) {
		return GetFileInfo(path)
	}
	cmd.FormatFileInfoFunc = func(info interface{}) string {
		return FormatFileInfo(info.(*FileInfo))
	}
	cmd.FormatLinkedLibrariesOnlyFunc = func(info interface{}) string {
		return FormatLinkedLibrariesOnly(info.(*FileInfo))
	}
	cmd.SetDisableColorsFunc = func(disable bool) {
		DisableColors = disable
	}
	cmd.ResolveCommandFunc = func(name string) (string, error) {
		return ResolveCommand(name)
	}
	cmd.FindLibraryFunc = func(name string) ([]string, error) {
		return FindLibrary(name)
	}
	cmd.SetCalculateHashesFunc = func(enable bool) {
		CalculateHashesFlag = enable
	}
	cmd.CompareFilesFunc = func(path1, path2 string) (string, error) {
		// Import color package functions
		labelFn := func(a ...interface{}) string { return labelColor.Sprint(a...) }
		treeFn := func(a ...interface{}) string { return treeColor.Sprint(a...) }
		matchFn := func(a ...interface{}) string { return pathColor.Sprint(a...) } // Green for matches
		diffFn := func(a ...interface{}) string { return warnColor.Sprint(a...) }  // Red for differences
		valueFn := func(a ...interface{}) string { return valueColor.Sprint(a...) }

		return CompareFiles(path1, path2, labelFn, treeFn, matchFn, diffFn, valueFn)
	}

	cmd.Execute()
}
