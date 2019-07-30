// +build linux

package memory_stats

import (
	"syscall"
)

// GetMemoryAvailable Returns currently available system memory in MB
func GetMemoryAvailable() uint64 {
	sysinfo := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(sysinfo)
	if err != nil {
		return 0
	}
	return (uint64(sysinfo.Freeram) * uint64(sysinfo.Unit)) / 1024 / 1024
}
