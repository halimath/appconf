package appconf

import (
	"testing"
	"time"
)

func TestAppConfig_Bind_struct(t *testing.T) {
	type (
		DB struct {
			Engine string `appconf:"type"`
			Host   string
			Port   int
			User   string `appconf:",ignore"`
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

	assertEqual(t, want, config)
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

	assertEqual(t, want, config)
}

// func ExampleAppConfig_Bind_map() {
// 	c, err := appconf.New(appconf.JSONFile("./testdata/config.json"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	var config map[string]interface{}
// 	if err := c.Bind(&config); err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("%v\n", config)

// 	// Output:
// }
