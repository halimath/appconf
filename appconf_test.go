package appconf

import (
	"testing"
	"time"
)

var standardConfig = &Node{
	Children: map[Key]*Node{
		Key("web"): {
			Children: map[Key]*Node{
				Key("address"):   {Value: "localhost:8080"},
				Key("timeout"):   {Value: "2s"},
				Key("authorize"): {Value: "true"},
			},
		},
		Key("db"): {
			Children: map[Key]*Node{
				Key("type"):     {Value: "mysql"},
				Key("host"):     {Value: "localhost"},
				Key("port"):     {Value: "3306"},
				Key("user"):     {Value: "test"},
				Key("password"): {Value: "secret"},
			},
		},

		Key("backends"): {
			Children: map[Key]*Node{
				"0": {
					Children: map[Key]*Node{
						Key("host"): {Value: "alpha"},
						Key("port"): {Value: "8080"},
						Key("tags"): {
							Children: map[Key]*Node{
								Key("0"): {Value: "a"},
								Key("1"): {Value: "1"},
							},
						},
					},
				},
				"1": {
					Children: map[Key]*Node{
						Key("host"): {Value: "beta"},
						Key("port"): {Value: "8081"},
						Key("tags"): {
							Children: map[Key]*Node{
								Key("0"): {Value: "b"},
								Key("1"): {Value: "2"},
							},
						},
					},
				},
			},
		},
	},
}

type (
	db struct {
		Engine string `appconf:"type"`
		Host   string
		Port   int
		User   string `appconf:",ignore"`
	}

	web struct {
		Address   string
		Timeout   time.Duration
		Authorize bool
	}

	config struct {
		DB  db
		Web web
	}
)

func TestAppConfig_Bind_nonNestedStruct(t *testing.T) {
	c := &AppConfig{
		n: standardConfig.Children[Key("db")],
	}

	var d db
	if err := c.Bind(&d); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, db{
		Engine: "mysql",
		Host:   "localhost",
		Port:   3306,
	}, d)
}

func TestAppConfig_Bind_nestedStruct(t *testing.T) {
	c := &AppConfig{
		n: standardConfig,
	}

	var cfg config
	if err := c.Bind(&cfg); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, config{
		DB: db{
			Engine: "mysql",
			Host:   "localhost",
			Port:   3306,
		},
		Web: web{
			Address:   "localhost:8080",
			Timeout:   2 * time.Second,
			Authorize: true,
		},
	}, cfg)
}

func TestAppConfig_Sub(t *testing.T) {
	c := &AppConfig{n: standardConfig}
	s := c.Sub("web")

	assertEqual(t, "localhost:8080", s.GetString("address"))
}

func TestAppConfig_Get(t *testing.T) {
	n, err := ConvertToNode(map[string]interface{}{
		"string":   "foo",
		"int":      67,
		"float":    17.2,
		"bool":     true,
		"complex":  complex(1, 2),
		"duration": time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	c := &AppConfig{
		n: n,
	}

	assertEqual(t, "foo", c.GetString("string"))
	assertEqual(t, "", c.GetString("stringnotfound"))

	assertEqual(t, 67, c.GetInt("int"))
	assertEqual(t, 0, c.GetInt("intnotfound"))
	assertEqual(t, 67, c.GetInt64("int"))
	assertEqual(t, 0, c.GetInt64("intnotfound"))

	assertEqual(t, 67, c.GetUint("int"))
	assertEqual(t, 0, c.GetUint("intnotfound"))
	assertEqual(t, 67, c.GetUint64("int"))
	assertEqual(t, 0, c.GetUint64("intnotfound"))

	assertEqual(t, 17.2, c.GetFloat32("float"))
	assertEqual(t, 0, c.GetFloat32("floatnotfound"))
	assertEqual(t, 17.2, c.GetFloat64("float"))
	assertEqual(t, 0, c.GetFloat64("floatnotfound"))

	assertEqual(t, true, c.GetBool("bool"))
	assertEqual(t, false, c.GetBool("boolnotfound"))

	assertEqual(t, complex(1, 2), c.GetComplex128("complex"))
	assertEqual(t, complex(0, 0), c.GetComplex128("boolnotfound"))

	assertEqual(t, time.Second, c.GetDuration("duration"))
	assertEqual(t, 0, c.GetDuration("durationnotfound"))
}
