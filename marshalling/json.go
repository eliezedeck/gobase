package marshalling

import (
	"bytes"
	"encoding/json"
)

func StructToJSON(i interface{}) ([]byte, error) {
	buff := &bytes.Buffer{}
	encoder := json.NewEncoder(buff)
	if err := encoder.Encode(i); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
