package tzlib

const citiesCap = 2

// City represents a city
type City struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	parent *Country
}

// NewCity creates a new city and returns a pointer to it
func NewCity(name string, weight float64, parent *Country) *City {
	return &City{name, weight, parent}
}

// Offset returns the UTC offset, in seconds, of this city
func (y City) Offset() int {
	if y.parent == nil {
		return -1
	}

	return y.parent.Offset()
}
