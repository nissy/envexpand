package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nissy/envexpand"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

var version string

var (
	isVersion = flag.Bool("v", false, "show version and exit")
	isHelp    = flag.Bool("h", false, "this help")
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func run() (err error) {
	flag.Parse()
	args := flag.Args()

	if *isVersion {
		fmt.Println("v" + version)
		return nil
	}

	if *isHelp {
		_, err = fmt.Fprintf(os.Stderr, "Usage: %s [options] [file JSON|YAML|TOML]\n", os.Args[0])
		flag.PrintDefaults()
		return err
	}

	if len(args) == 0 {
		return nil
	}

	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	in, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var out interface{}
	data, err := expand(filename, in, &out)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func expand(filename string, in []byte, out interface{}) ([]byte, error) {
	switch {
	case isFileExtension(filename, "JSON"):
		if err := json.Unmarshal(in, out); err != nil {
			return nil, err
		}
		if err := envexpand.Do(out); err != nil {
			return nil, err
		}
		return json.Marshal(out)
	case isFileExtension(filename, "TOML"):
		if err := toml.Unmarshal(in, out); err != nil {
			return nil, err
		}
		if err := envexpand.Do(out); err != nil {
			return nil, err
		}
		return toml.Marshal(out)
	case isFileExtension(filename, "YAML", "YML"):
		if err := yaml.Unmarshal(in, out); err != nil {
			return nil, err
		}
		if err := envexpand.Do(out); err != nil {
			return nil, err
		}
		return yaml.Marshal(out)
	}

	return nil, errors.New("Format is not supported.")
}

func isFileExtension(s string, ex ...string) bool {
	for _, v := range ex {
		if vv := fmt.Sprintf(".%s", v); strings.HasSuffix(s, vv) || strings.HasSuffix(s, strings.ToLower(vv)) {
			return true
		}
	}

	return false
}
