package appconf

import (
	"encoding/json"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Loader interface {
	Load() (*Node, error)
}

type LoaderFunc func() (*Node, error)

func (l LoaderFunc) Load() (*Node, error) {
	return l()
}

type ReaderLoaderFunc func(io.Reader) (*Node, error)

func Static(m map[string]interface{}) Loader {
	return LoaderFunc(func() (*Node, error) {
		return ConvertToNode(m)
	})
}

func File(filename string, l ReaderLoaderFunc) Loader {
	return LoaderFunc(func() (*Node, error) {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		return l(f)
	})
}

// --

func JSON(r io.Reader) (*Node, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return ConvertToNode(m)
}

func JSONFile(name string) Loader {
	return File(name, JSON)
}

// --

func YAML(r io.Reader) (*Node, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return ConvertToNode(m)
}

func YAMLFile(name string) Loader {
	return File(name, YAML)
}

// --

func TOML(r io.Reader) (*Node, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := toml.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return ConvertToNode(m)
}

func TOMLFile(name string) Loader {
	return File(name, TOML)
}
