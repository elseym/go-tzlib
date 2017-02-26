package tzlib

import (
	"sort"
	"time"
)

const timezonesCap = 40

// Timezone holds information about Offset and Countries
type Timezone struct {
	Utc       string         `json:"utc"`
	Offset    int            `json:"offset"`
	Countries []*Country     `json:"countries"`
	Location  *time.Location `json:"-"`
	parent    *Tzlib
}

// NewTimezone creates a new timezone and returns a pointer to it
func NewTimezone(utc string, parent *Tzlib) *Timezone {
	offset := utcToOffset(utc)
	location := time.FixedZone(utc, offset)
	countries := make([]*Country, 0, countriesCap)

	return &Timezone{utc, offset, countries, location, parent}
}

// AddNewCountry creates a new country from id and name
// and appends it to this timezone
// and returns a pointer to the newly added country
func (z *Timezone) AddNewCountry(name string, weight float64) *Country {
	return z.AddCountry(NewCountry(name, weight, z))
}

// AddCountry appends c to this timezone
// and returns a pointer to the newly added country
func (z *Timezone) AddCountry(c *Country) *Country {
	c.parent = z

	p := sort.Search(len(z.Countries), func(i int) bool { return c.Weight > z.Countries[i].Weight })

	z.Countries = append(z.Countries, c)
	copy(z.Countries[p+1:], z.Countries[p:])
	z.Countries[p] = c

	return z.Countries[p]
}

// Localtime returns the local time in this timezone
func (z Timezone) Localtime() (t time.Time) {
	return time.Now().In(z.Location)
}

// Until returns the duration until t is passed in this timezone
func (z Timezone) Until(t time.Time) (d time.Duration) {
	l := z.Localtime()
	h := (t.Hour() - l.Hour()) * 60 * 60
	m := (t.Minute() - l.Minute()) * 60
	s := t.Second() - l.Second()

	return time.Duration((h + m + s) * int(time.Second))
}
