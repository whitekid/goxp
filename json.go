package goxp

import (
	"encoding/json"
	"io"
)

// ReadJSON read json from response with generics
func ReadJSON[T any](r io.Reader) (*T, error) {
	v := new(T)
	return v, json.NewDecoder(r).Decode(v)
}
