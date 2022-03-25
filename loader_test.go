package appconf

import (
	"testing"

	"github.com/go-test/deep"
)

func TestStatic(t *testing.T) {
	assertLoader(t, Static(map[string]interface{}{
		"web": map[string]interface{}{
			"address":   "localhost:8080",
			"timeout":   "2s",
			"authorize": true,
		},
		"db": map[string]interface{}{
			"type":     "mysql",
			"host":     "localhost",
			"port":     3306,
			"user":     "test",
			"password": "secret",
		},
		"backends": []interface{}{
			map[string]interface{}{
				"host": "alpha",
				"port": 8080,
				"tags": []interface{}{"a", "1"},
			},
			map[string]interface{}{
				"host": "beta",
				"port": 8081,
				"tags": []interface{}{"b", "2"},
			},
		},
	}))
}

func TestJSONFile(t *testing.T) {
	assertLoader(t, JSONFile("./testdata/config.json"))
}

func TestYAMLFile(t *testing.T) {
	assertLoader(t, YAMLFile("./testdata/config.yaml"))
}

func TestTOMLFile(t *testing.T) {
	assertLoader(t, TOMLFile("./testdata/config.toml"))
}

func assertLoader(t *testing.T, l Loader) {
	got, err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(standardConfig, got); diff != nil {
		t.Error(diff)
	}
}
