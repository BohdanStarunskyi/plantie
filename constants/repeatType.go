package constants

import (
	"encoding/json"
	"fmt"
)

type RepeatType int64

const (
	RepeatDaily RepeatType = iota
	RepeatWeekly
	RepeatMonthly
)

func (r RepeatType) String() string {
	return [...]string{"daily", "weekly", "monthly"}[r]
}

func (r *RepeatType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		switch s {
		case "daily":
			*r = RepeatDaily
		case "weekly":
			*r = RepeatWeekly
		case "monthly":
			*r = RepeatMonthly
		default:
			return fmt.Errorf("invalid repeatType: %s", s)
		}
		return nil
	}
	var i int64
	if err := json.Unmarshal(b, &i); err == nil {
		*r = RepeatType(i)
		return nil
	}
	return fmt.Errorf("invalid repeatType format")
}

func (r RepeatType) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}
