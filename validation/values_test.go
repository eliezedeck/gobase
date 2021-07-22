package validation

import (
	"testing"
)

func TestIsValidHostnamePort(t *testing.T) {
	if IsValidHostnamePort("http://localhost:1900") {
		t.Failed()
	}
	if !IsValidHostnamePort("localhost:80") {
		t.Failed()
	}
}

func TestISO8601DateToTime(t *testing.T) {
	parsed, err := ISO8601DateToTime("2021-07-22")
	if err != nil {
		panic(err)
	}
	t.Log(parsed)
}
