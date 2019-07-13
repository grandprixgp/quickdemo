// +build linux

package creation_time

import (
	"os"
	"syscall"
	"time"
)

func Get_creation_time(filename string) time.Time {
	demo_file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer demo_file.Close()
	file_stats, err := demo_file.Stat()
	file_stat := file_stats.Sys().(*syscall.Stat_t)

	return time.Unix(int64(file_stat.Ctim.Sec), int64(file_stat.Ctim.Nsec))
}
