package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStrToTime(t *testing.T) {
	type args struct {
		str  string
		want time.Time
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ANSIC", args{"Mon Jan 2 15:04:05 2006", time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)}, false},
		{"UnixDate", args{"Mon Jan 2 15:04:05 MST 2006", time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("MST", 0))}, false},
		{"RubyDate", args{"Mon Jan 02 15:04:05 -0700 2006", time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -25200))}, false},
		{"RFC822", args{"02 Jan 06 15:04 MST", time.Date(2006, 1, 2, 15, 4, 0, 0, time.FixedZone("MST", 0))}, false},
		{"RFC822Z", args{"02 Jan 06 15:04 -0700", time.Date(2006, 1, 2, 15, 4, 0, 0, time.FixedZone("", -25200))}, false},
		{"RFC850", args{"Monday, 02-Jan-06 15:04:05 MST", time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("MST", 0))}, false},
		{"RFC1123", args{"Mon, 02 Jan 2006 15:04:05 MST", time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("MST", 0))}, false},
		{"RFC1123Z", args{"Mon, 02 Jan 2006 15:04:05 -0700", time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -25200))}, false},
		{"RFC3339", args{"2006-01-02T15:04:05+07:00", time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 25200))}, false},
		{"RFC3339Nano", args{"2006-01-02T15:04:05.999999999+07:00", time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.FixedZone("", 25200))}, false},
		{"Kitchen", args{"3:04PM", time.Date(0, 1, 1, 15, 4, 0, 0, time.UTC)}, false},
		{"Stamp", args{"Jan 2 15:04:05", time.Date(0, 1, 2, 15, 4, 5, 0, time.UTC)}, false},
		{"StampMilli", args{"Jan 2 15:04:05.000", time.Date(0, 1, 2, 15, 4, 5, 0, time.UTC)}, false},
		{"StampMicro", args{"Jan 2 15:04:05.000000", time.Date(0, 1, 2, 15, 4, 5, 0, time.UTC)}, false},
		{"StampNano", args{"Jan 2 15:04:05.000000000", time.Date(0, 1, 2, 15, 4, 5, 0, time.UTC)}, false},
		{"String", args{"2006-01-02 15:04:05 +0000 UTC", time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToTime(tt.args.str)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `StrToTime() failed`, `error = %+v, wantErr = %v`, err, tt.wantErr)
			}

			require.Equal(t, tt.args.want, got, "want=%s, got=%s", tt.args.want, got)
		})
	}
}
