// +build windows

package creation_time

import (
	"time"

	times "gopkg.in/djherbis/times.v1"
)

func Get_creation_time(filename string) time.Time {
	file_stats, err := times.Stat(filename)
	if err != nil {
		panic(err)
	}

	return file_stats.BirthTime()
}
