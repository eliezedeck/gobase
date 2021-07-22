package validation

import "testing"

func TestISO8601Validation(t *testing.T) {
	date := "2021-07-10"
	if err := V.Var(date, "date-iso8601"); err != nil {
		t.Error(err)
		return
	}
}
