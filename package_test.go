package appconf_test

import (
	"fmt"
	"time"

	"github.com/halimath/appconf"
)

func Example() {
	c, err := appconf.New(
		appconf.Static(map[string]interface{}{
			"web.address": "http://example.com",
			"DB": map[string]interface{}{
				"Timeout": 2 * time.Second,
			},
		}),
		appconf.JSONFile("./testdata/config.json"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(c.HasKey("Web"))
	fmt.Println(c.HasKey("Foo"))
	fmt.Println(c.GetString("Web.Address"))
	fmt.Println(c.GetDuration("db.timeout"))

	// Output:
	// true
	// false
	// localhost:8080
	// 2s
}

func ExampleAppConfig_Bind_struct() {
	type Config struct {
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
		Backends []struct {
			Host string
			Port int
			Tags []string
		}
	}

	c, err := appconf.New(appconf.JSONFile("./testdata/config.json"))
	if err != nil {
		panic(err)
	}

	var config Config
	if err := c.Bind(&config); err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", config)

	// Output:
	// {{mysql localhost 3306 } {localhost:8080 2s true} [{alpha 8080 [a 1]} {beta 8081 [b 2]}]}
}

func ExampleAppConfig_Bind_map() {
	c, err := appconf.New(appconf.JSONFile("./testdata/config.json"))
	if err != nil {
		panic(err)
	}

	var config map[string]interface{}
	if err := c.Bind(&config); err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", config)

	// Output:
	// map[backends:map[0:map[host:alpha port:8080 tags:map[0:a 1:1]] 1:map[host:beta port:8081 tags:map[0:b 1:2]]] db:map[host:localhost password:secret port:3306 type:mysql user:test] web:map[address:localhost:8080 authorize:true timeout:2s]]
}
