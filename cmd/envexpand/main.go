package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"github.com/nissy/envexpand"
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
}

func run() error {
	flag.Parse()
	args := flag.Args()

	if *isVersion {
		fmt.Println("v" + version)
		return nil
	}

	if *isHelp || len(args) == 0 {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s [options] [file JSON|YAML|TOML]\n", os.Args[0])
		flag.PrintDefaults()
		return err
	}

	var out interface{}
	data, err := expand(args[0], &out)
	if err != nil {
		return err
	}

	fmt.Print(string(data))

	return nil
}

func expand(filename string, out interface{}) ([]byte, error) {
	switch {
	case isFileExtension(filename, "JSON"):
		if err := envexpand.Open(filename, out, json.Unmarshal); err != nil {
			return nil, err
		}
		return json.Marshal(out)
	case isFileExtension(filename, "TOML"):
		if err := envexpand.Open(filename, out, toml.Unmarshal); err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		err := toml.NewEncoder(&buf).Encode(out)
		return buf.Bytes(), err
	case isFileExtension(filename, "YAML", "YML"):
		if err := envexpand.Open(filename, out, yaml.Unmarshal); err != nil {
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
