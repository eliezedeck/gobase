package validation

import (
	"encoding/json"
	"io"
)

// DecodeJSONStrict decodes into the given `i` while making sure not to allow extra fields that are not present in the
// struct for a strict decoding.
func DecodeJSONStrict(readable io.Reader, i interface{}) error {
	decoder := json.NewDecoder(readable)
	decoder.DisallowUnknownFields()
	return decoder.Decode(i)
}
