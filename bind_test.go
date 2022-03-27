package appconf

import (
	"testing"
	"time"

	"github.com/halimath/assertthat-go/assert"
	"github.com/halimath/assertthat-go/is"
)

func TestAppConfig_Bind_struct(t *testing.T) {
	type (
		DB struct {
			Engine      string `appconf:"type"`
			Host        string
			Port        int
			User        string `appconf:",ignore"`
			KeyNotFound string
		}

		Web struct {
			Address   string
			Timeout   time.Duration
			Authorize bool
		}

		Backend struct {
			Host string
			Port int
			Tags []string
		}

		Config struct {
			DB       DB
			Web      Web
			Backends []Backend
		}
	)

	ac := &AppConfig{n: standardConfig}

	var config Config
	if err := ac.Bind(&config); err != nil {
		panic(err)
	}

	want := Config{
		DB: DB{
			Engine: "mysql",
			Host:   "localhost",
			Port:   3306,
		},
		Web: Web{
			Address:   "localhost:8080",
			Timeout:   2 * time.Second,
			Authorize: true,
		},
		Backends: []Backend{
			{
				Host: "alpha",
				Port: 8080,
				Tags: []string{"a", "1"},
			},
			{
				Host: "beta",
				Port: 8081,
				Tags: []string{"b", "2"},
			},
		},
	}

	assert.That(t, config, is.DeepEqual(want))
}

func TestAppConfig_Bind_map(t *testing.T) {
	ac := &AppConfig{n: standardConfig}

	var config map[string]interface{}
	if err := ac.Bind(&config); err != nil {
		panic(err)
	}

	want := map[string]interface{}{
		"db": map[string]interface{}{
			"type":     "mysql",
			"host":     "localhost",
			"port":     "3306",
			"user":     "test",
			"password": "secret",
		},
		"web": map[string]interface{}{
			"address":   "localhost:8080",
			"timeout":   "2s",
			"authorize": "true",
		},
		"backends": map[string]interface{}{
			"0": map[string]interface{}{
				"host": "alpha",
				"port": "8080",
				"tags": map[string]interface{}{"0": "a", "1": "1"},
			},
			"1": map[string]interface{}{
				"host": "beta",
				"port": "8081",
				"tags": map[string]interface{}{"0": "b", "1": "2"},
			},
		},
	}

	assert.That(t, config, is.DeepEqual(want))
}

func TestAppConfig_Bind_scalars(t *testing.T) {
	type C struct {
		Int   int   `appconf:"v"`
		Int8  int8  `appconf:"v"`
		Int16 int16 `appconf:"v"`
		Int32 int32 `appconf:"v"`
		Int64 int64 `appconf:"v"`

		Uint   uint   `appconf:"v"`
		Uint8  uint8  `appconf:"v"`
		Uint16 uint16 `appconf:"v"`
		Uint32 uint32 `appconf:"v"`
		Uint64 uint64 `appconf:"v"`

		Float32 float32 `appconf:"v"`
		Float64 float64 `appconf:"v"`

		Complex64  complex64  `appconf:"v"`
		Complex128 complex128 `appconf:"v"`
	}

	n, err := ConvertToNode(map[string]interface{}{
		"v": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	ac := &AppConfig{
		n: n,
	}

	var c C
	if err := ac.Bind(&c); err != nil {
		t.Fatal(err)
	}

	assert.That(t, c,
		is.Equal(C{
			Int:        1,
			Int8:       1,
			Int16:      1,
			Int32:      1,
			Int64:      1,
			Uint:       1,
			Uint8:      1,
			Uint16:     1,
			Uint32:     1,
			Uint64:     1,
			Float32:    1,
			Float64:    1,
			Complex64:  1,
			Complex128: 1,
		}),
	)
}
