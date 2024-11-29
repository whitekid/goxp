package goxp

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/whitekid/goxp/errors"
)

var wellKnownDateTimeLayouts = [...]string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,

	"2006-01-02 15:04:05.999999999 -0700 MST", // String() format
	// usually used
	"2006-01-02T15:04:05Z0700", // RFC3339Z without colon(:)
	"2006. 1. 2.",
	"January 2, 2006",
}

// ParseDateTime parse standard layout string to time
func ParseDateTime(s string) (time.Time, error) {
	for _, layout := range wellKnownDateTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, err
		}
	}

	return time.Time{}, errors.Errorf(nil, "parse failed: %s", s)
}

// TimeWithLayout time type for struct encoding/decodings
type TimeWithLayout struct {
	time.Time
}

func (t *TimeWithLayout) parseWithLayout(value, layout string) error {
	tm, err := time.Parse(layout, value)
	if err != nil {
		return err
	}
	t.Time = tm
	return nil
}

func (t *TimeWithLayout) marshalJSONWithLayout(layout string) ([]byte, error) {
	return json.Marshal(t.Format(layout))
}

func (t *TimeWithLayout) unmarshalJSONWithLayout(data []byte, layout string) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	tm, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	t.Time = tm
	return nil
}

func (t *TimeWithLayout) marshalXMLWithLayout(e *xml.Encoder, start xml.StartElement, layout string) error {
	return e.EncodeElement(t.Format(layout), start)
}

func (t *TimeWithLayout) unmarshalXMLWithLayout(d *xml.Decoder, start xml.StartElement, layout string) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	tm, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	t.Time = tm
	return nil
}

func (t *TimeWithLayout) marshalYAMLWithLayout(layout string) (any, error) {
	return t.Format(layout), nil
}

func (t *TimeWithLayout) unmarshalYAMLWithLayout(value *yaml.Node, layout string) error {
	tm, err := time.Parse(layout, value.Value)
	if err != nil {
		return err
	}

	t.Time = tm
	return nil
}

type RFC1123ZTime struct {
	TimeWithLayout
}

var (
	_ json.Marshaler   = (*RFC1123ZTime)(nil)
	_ json.Unmarshaler = (*RFC1123ZTime)(nil)
	_ yaml.Marshaler   = (*RFC1123ZTime)(nil)
	_ yaml.Unmarshaler = (*RFC1123ZTime)(nil)
	_ xml.Marshaler    = (*RFC1123ZTime)(nil)
	_ xml.Unmarshaler  = (*RFC1123ZTime)(nil)
)

func NewRFC1123ZTime(t time.Time) *RFC1123ZTime {
	return &RFC1123ZTime{
		TimeWithLayout: TimeWithLayout{
			Time: t,
		},
	}
}

func (t *RFC1123ZTime) Layout() string       { return time.RFC1123Z }
func (t *RFC1123ZTime) String() string       { return t.Format(t.Layout()) }
func (t *RFC1123ZTime) Parse(s string) error { return t.parseWithLayout(s, t.Layout()) }

func (t *RFC1123ZTime) UnmarshalJSON(data []byte) error {
	return t.unmarshalJSONWithLayout(data, t.Layout())
}

func (t *RFC1123ZTime) MarshalJSON() ([]byte, error) {
	return t.marshalJSONWithLayout(t.Layout())
}

func (t *RFC1123ZTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return t.unmarshalXMLWithLayout(d, start, t.Layout())
}

func (t *RFC1123ZTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return t.marshalXMLWithLayout(e, start, t.Layout())
}

func (t *RFC1123ZTime) MarshalYAML() (any, error) {
	return t.marshalYAMLWithLayout(t.Layout())
}

func (t *RFC1123ZTime) UnmarshalYAML(value *yaml.Node) error {
	return t.unmarshalYAMLWithLayout(value, t.Layout())
}

type RFC3339Time struct {
	TimeWithLayout
}

var (
	_ json.Marshaler   = (*RFC3339Time)(nil)
	_ json.Unmarshaler = (*RFC3339Time)(nil)
	_ yaml.Marshaler   = (*RFC3339Time)(nil)
	_ yaml.Unmarshaler = (*RFC3339Time)(nil)
	_ xml.Marshaler    = (*RFC3339Time)(nil)
	_ xml.Unmarshaler  = (*RFC3339Time)(nil)
)

func NewRFC3339Time(t time.Time) *RFC3339Time {
	return &RFC3339Time{
		TimeWithLayout: TimeWithLayout{
			Time: t,
		},
	}
}

func (t *RFC3339Time) Layout() string       { return time.RFC3339 }
func (t *RFC3339Time) String() string       { return t.Format(t.Layout()) }
func (t *RFC3339Time) Parse(s string) error { return t.parseWithLayout(s, t.Layout()) }

func (t *RFC3339Time) UnmarshalJSON(data []byte) error {
	return t.unmarshalJSONWithLayout(data, t.Layout())
}

func (t *RFC3339Time) MarshalJSON() ([]byte, error) {
	return t.marshalJSONWithLayout(t.Layout())
}

func (t *RFC3339Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return t.unmarshalXMLWithLayout(d, start, t.Layout())
}

func (t *RFC3339Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return t.marshalXMLWithLayout(e, start, t.Layout())
}

func (t *RFC3339Time) MarshalYAML() (any, error) {
	return t.marshalYAMLWithLayout(t.Layout())
}

func (t *RFC3339Time) UnmarshalYAML(value *yaml.Node) error {
	return t.unmarshalYAMLWithLayout(value, t.Layout())
}

type UnixTimestamp struct {
	time.Time
}

var (
	_ json.Marshaler   = (*UnixTimestamp)(nil)
	_ json.Unmarshaler = (*UnixTimestamp)(nil)
	_ yaml.Marshaler   = (*UnixTimestamp)(nil)
	_ yaml.Unmarshaler = (*UnixTimestamp)(nil)
	_ xml.Marshaler    = (*UnixTimestamp)(nil)
	_ xml.Unmarshaler  = (*UnixTimestamp)(nil)
)

func NewUnixtimestamp(t int64) *UnixTimestamp {
	return &UnixTimestamp{
		Time: time.Unix(t, 0),
	}
}

func (t *UnixTimestamp) String() string { return strconv.FormatInt(t.Unix(), 10) }

func (t *UnixTimestamp) Parse(s string) error {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(v, 0)
	return nil
}

func (t *UnixTimestamp) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.Time = time.Unix(v, 0)
	return nil
}

func (t *UnixTimestamp) MarshalJSON() ([]byte, error) { return json.Marshal(t.Unix()) }

func (t *UnixTimestamp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int64

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	t.Time = time.Unix(v, 0)

	return nil
}

func (t *UnixTimestamp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(t.Unix(), start)
}

func (t *UnixTimestamp) UnmarshalYAML(value *yaml.Node) error {
	v, err := strconv.ParseInt(value.Value, 10, 64)
	if err != nil {
		return err
	}

	t.Time = time.Unix(v, 0)
	return nil
}

func (t *UnixTimestamp) MarshalYAML() (any, error) { return t.Unix(), nil }
