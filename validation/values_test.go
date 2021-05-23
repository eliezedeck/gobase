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
