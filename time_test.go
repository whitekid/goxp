package goxp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestParseDateTime(t *testing.T) {
	type args struct {
		str  string
		want time.Time
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ANSIC", args{time.Layout, time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -25200))}, false},
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
		{"", args{"2021. 3. 2.", time.Date(2021, 3, 2, 0, 0, 0, 0, time.UTC)}, false},
		{"", args{"January 2, 2006", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateTime(tt.args.str)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `StrToTime() failed`, `error = %+v, wantErr = %v`, err, tt.wantErr)
			}

			require.Equal(t, tt.args.want, got, "want=%s, got=%s", tt.args.want, got)
		})
	}
}

type layter interface {
	Layout() string
}

func TestTimeWithLayout(t *testing.T) {
	now := time.Now()
	type args struct {
		v any
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{`RFC1123ZTime`, args{&RFC1123ZTime{}}, false},
		{`RFC3339Time`, args{&RFC3339Time{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name+"-json", func(t *testing.T) {
			want := fmt.Sprintf(`"%s"`, now.Format(tt.args.v.(layter).Layout()))
			err := json.Unmarshal([]byte(want), tt.args.v)
			require.NoError(t, err)

			got, err := json.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `json.Marshal() failed`, `error = %+v, wantErr = %v`, err, tt.wantErr)
			}
			require.Equal(t, want, string(got))
		})

		t.Run(tt.name+"-yaml", func(t *testing.T) {
			want := fmt.Sprintf("%s\n", now.Format(tt.args.v.(layter).Layout()))
			err := yaml.Unmarshal([]byte(want), tt.args.v)
			require.NoError(t, err)

			got, err := yaml.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `yaml.Marshal() failed`, `error = %+v, wantErr = %v`, err, tt.wantErr)
			}

			if bytes.HasPrefix(got, []byte(`"`)) {
				want = fmt.Sprintf("\"%s\"\n", now.Format(tt.args.v.(layter).Layout()))
			}
			require.Equal(t, want, string(got))
		})

		t.Run(tt.name+"-xml", func(t *testing.T) {
			typName := reflect.TypeOf(tt.args.v).Elem().Name()
			want := fmt.Sprintf(`<%s>%s</%s>`, typName, now.Format(tt.args.v.(layter).Layout()), typName)
			err := xml.Unmarshal([]byte(want), tt.args.v)
			require.NoError(t, err)

			got, err := xml.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `xml.Marshal() failed`, `error = %+v, wantErr = %v`, err, tt.wantErr)
			}
			require.Equal(t, want, string(got))
		})
	}
}

func TestUnixTimestamp(t *testing.T) {
	now := time.Now()
	type args struct {
		v any
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{`UnixTimestamp`, args{&UnixTimestamp{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name+"-json", func(t *testing.T) {
			want := fmt.Sprintf("%d", now.Unix())
			err := json.Unmarshal([]byte(want), tt.args.v)
			require.NoError(t, err)

			got, err := json.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `json.Marshal() failed`, `error = %+v, wantErr = %v`, err, tt.wantErr)
			}
			require.Equal(t, want, string(got))
		})

		t.Run(tt.name+"-yaml", func(t *testing.T) {
			want := fmt.Sprintf("%d\n", now.Unix())
			err := yaml.Unmarshal([]byte(want), tt.args.v)
			require.NoError(t, err)

			got, err := yaml.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `yaml.Marshal() failed`, `error = %+v, wantErr = %v`, err, tt.wantErr)
			}

			require.Equal(t, want, string(got))
		})

		t.Run(tt.name+"-xml", func(t *testing.T) {
			typName := reflect.TypeOf(tt.args.v).Elem().Name()
			want := fmt.Sprintf(`<%s>%d</%s>`, typName, now.Unix(), typName)
			err := xml.Unmarshal([]byte(want), tt.args.v)
			require.NoError(t, err)

			got, err := xml.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `xml.Marshal() failed`, `error = %+v, wantErr = %v`, err, tt.wantErr)
			}
			require.Equal(t, want, string(got))
		})
	}
}
