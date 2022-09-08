package goxp

import (
	"time"

	"github.com/pkg/errors"
)

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

		// usually used
		"2006-01-02T15:04:05Z0700", // RFC3339Z without colon(:)
		"2006. 1. 2.",
		"January 2, 2006",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, err
		}
	}

	return time.Time{}, errors.Errorf("parse failed: %s", s)
}
