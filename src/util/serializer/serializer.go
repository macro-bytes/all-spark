package serializer

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// DeserializePath - reads file specified in the path argument
// and deserializes the JSON object to the destination object
func DeserializePath(path string, destination interface{}) error {
	template, err := os.Open(path)
	if err != nil {
		return err
	}
	defer template.Close()

	byteValue, err := ioutil.ReadAll(template)
	if err != nil {
		return err
	}

	return Deserialize([]byte(byteValue), &destination)
}

// Serialize the object parameter to a JSON byte array
func Serialize(object interface{}) ([]byte, error) {
	serialized, err := json.Marshal(object)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}

// Deserialize the byte array to the destination interface
func Deserialize(buffer []byte, destination interface{}) error {
	return json.Unmarshal(buffer, &destination)
}
