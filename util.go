package tzlib

import (
	"log"
	"strings"
	"time"
)

// Import returns a new Tzlib imported via the specified Importer
func Import(i Importer) (l *Tzlib, err error) {
	return i.Import()
}

// Time returns a new Time with custom hours, minutes, seconds
func Time(h, m, s int) time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, now.Location())
}

// utcToOffset transforms "UTC+12:34" into an absolute utc offset in seconds
func utcToOffset(s string) int {
	if !strings.Contains(s, ":") {
		s += ":00"
	}

	dur, err := time.ParseDuration(strings.NewReplacer("UTC", "", ":", "h").Replace(s) + "m")
	if err != nil {
		log.Fatal(err)
	}

	return int(dur.Seconds())
}

// secondsFromTime
func secondsFromTime(t time.Time) int {
	return secondsFromHMS(t.Hour(), t.Minute(), t.Second())
}

// secondsFromHMS
func secondsFromHMS(h, m, s int) int {
	return s + 60*m + 60*60*h
}
