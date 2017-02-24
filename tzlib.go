package tzlib

import (
	"fmt"
	"math"
	"sort"
	"time"
)

const (
	after  = 0
	before = -1
)

// Importer specifies the importer interface
type Importer interface {
	Import() (*Tzlib, error)
}

// Exporter specifies the exporter interface
type Exporter interface {
	Export(l *Tzlib) error
}

// Tzlib contains timezones and metadata
type Tzlib struct {
	Timezones []*Timezone `json:"timezones"`
	Expires   time.Time   `json:"expires"`
	Created   time.Time   `json:"created"`
	offsets   map[int]int
}

// NewTzlib returns a new Tzlib with expiry date
// and returns a pointer to it
func NewTzlib(created, expires time.Time) *Tzlib {
	return &Tzlib{
		Timezones: make([]*Timezone, 0, timezonesCap),
		Expires:   expires.Truncate(time.Second).UTC(),
		Created:   created.Truncate(time.Second).UTC(),
		offsets:   make(map[int]int, timezonesCap),
	}
}

// AddNewTimezone creates a new timezone from utc
// and appends it to this tzlib
// and returns a pointer to the newly added timezone
func (l *Tzlib) AddNewTimezone(utc string) *Timezone {
	return l.AddTimezone(NewTimezone(utc, l))
}

// AddTimezone appends z to this tzlib
// and returns a pointer to the newly added timezone
func (l *Tzlib) AddTimezone(z *Timezone) *Timezone {
	z.parent = l

	l.offsets[z.Offset] = sort.Search(len(l.Timezones), func(i int) bool { return z.Offset < l.Timezones[i].Offset })

	l.Timezones = append(l.Timezones, z)
	copy(l.Timezones[l.offsets[z.Offset]+1:], l.Timezones[l.offsets[z.Offset]:])
	l.Timezones[l.offsets[z.Offset]] = z

	return l.Timezones[l.offsets[z.Offset]]
}

// Export exports via the specified Exporter
func (l *Tzlib) Export(e Exporter) error {
	return e.Export(l)
}

// WhereWasIt returns the timezone where t most recently passed
func (l Tzlib) WhereWasIt(t time.Time) (z *Timezone, err error) {
	return l.where(after, t)
}

// WhereWillItBe returns the timezone t will pass next
func (l Tzlib) WhereWillItBe(t time.Time) (z *Timezone, err error) {
	return l.where(before, t)
}

// where oh where
func (l Tzlib) where(d int, t time.Time) (z *Timezone, err error) {
	if 0 == len(l.offsets) || len(l.offsets) != len(l.Timezones) {
		return nil, fmt.Errorf("inconsistent data or no timezones loaded")
	}

	index := sort.Search(len(l.Timezones), l.tzSearcher(l.calcOffset(t))) + d

	if index < 0 {
		err = fmt.Errorf("not enough timezones loaded")
	}

	z = l.Timezones[index]

	return
}

// tzSearcher
func (l Tzlib) tzSearcher(offset int) func(int) bool {
	return func(i int) bool { return offset < l.Timezones[i].Offset }
}

// calcOffset
func (l Tzlib) calcOffset(t time.Time) (offset int) {
	then := secondsFromTime(t)
	offset = then - secondsFromTime(time.Now().UTC())

	if !l.offsetInRange(offset) {
		offset += int(math.Copysign(float64(24*time.Hour/time.Second), float64(then)))
	}

	return
}

// offsetInRange
func (l Tzlib) offsetInRange(offset int) bool {
	return len(l.Timezones) > 0 &&
		l.Timezones[0].Offset <= offset &&
		offset <= l.Timezones[len(l.Timezones)-1].Offset
}
