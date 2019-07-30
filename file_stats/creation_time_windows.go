// +build windows

package file_stats

import (
	"os"
	"time"

	times "gopkg.in/djherbis/times.v1"
)

// GetCreationTime returns creation time
func GetCreationTime(filename string) time.Time {
	fileStats, err := times.Stat(filename)
	if err != nil {
		panic(err)
	}

	return fileStats.BirthTime()
}

// GetSize returns file size
func GetSize(filename string) uint64 {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		panic(err)
	}

	return uint64(fileStats.Size())
}
