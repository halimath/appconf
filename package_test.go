package appconf_test

import (
	"fmt"
	"time"

	"github.com/halimath/appconf"
)

func ExampleAppConfig_Bind() {
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

		Config struct {
			DB  DB
			Web Web
		}
	)

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
	// {{mysql localhost 3306 } {localhost:8080 2s true}}
}

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
