package json

import (
	"encoding/json"
	"io"
)

// DeserializeStruct: creates a struct from io.Reader interface
func DeserializeStruct[T interface{}](reader io.Reader) (*T, error) {
	var C T

	err := json.NewDecoder(reader).Decode(&C)

	return &C, err
}

// SerializeStruct: Creates a []byte from a struct
func SerializeStruct[T interface{}](c T) ([]byte, error) {
	return json.Marshal(c)
}
