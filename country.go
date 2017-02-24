package tzlib

import "sort"

const countriesCap = 4

// Country represents a country
type Country struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	Cities []*City `json:"cities"`
	parent *Timezone
}

// NewCountry creates a new Country and returns a pointer to it
func NewCountry(name string, weight float64, parent *Timezone) *Country {
	return &Country{name, weight, make([]*City, 0, citiesCap), parent}
}

// AddNewCity creates a new city from name
// and appends it to this country
// and returns a pointer to the newly added city
func (c *Country) AddNewCity(name string, weight float64) *City {
	return c.AddCity(NewCity(name, weight, c))
}

// AddCity appends y to this country
// and returns a pointer to the newly added city (i.e. &y)
func (c *Country) AddCity(y *City) *City {
	y.parent = c

	p := sort.Search(len(c.Cities), func(i int) bool { return y.Weight > c.Cities[i].Weight })

	c.Cities = append(c.Cities, y)
	copy(c.Cities[p+1:], c.Cities[p:])
	c.Cities[p] = y

	return c.Cities[p]
}

// Offset returns the UTC offset, in seconds, of this country
func (c Country) Offset() int {
	if c.parent == nil {
		return -1
	}

	return c.parent.Offset
}
