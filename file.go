package envexpand

import (
	"io/ioutil"
	"os"
)

func Open(filename string, out interface{}, fn func(in []byte, out interface{}) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	in, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if err := fn(in, out); err != nil {
		return err
	}

	return Do(out)
}
