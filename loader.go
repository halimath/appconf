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

// Loader defines the interface implemented by all loaders.
type Loader interface {
	// Load loads the configuration values and returns them as a tree of Nodes. If the loading was not
	// successful an error should be returned.
	Load() (*Node, error)
}

// LoaderFunc is a convenience type to convert a function to a Loader.
type LoaderFunc func() (*Node, error)

func (l LoaderFunc) Load() (*Node, error) {
	return l()
}

// ReaderLoaderFunc is function type to implement Loaders that consume an io.Reader.
type ReaderLoaderFunc func(io.Reader) (*Node, error)

// Static creates a Loader that returns static configuration values from the given map structure. The map's
// values are limited to strings, map[string]interface{} (with the same value constraints applied) or slices
// of either strings or maps.
func Static(m map[string]interface{}) Loader {
	return LoaderFunc(func() (*Node, error) {
		return ConvertToNode(m)
	})
}

// File creates a Loader that reads the file named filename and forwards the content to l. If mandatory is
// set to false, an empty configuration will be returned when filename does not exist. Otherwise this is is
// reported as an error.
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

// JSON parses the r's content as JSON and converts it to a Node tree.
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

// JSONFile creates a Loader which loads JSON configuration from a file name.
func JSONFile(name string, mandatory bool) Loader {
	return File(name, mandatory, JSON)
}

// --

// YAML loades the content from r and converts it to a Node tree.
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

// YAMLFile creates a Loader which loads YAML configuration from a file name.
func YAMLFile(name string, mandatory bool) Loader {
	return File(name, mandatory, YAML)
}

// --

// TOML loads the content from r and converts it to a Node tree.
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

// TOMLFile creates a Loader which loads TOML configuration from a file name.
func TOMLFile(name string, mandatory bool) Loader {
	return File(name, mandatory, TOML)
}

// --

// Env creates a Loader which reads configuration values from the environment. Only env variables with a
// name starting with prefix are considered. Use the empty string to select all variables.
func Env(prefix string) Loader {
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "_") {
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
