package goxp

import (
	"encoding/xml"
	"io"
)

// ReadXML decode xml with generics
func ReadXML[T any](r io.Reader) (*T, error) {
	v := new(T)
	return v, xml.NewDecoder(r).Decode(v)
}
