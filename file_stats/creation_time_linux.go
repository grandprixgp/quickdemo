// +build linux

package file_stats

import (
	"os"
	"syscall"
	"time"
)

// GetCreationTime returns creation time
func GetCreationTime(filename string) time.Time {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file_stats, err := file.Stat()
	file_stat := file_stats.Sys().(*syscall.Stat_t)

	return time.Unix(int64(file_stat.Ctim.Sec), int64(file_stat.Ctim.Nsec))
}

// GetSize returns file size
func GetSize(filename string) uint64 {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file_stats, err := file.Stat()
	file_stat := file_stats.Sys().(*syscall.Stat_t)

	return uint64(file_stat.Size)
}
