package types

import "sort"

// Strings represents string array
type Strings []string

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

// Sort string array
func (s Strings) Sort() {
	sort.Strings(s)
}
