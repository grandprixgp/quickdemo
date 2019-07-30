// +build windows

package memory_stats

import (
	"syscall"
	"unsafe"
)

type tMemStatusEx struct {
	dwLength     uint32
	unused0      [1]uint32
	unused1      [1]uint64
	ullAvailPhys uint64
	unused2      [5]uint64
}

// GetMemoryAvailable Returns currently available system memory in MB
func GetMemoryAvailable() uint64 {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return 0
	}

	globalMemoryStatusEx, err := kernel32.FindProc("GlobalMemoryStatusEx")
	if err != nil {
		return 0
	}

	memStatusEx := &tMemStatusEx{dwLength: 64}

	retval, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(memStatusEx)))
	if retval == 0 {
		return 0
	}

	return memStatusEx.ullAvailPhys / 1024 / 1024
}
