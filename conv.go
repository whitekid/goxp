package goxp

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func AtoiDef[T constraints.Integer](s string, defValue T) T {
	value, err := strconv.Atoi(s)
	if err != nil {
		return defValue
	}

	return T(value)
}

func ParseBoolDef(s string, def bool) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return v
}
