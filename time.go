package main

import (
	"fmt"
	"strings"
	"time"
)

func DurationIntFormat(duration time.Duration) string {
	n := duration.Nanoseconds()
	micros := 0
	millis := 0
	seconds := 0
	minutes := 0
	hours := 0
	output := ""
	for n > int64(time.Hour) {
		hours++
		n -= int64(time.Hour)
	}
	for n > int64(time.Minute) {
		minutes++
		n -= int64(time.Minute)
	}
	for n > int64(time.Second) {
		seconds++
		n -= int64(time.Second)
	}
	for n > int64(time.Millisecond) {
		millis++
		n -= int64(time.Millisecond)
	}
	for n > int64(time.Microsecond) {
		micros++
		n -= int64(time.Microsecond)
	}
	if hours != 0 {
		output += fmt.Sprint(hours) + "h "
	}
	if minutes != 0 {
		output += fmt.Sprint(minutes) + "m "
	}
	if seconds != 0 {
		output += fmt.Sprint(seconds) + "s "
	}
	if millis != 0 {
		output += fmt.Sprint(millis) + "ms "
	}
	if micros != 0 {
		output += fmt.Sprint(micros) + "µs "
	}
	if n != 0 {
		output += fmt.Sprint(n) + "ns "
	}
	output = strings.Trim(output, " ") // remove space if last rune is " "
	return output
}
