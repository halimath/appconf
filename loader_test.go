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
	assertLoader(t, JSONFile("./testdata/config.json", true))
}

func TestYAMLFile(t *testing.T) {
	assertLoader(t, YAMLFile("./testdata/config.yaml", true))
}

func TestTOMLFile(t *testing.T) {
	assertLoader(t, TOMLFile("./testdata/config.toml", true))
}

func TestEnv(t *testing.T) {
	t.Setenv("FOO_WEB_ADDRESS", "localhost:8080")
	t.Setenv("FOO_WEB_TIMEOUT", "2s")
	t.Setenv("FOO_WEB_AUTHORIZE", "true")
	t.Setenv("FOO_DB_TYPE", "mysql")
	t.Setenv("FOO_DB_HOST", "localhost")
	t.Setenv("FOO_DB_PORT", "3306")
	t.Setenv("FOO_DB_USER", "test")
	t.Setenv("FOO_DB_PASSWORD", "secret")
	t.Setenv("FOO_BACKENDS_0_HOST", "alpha")
	t.Setenv("FOO_BACKENDS_0_PORT", "8080")
	t.Setenv("FOO_BACKENDS_0_TAGS_0", "a")
	t.Setenv("FOO_BACKENDS_0_TAGS_1", "1")
	t.Setenv("FOO_BACKENDS_1_HOST", "beta")
	t.Setenv("FOO_BACKENDS_1_PORT", "8081")
	t.Setenv("FOO_BACKENDS_1_TAGS_0", "b")
	t.Setenv("FOO_BACKENDS_1_TAGS_1", "2")

	assertLoader(t, Env("FOO"))
}

func assertLoader(t *testing.T, l Loader) {
	got, err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(standardConfig, got); diff != nil {
		got.Dump(0)
		t.Error(diff)
	}
}
