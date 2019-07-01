package yaml

import (
	"github.com/go-yaml/yaml"
	"github.com/nissy/envexpand"
)

func Open(filename string, out interface{}) error {
	return envexpand.Open(filename, out, yaml.Unmarshal)
}
