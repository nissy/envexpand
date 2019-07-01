package json

import (
	"encoding/json"

	"github.com/nissy/envexpand"
)

func Open(filename string, out interface{}) error {
	return envexpand.Open(filename, out, json.Unmarshal)
}
