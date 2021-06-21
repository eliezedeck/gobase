package validation

import (
	"context"
	"encoding/json"
	"io"

	"github.com/eliezedeck/gobase/logging"
	"github.com/eliezedeck/gobase/marshalling"
	"go.uber.org/zap"
)

// DecodeJSONStrict decodes into the given `i` while making sure not to allow extra fields that are not present in the
// struct for a strict decoding.
func DecodeJSONStrict(readable io.Reader, i interface{}) error {
	decoder := json.NewDecoder(readable)
	decoder.DisallowUnknownFields()
	return decoder.Decode(i)
}

// ValidateJSONBody parses the given JSON from the request, unmarshalls it into the given `dest`. If drop != false, the
// JSON is recreated based on the struct `dest` and returned as `validated`.
func ValidateJSONBody(jbody io.Reader, dest interface{}, drop ...bool) (validated []byte, err error) {
	// Decode the JSON, not allowing extra fields
	if err := DecodeJSONStrict(jbody, dest); err != nil {
		logging.L.Error("JSON decoding", zap.Error(err))
		return nil, err // HTTP 500
	}

	// Mold the data (dependent on `mod` tags)
	if err := M.Struct(context.TODO(), dest); err != nil {
		logging.L.Error("JSON data molding", zap.Error(err))
		return nil, err
	}

	// Validate based on the dest struct
	if err := V.Struct(dest); err != nil {
		logging.L.Error("JSON validation", zap.Error(err))
		return nil, err
	}

	if len(drop) == 0 || drop[0] {
		// user doesn't want the repacked JSON to be returned
		return nil, nil
	}

	// Convert back to a safe JSON
	asjson, err := marshalling.StructToJSON(dest)
	if err != nil {
		logging.L.Error("JSON rebuild", zap.Error(err))
		return nil, err // HTTP 500
	}

	return asjson, nil
}
