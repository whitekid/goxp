package types

import (
	"sort"
	"strings"
)

// Strings represents string array
type Strings []string

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

// Join ...
func (s *Strings) Join(sep string) string {
	return strings.Join([]string(*s), sep)
}

// Sort string array
func (s Strings) Sort() {
	sort.Strings(s)
}
