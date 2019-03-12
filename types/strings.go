package types

import (
	"io"
	"sort"
	"strings"
)

// Strings represents string array
type Strings []string

// Slice ...
func (s *Strings) Slice() []string {
	return []string(*s)
}

// Add new string
func (s *Strings) Add(e ...string) {
	for _, ee := range e {
		*s = append(*s, ee)
	}
}

// Contains returns if strings contains s1
func (s Strings) Contains(s1 string) bool {
	return s.Index(s1) >= 0
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
func (s Strings) Copy() (r Strings) {
	for _, e := range s {
		r = append(r, e)
	}

	return r
}

// Remove remove index
func (s Strings) Remove(e string) Strings {
	i := s.Index(e)
	if i == -1 {
		return s
	}

	return s.RemoveAt(i)
}

// RemoveAt remove index
func (s Strings) RemoveAt(i int) Strings {
	return append(s[:i], s[i+1:]...)
}

// Equals return true if equals else return false
func (s Strings) Equals(s1 Strings) bool {
	if len(s) != len(s1) {
		return false
	}

	for i, e := range s {
		if e != s1[i] {
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

// Reader returns new concat readers
func (s Strings) Reader(sep string) io.Reader {
	n := len(s)
	switch n {
	case 0:
		return strings.NewReader("")
	case 1:
		return strings.NewReader(s[0])
	}

	readers := make([]io.Reader, n*2-1)

	readers[0] = strings.NewReader(s[0])
	for i := 1; i < n; i++ {
		readers[i*2-1] = strings.NewReader(sep)
		readers[i*2] = strings.NewReader(s[i])
	}

	return io.MultiReader(readers...)
}

// Join ...
func (s *Strings) Join(sep string) string {
	return strings.Join([]string(*s), sep)
}

// Sort string array
func (s Strings) Sort() {
	sort.Strings(s)
}
