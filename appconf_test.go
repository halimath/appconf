package appconf

import (
	"testing"
	"time"

	"github.com/halimath/assertthat-go/assert"
	"github.com/halimath/assertthat-go/is"
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

	assert.That(t, d, is.DeepEqual(db{
		Engine: "mysql",
		Host:   "localhost",
		Port:   3306,
	}))
}

func TestAppConfig_Bind_nestedStruct(t *testing.T) {
	c := &AppConfig{
		n: standardConfig,
	}

	var cfg config
	if err := c.Bind(&cfg); err != nil {
		t.Fatal(err)
	}

	assert.That(t, cfg, is.DeepEqual(config{
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
	}))
}

func TestAppConfig_Sub(t *testing.T) {
	c := &AppConfig{n: standardConfig}
	s := c.Sub("web")

	assert.That(t, s.GetString("address"), is.Equal("localhost:8080"))
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

	assert.That(t, c.GetString("string"), is.Equal("foo"))
	assert.That(t, c.GetString("stringnotfound"), is.Equal(""))

	assert.That(t, c.GetInt("int"), is.Equal(67))
	assert.That(t, c.GetInt("intnotfound"), is.Equal(0))
	assert.That(t, c.GetInt64("int"), is.Equal[int64](67))
	assert.That(t, c.GetInt64("intnotfound"), is.Equal[int64](0))

	assert.That(t, c.GetUint("int"), is.Equal[uint](67))
	assert.That(t, c.GetUint("intnotfound"), is.Equal[uint](0))
	assert.That(t, c.GetUint64("int"), is.Equal[uint64](67))
	assert.That(t, c.GetUint64("intnotfound"), is.Equal[uint64](0))

	assert.That(t, c.GetFloat32("float"), is.Equal[float32](17.2))
	assert.That(t, c.GetFloat32("floatnotfound"), is.Equal[float32](0))
	assert.That(t, c.GetFloat64("float"), is.Equal(17.2))
	assert.That(t, c.GetFloat64("floatnotfound"), is.Equal(0.0))

	assert.That(t, c.GetBool("bool"), is.Equal(true))
	assert.That(t, c.GetBool("boolnotfound"), is.Equal(false))

	assert.That(t, c.GetComplex128("complex"), is.Equal(complex(1, 2)))
	assert.That(t, c.GetComplex128("boolnotfound"), is.Equal(complex(0, 0)))

	assert.That(t, c.GetDuration("duration"), is.Equal(time.Second))
	assert.That(t, c.GetDuration("durationnotfound"), is.Equal[time.Duration](0))
}
