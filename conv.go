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

func ParseIntDef[T constraints.Integer](s string, defaultValue, minValue, maxValue T) T {
	value, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return defaultValue
	}

	return min(max(T(value), minValue), maxValue)
}
