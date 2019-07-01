package toml

import (
	"github.com/nissy/envexpand"
	"github.com/pelletier/go-toml"
)

func Open(filename string, out interface{}) error {
	return envexpand.Open(filename, out, toml.Unmarshal)
}
