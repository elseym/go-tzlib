package exporters

import (
	"encoding/json"
	"io/ioutil"

	tzlib "github.com/elseym/go-tzlib"
)

// JSONFile exports to a json file to the filesystem
type JSONFile struct {
	filename string
}

// NewJSONFile returns a new exporter which will write to filename to the filesystem
func NewJSONFile(filename string) JSONFile {
	return JSONFile{filename}
}

// Export writes the marshalled Tzlib to disk
func (e JSONFile) Export(l *tzlib.Tzlib) error {
	j, err := json.Marshal(&l)
	if err == nil {
		err = ioutil.WriteFile(e.filename, j, 0644)
	}

	return err
}
