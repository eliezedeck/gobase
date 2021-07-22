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

// ISO8601DateToTime converts a date like 2021-07-22 as a time.Time. Note that the resulting timezone is +00:00.
func ISO8601DateToTime(date string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05-0700", fmt.Sprintf("%sT00:00:00-0000", date))
}

// CompareISO8601Dates is a convenient method to compare that one ISO8601 date is older than the other.
func CompareISO8601Dates(older, recent string) (bool, error) {
	old, err := ISO8601DateToTime(older)
	if err != nil {
		return false, err
	}
	fresh, err := ISO8601DateToTime(recent)
	if err != nil {
		return false, err
	}
	return old.Before(fresh), nil
}
