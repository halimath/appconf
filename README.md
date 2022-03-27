# appconf

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]

Opinonated configuration loader for golang applications.

# Usage

## Installation

This module uses golang modules and can be installed with

```shell
go get github.com/halimath/appconf@main
```

## Basic Usage

`appconf` provides a type `AppConfig` which collects configuration input from different _loaders_, merges them
together and provides a type-safe API to query fields from the config. Consider this example:


```go
c, err := appconf.New(
	appconf.Static(map[string]interface{}{
		"web.address": "http://example.com",
		"DB": map[string]interface{}{
			"Timeout": 2 * time.Second,
		},
	}),
	appconf.JSONFile("./testdata/config.json", true),
)
if err != nil {
	panic(err)
}

fmt.Println(c.HasKey("Web"))
fmt.Println(c.HasKey("Foo"))
fmt.Println(c.GetString("Web.Address"))
fmt.Println(c.GetDuration("db.timeout"))
```

### Loaders

The example uses two loaders to provide input for the configuration: a `Static` loader providing default
values and a `JSONFile` loader overriding some (or all of the values). 

The module contains loaders for the following kind of input:
* JSON (both from a `Reader` and from a file)
* YAML (both from a `Reader` and from a file)
* TOML (both from a `Reader` and from a file)
* environment variables

You can create your own loader by implementing the `Loader` interface. See below for details.

Once the configuration is loaded the individual values can be queried using their _key_.

### Keys

Keys are used to access configuration values. Keys are considered case-insensitive due to the fact that not 
all loaders are able to deliver case-sensitive keys (such as environment variables). Keys can be nested.
When queriying nested keys use a single dot to separate the parts (this is called a _key path_).

### Getters

When queriying values you can use different getters to convert the value to a desired type. The following
getters are available:

* `GetString`
* `GetInt`
* `GetInt64`
* `GetUint`
* `GetUint64`
* `GetFloat32`
* `GetFloat64`
* `GetComplex128`
* `GetBool`
* `GetDuration`

Each of these getters always returns a valid value. If the key is not defined or if the underlying value can
not be converted to the given type, they return the type's default value. There is a corresponding
`Get...E` version of the getter, which returns an `error` in addition to the value.

### Using sub-configurations

You can call the configuration's `Sub` method to query a key and return the configuration structure rooted at
that key. This is usefull if you want to pass only parts of the configuration to downstream code, such as the
database configuration rooted at key `db`: Simply call `conf.Sub("db")`.

As with the getters described above, `Sub` always returns a valid configuration. It is empty when the given
key is not found. There is a corresponding `SubE` method, which returns a configuration and an optional
`error`.

### Binding

`appconf` supports binding configuration to `struct` values. This is done using reflection and it works
similar to unmarshaling of i.e. JSON code:

```go
c, _ := appconf.Load(...)

var config ConfigStruct

if err := c.Bind(&config); err != nil {
	panic(err)
}
```

The above code shows how to bind to a `ConfigStruct` value. By default each struct field is assigned the
value of the config value with a key formed by converting the field name to lower case. If you want to bind
a different key, add a field tag of the form 

```go
type Config struct {
	SomeValue string `appconf:"somekey"`
}
```

If you want to ignore a struct field during binding add the field tag `appconf:",ignore"`. Note the comma 
before `ignore` which is important as otherwise the field would be bound to a key named `ignore`.

Bindings works with nested structs and nested slices. The keys for slice elements are formed by putting the
index as a single key path element, i.e. `db.hosts.0.name`.

You can also bind the configuration to a `map[string]interface{}`. Keep in mind, that all leaf values are
added as `string`s.

### Implementing a custom loader

To implement a custom configuration loader you create a type which implements the `Loader` interface. This 
interface contains a single method which loads the configuration and returns it as a `Node` in addition to any
`error`. `Node`s form a tree with keys annotated to each `Node`. Leaf `Node`s carry the configuration values.
This is the internal representation this modules uses to store, merge and query values. Trees of `Node`s can
be constructed manually or by using the factory function `ConvertToNode` which accepts a 
`map[string]interface{}`.

# License

Copyright 2022 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[ci-img-url]: https://github.com/halimath/appconf/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/appconf
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/appconf
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/appconf
[release-img-url]: https://img.shields.io/github/v/release/halimath/appconf.svg
[release-url]: https://github.com/halimath/appconf/releases