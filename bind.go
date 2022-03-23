package appconf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	FieldTagKey            = "appconf"
	FieldTagValueSeparator = ","
	FieldTagIgnore         = "ignore"
)

var (
	ErrInvalidBindingType = errors.New("invalid binding type")
)

func (c *AppConfig) Bind(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("%w: cannot bind to non-pointer %s", ErrInvalidBindingType, reflect.TypeOf(v))
	}

	switch reflect.Indirect(rv).Kind() {
	case reflect.Struct:
		return c.bindStruct(rv)
	case reflect.Map:
		return c.bindMap(rv)
	default:
		return fmt.Errorf("%w: cannot bind to value of type %q", ErrInvalidBindingType, rv.Type())
	}
}

func (c *AppConfig) bindStruct(rv reflect.Value) error {
	rt := reflect.Indirect(rv).Type()

	n := rt.NumField()

	for i := 0; i < n; i++ {
		f := rt.Field(i)
		opts := determineBindOpts(f)

		if opts.ignore {
			// No matching field tag; skip field
			continue
		}

		var v reflect.Value

		if f.Type == reflect.TypeOf(time.Second) {
			v = reflect.ValueOf(c.GetDuration(opts.key))
		} else {
			switch f.Type.Kind() {
			case reflect.Struct:
				ptr := reflect.New(f.Type)
				c.Sub(opts.key).bindStruct(ptr)
				v = ptr.Elem()

			case reflect.String:
				v = reflect.ValueOf(c.GetString(opts.key))
			case reflect.Bool:
				v = reflect.ValueOf(c.GetBool(opts.key))
			case reflect.Int:
				v = reflect.ValueOf(c.GetInt(opts.key))
			case reflect.Int8:
				v = reflect.ValueOf(int8(c.GetInt(opts.key)))
			case reflect.Int16:
				v = reflect.ValueOf(int16(c.GetInt(opts.key)))
			case reflect.Int32:
				v = reflect.ValueOf(int32(c.GetInt64(opts.key)))
			case reflect.Int64:
				v = reflect.ValueOf(c.GetInt64(opts.key))
			case reflect.Uint:
				v = reflect.ValueOf(c.GetUint(opts.key))
			case reflect.Uint8:
				v = reflect.ValueOf(uint8(c.GetUint(opts.key)))
			case reflect.Uint16:
				v = reflect.ValueOf(uint16(c.GetUint(opts.key)))
			case reflect.Uint32:
				v = reflect.ValueOf(uint32(c.GetUint64(opts.key)))
			case reflect.Uint64:
				v = reflect.ValueOf(c.GetUint64(opts.key))
			case reflect.Complex64:
				v = reflect.ValueOf(complex64(c.GetComplex128(opts.key)))
			case reflect.Complex128:
				v = reflect.ValueOf(c.GetComplex128(opts.key))
			case reflect.Float32,
				reflect.Float64:
				v = reflect.ValueOf(c.GetFloat64(opts.key))
			default:
				return fmt.Errorf("%w: unsupported struct field %s: type not supported: %s", ErrInvalidBindingType, f.Name, f.Type)
			}
		}

		rv.Elem().Field(i).Set(v)
	}

	return nil
}

func (c *AppConfig) bindMap(rv reflect.Value) error {
	return fmt.Errorf("not implemented")
}

type (
	structFieldBindOpts struct {
		key    string
		ignore bool
	}
)

func determineBindOpts(f reflect.StructField) structFieldBindOpts {

	opts := structFieldBindOpts{
		key: strings.ToLower(f.Name),
	}

	t := f.Tag.Get(FieldTagKey)
	parts := strings.Split(t, FieldTagValueSeparator)

	for i, p := range parts {
		p = strings.TrimSpace(p)
		if i == 0 && len(p) > 0 {
			opts.key = parts[0]
		} else if i > 0 && p == FieldTagIgnore {
			opts.ignore = true
		}
	}

	return opts
}
