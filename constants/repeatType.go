package constants

import (
	"encoding/json"
	"errors"
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

func ValidateReminderFields(repeatType RepeatType, dayOfWeek, dayOfMonth *int16) error {
	switch repeatType {
	case RepeatDaily:
		if dayOfWeek != nil {
			return errors.New("dayOfWeek should not be set for daily reminders")
		}
		if dayOfMonth != nil {
			return errors.New("dayOfMonth should not be set for daily reminders")
		}

	case RepeatWeekly:
		if dayOfWeek == nil {
			return errors.New("dayOfWeek is required for weekly reminders")
		}
		if dayOfMonth != nil {
			return errors.New("dayOfMonth should not be set for weekly reminders")
		}

	case RepeatMonthly:
		if dayOfMonth == nil {
			return errors.New("dayOfMonth is required for monthly reminders")
		}
		if dayOfWeek != nil {
			return errors.New("dayOfWeek should not be set for monthly reminders")
		}

	default:
		return fmt.Errorf("invalid repeatType: %s", repeatType)
	}

	return nil
}
