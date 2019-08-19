package template_reader

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Deserialize - deserializes cloud creation templates
func Deserialize(path string, destination interface{}) error {
	template, err := os.Open(path)
	if err != nil {
		return err
	}
	defer template.Close()

	byteValue, err := ioutil.ReadAll(template)
	if err != nil {
		return err
	}

	if json.Unmarshal([]byte(byteValue), &destination) != nil {
		return err
	}
	return nil
}
