package utils

import (
	"time"

	"github.com/pkg/errors"
)

// SetBit Sets the bit at pos in the  n.
func SetBit(n byte, pos int) byte {
	n |= (1 << pos)
	return n
}

// ClearBit Clears the bit at pos in n.
func ClearBit(n byte, pos int) byte { return n &^ (1 << pos) }

// StrToTime parse standard layout string to time
func StrToTime(s string) (time.Time, error) {
	for _, layout := range [...]string{
		time.RFC3339,
		"2006-01-02 15:04:05.999999999 -0700 MST", // String() format
		time.RFC1123Z,
		time.RFC1123,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, err
		}
	}

	return time.Time{}, errors.Errorf("parse failed: %s", s)
}
