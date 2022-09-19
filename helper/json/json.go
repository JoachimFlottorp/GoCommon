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

func WriteStruct[T interface{}](writer io.Writer, c T) error {
	serialized, err := SerializeStruct(c)
	if err != nil {
		return err
	}

	_, err = writer.Write(serialized)
	return err
}
