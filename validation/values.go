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
