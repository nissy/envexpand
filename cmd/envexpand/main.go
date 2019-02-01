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

var (
	version   = "0.0.1"
	format    = flag.String("format", "", "show version and exit")
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
		_, err = fmt.Fprintf(os.Stderr, "Usage: %s [options] file\n", os.Args[0])
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
	if len(*format) > 0 {
		filename = fmt.Sprintf(".%s", *format)
	}

	switch {
	case strings.HasSuffix(filename, ".json"), strings.HasSuffix(filename, ".JSON"):
		if err := json.Unmarshal(in, out); err != nil {
			return nil, err
		}
		if err := envexpand.Do(out); err != nil {
			return nil, err
		}
		return json.Marshal(out)
	case strings.HasSuffix(filename, ".toml"), strings.HasSuffix(filename, ".TOML"):
		if err := toml.Unmarshal(in, out); err != nil {
			return nil, err
		}
		if err := envexpand.Do(out); err != nil {
			return nil, err
		}
		return toml.Marshal(out)
	case strings.HasSuffix(filename, ".yaml"), strings.HasSuffix(filename, ".YAML"), strings.HasSuffix(filename, ".yml"), strings.HasSuffix(filename, ".YML"):
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
