package goxp

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"gopkg.in/yaml.v3"
)

type decoder interface {
	Decode(v any) error
}

func decode[T any](d decoder) (*T, error) {
	v := new(T)
	return v, d.Decode(v)
}

// ReadJSON read json from response with generics
func ReadJSON[T any](r io.Reader) (*T, error) { return decode[T](json.NewDecoder(r)) }

// ReadXML decode xml with generics
func ReadXML[T any](r io.Reader) (*T, error) { return decode[T](xml.NewDecoder(r)) }

// ReadYAML decode yaml with generics
func ReadYAML[T any](r io.Reader) (*T, error) { return decode[T](yaml.NewDecoder(r)) }