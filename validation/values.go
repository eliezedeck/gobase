package validation

import (
	"fmt"
	"time"
)

// IsValidHostnamePort returns true if the given hostport is a valid host:port
func IsValidHostnamePort(hostport string) bool {
	if err := V.Var(hostport, "required,hostname_port"); err != nil {
		return false
	}
	return true
}

// ISO8601ShortDateToTime converts a date like 2021-07-22 as a time.Time. Note that the resulting timezone is +00:00.
func ISO8601ShortDateToTime(date string, location ...*time.Location) (time.Time, error) {
	if len(location) > 0 {
		return time.ParseInLocation("2006-01-02T15:04:05", fmt.Sprintf("%sT00:00:00", date), location[0])
	}
	return time.ParseInLocation("2006-01-02T15:04:05", fmt.Sprintf("%sT00:00:00", date), time.Local)
}

// CompareISO8601Dates is a convenient method to compare that one ISO8601 date is older than the other.
func CompareISO8601Dates(older, recent string) (bool, error) {
	old, err := ISO8601ShortDateToTime(older)
	if err != nil {
		return false, err
	}
	fresh, err := ISO8601ShortDateToTime(recent)
	if err != nil {
		return false, err
	}
	return old.Before(fresh), nil
}
