package types

import (
	"io"
	"sort"
	"strings"
)

// Strings represents string array
type Strings []string

func NewStrings(s []string) *Strings { r := Strings(s); return &r }

// Slice ...
func (s Strings) Slice() []string {
	return []string(s)
}

// Append new string
func (s Strings) Append(e ...string) Strings {
	s = append(s, e...)
	return s
}

// Contains returns if strings contains s1
func (s Strings) Contains(s1 string) bool {
	return s.Index(s1) > -1
}

// Index returns index of s1
func (s Strings) Index(s1 string) int {
	for i, e := range s {
		if e == s1 {
			return i
		}
	}
	return -1
}

// Copy return copied string array
func (s Strings) Copy() *Strings {
	sl := Strings(make([]string, len(s)))
	copy(sl, s)

	return &sl
}

// Remove remove string element
func (s *Strings) Remove(e string) *Strings {
	i := s.Index(e)
	if i == -1 {
		return s
	}

	return s.RemoveAt(i)
}

// RemoveAt remove indexed string
func (s *Strings) RemoveAt(i int) *Strings {
	*s = Strings(append((*s)[:i], (*s)[i+1:]...))
	return s
}

// Equals return true if equals else return false
func (s Strings) Equals(s1 Strings) bool {
	if len(s) != len(s1) {
		return false
	}

	for i := range s {
		if s[i] != s1[i] {
			return false
		}
	}

	return true
}

// EqualFold return true if equalFold else return false
func (s Strings) EqualFold(s1 Strings) bool {
	if len(s) != len(s1) {
		return false
	}

	for i := range s {
		if strings.EqualFold(s[i], s1[i]) {
			return false
		}
	}

	return true
}

// ToInterface returns []interfaces{}
func (s Strings) ToInterface() (result []interface{}) {
	result = make([]interface{}, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = s[i]
	}

	return
}

// Reader returns new concat reader
func (s Strings) Reader(sep string) io.Reader { return strings.NewReader(strings.Join(s, sep)) }

// Join ...
func (s Strings) Join(sep string) string { return strings.Join(s, sep) }

// Sort string array
func (s Strings) Sort() { sort.Strings(s) }

func (s Strings) Map(mapper func(s string) string) Strings {
	r := make([]string, 0, len(s))
	for _, e := range s {
		r = append(r, mapper(e))
	}
	return r
}
