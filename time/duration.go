package time

import (
	"time"
	"fmt"
)

const Day = time.Hour * 24

func Duration(d time.Duration, num int) time.Duration {
	bb := d.Nanoseconds() * int64(num)
	duration, _ := time.ParseDuration(fmt.Sprintf("%dns", bb))
	return duration
}

func Days(num int) time.Duration {
	return Duration(Day, num)
}