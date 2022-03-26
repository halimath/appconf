package appconf

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

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

func File(filename string, mandatory bool, l ReaderLoaderFunc) Loader {
	return LoaderFunc(func() (*Node, error) {
		f, err := os.Open(filename)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) && !mandatory {
				return NewNode(""), nil
			}
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

func JSONFile(name string, mandatory bool) Loader {
	return File(name, mandatory, JSON)
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

func YAMLFile(name string, mandatory bool) Loader {
	return File(name, mandatory, YAML)
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

func TOMLFile(name string, mandatory bool) Loader {
	return File(name, mandatory, TOML)
}

// --

func Env(prefix string) Loader {
	if !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	return LoaderFunc(func() (*Node, error) {
		envMap := make(map[string]interface{})

		for _, envVar := range os.Environ() {
			if !strings.HasPrefix(envVar, prefix) {
				continue
			}
			keyVal := strings.Split(envVar, "=")

			envMap[envKeyToMapKey(keyVal[0], prefix)] = keyVal[1]
		}

		return ConvertToNode(envMap)
	})
}

func envKeyToMapKey(k, prefix string) string {
	k = strings.Replace(k, prefix, "", 1)
	return strings.ReplaceAll(strings.ToLower(k), "_", KeySeparator)
}
