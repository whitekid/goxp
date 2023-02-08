package timex

// package timex time.Time pointer operations

import "time"

// TODO move to goxp
func NowP() *time.Time {
	v := time.Now()
	return &v
}

func AddP(t time.Time, d time.Duration) *time.Time {
	tm := t.Add(d)
	return &tm
}
