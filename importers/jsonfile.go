package importers

import (
	"encoding/json"
	"io/ioutil"

	tzlib "github.com/elseym/go-tzlib"
)

// JSONFile imports from a json file from the filesystem
type JSONFile struct {
	filename string
}

// NewJSONFile returns a new JSONFile importer reading json from filename
func NewJSONFile(filename string) JSONFile {
	return JSONFile{filename}
}

// Import creates a new Tzlib from unmarshalled json
func (jf JSONFile) Import() (lib *tzlib.Tzlib, err error) {
	var l tzlib.Tzlib

	// rebuild and recalculate data
	if l, err = jf.loadFromDisk(); err == nil {
		lib = tzlib.NewTzlib(l.Created, l.Expires)
		for _, z := range l.Timezones {
			zon := lib.AddNewTimezone(z.Utc)
			for _, c := range z.Countries {
				cou := zon.AddNewCountry(c.Name, c.Weight)
				for _, y := range c.Cities {
					cou.AddNewCity(y.Name, y.Weight)
				}
			}
		}
	}

	return
}

// loadFromDisk
func (jf JSONFile) loadFromDisk() (l tzlib.Tzlib, err error) {
	var b []byte

	b, err = ioutil.ReadFile(jf.filename)
	if err == nil {
		err = json.Unmarshal(b, &l)
	}

	return
}
