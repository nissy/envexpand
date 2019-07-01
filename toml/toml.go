package toml

import (
	"github.com/BurntSushi/toml"
	"github.com/nissy/envexpand"
)

func Open(filename string, out interface{}) error {
	return envexpand.Open(filename, out, toml.Unmarshal)
}
