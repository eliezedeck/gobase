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

func TestISO8601ShortDateToTime(t *testing.T) {
	parsed, err := ISO8601ShortDateToTime("2021-07-22")
	if err != nil {
		panic(err)
	}
	t.Log(parsed)
}

func TestCompareISO8601Dates(t *testing.T) {
	yes, err := CompareISO8601Dates("2021-07-21", "2021-07-22")
	if err != nil {
		panic(err)
	}
	if !yes {
		t.Fatal("wrong comparison")
	}
}
