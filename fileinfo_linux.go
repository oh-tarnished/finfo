//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// getPlatformSpecificInfo retrieves Linux-specific file information
func getPlatformSpecificInfo(fi *FileInfo, info os.FileInfo) error {
	// Get architecture
	fi.Arch = runtime.GOARCH
	fi.OS = "Linux"

	// Get owner and group information
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("failed to get syscall.Stat_t")
	}

	// Get owner name
	ownerName, err := getUserName(stat.Uid)
	if err != nil {
		fi.Owner = fmt.Sprintf("uid:%d", stat.Uid)
	} else {
		fi.Owner = ownerName
	}

	// Get group name
	groupName, err := getGroupName(stat.Gid)
	if err != nil {
		fi.Group = fmt.Sprintf("gid:%d", stat.Gid)
	} else {
		fi.Group = groupName
	}

	// Check if writable by all
	mode := info.Mode()
	fi.IsWritableByAll = mode&0002 != 0

	// Check if requires sudo (owned by root and not writable by current user)
	currentUID := uint32(os.Getuid())
	fi.RequiresSudo = stat.Uid == 0 && currentUID != 0

	return nil
}

// getUserName gets the username for a given UID
func getUserName(uid uint32) (string, error) {
	cmd := exec.Command("id", "-un", fmt.Sprintf("%d", uid))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getGroupName gets the group name for a given GID
func getGroupName(gid uint32) (string, error) {
	cmd := exec.Command("id", "-gn", fmt.Sprintf("%d", gid))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
