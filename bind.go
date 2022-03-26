package appconf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

func bind(n *Node, v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("%w: cannot bind to non-pointer %s", ErrInvalidBindingType, reflect.TypeOf(v))
	}

	switch reflect.Indirect(rv).Kind() {
	case reflect.Struct:
		return bindStruct(n, rv)
	case reflect.Map:
		mptr, ok := v.(*map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid map: %t", ErrInvalidBindingType, v)
		}

		if *mptr == nil {
			m := make(map[string]interface{})
			*mptr = m
		}
		return bindMap(n, *mptr)
	default:
		return fmt.Errorf("%w: cannot bind to value of type %q", ErrInvalidBindingType, rv.Type())
	}
}

func bindStruct(n *Node, rv reflect.Value) error {
	rt := reflect.Indirect(rv).Type()

	numFields := rt.NumField()

	for i := 0; i < numFields; i++ {
		f := rt.Field(i)
		opts := determineBindOpts(f)

		if opts.ignore {
			continue
		}

		v, err := resolveReflectValue(n, f.Type, opts)
		if err != nil {
			return fmt.Errorf("%w: struct field %s: %s", ErrInvalidBindingType, f.Name, err)
		}

		rv.Elem().Field(i).Set(v)
	}

	return nil
}

func resolveReflectValue(n *Node, t reflect.Type, opts structFieldBindOpts) (reflect.Value, error) {
	keyPath := ParseKeyPath(opts.key)

	n = n.resolve(keyPath)
	if n == nil {
		return reflect.Value{}, fmt.Errorf("%w: %s", ErrNoSuchKey, opts.key)
	}

	if t == reflect.TypeOf(time.Second) {
		return reflect.ValueOf(n.GetDuration()), nil
	}

	switch t.Kind() {
	case reflect.Struct:
		ptr := reflect.New(t)
		if err := bindStruct(n, ptr); err != nil {
			return reflect.Value{}, err
		}
		return ptr.Elem(), nil

	case reflect.Slice:
		v := reflect.MakeSlice(t, 0, 10)
		v, err := bindSlice(n, v)
		if err != nil {
			return reflect.Value{}, nil
		}
		return v, nil

	case reflect.String,
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Float32,
		reflect.Float64:
		s, err := resolveScalar(n, t.Kind())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(s), nil
	default:
		return reflect.Value{}, fmt.Errorf("%w: type not supported: %s", ErrInvalidBindingType, t)
	}
}

func resolveScalar(n *Node, rk reflect.Kind) (interface{}, error) {
	switch rk {
	case reflect.String:
		return n.GetString(), nil
	case reflect.Bool:
		return n.GetBool(), nil
	case reflect.Int:
		return n.GetInt(), nil
	case reflect.Int8:
		return int8(n.GetInt()), nil
	case reflect.Int16:
		return int16(n.GetInt()), nil
	case reflect.Int32:
		return int32(n.GetInt64()), nil
	case reflect.Int64:
		return n.GetInt64(), nil
	case reflect.Uint:
		return n.GetUint(), nil
	case reflect.Uint8:
		return uint8(n.GetUint()), nil
	case reflect.Uint16:
		return uint16(n.GetUint()), nil
	case reflect.Uint32:
		return uint32(n.GetUint64()), nil
	case reflect.Uint64:
		return n.GetUint64(), nil
	case reflect.Complex64:
		return complex64(n.GetComplex128()), nil
	case reflect.Complex128:
		return n.GetComplex128(), nil
	case reflect.Float32:
		return n.GetFloat32(), nil
	case reflect.Float64:
		return n.GetFloat64(), nil
	default:
		return nil, fmt.Errorf("%w: unsupported kind: %s", ErrNotAScalar, rk)
	}
}

func bindSlice(n *Node, rv reflect.Value) (reflect.Value, error) {
	for idx := 0; idx < len(n.Children); idx++ {
		v, err := resolveReflectValue(n, rv.Type().Elem(), structFieldBindOpts{key: strconv.Itoa(idx)})
		if err != nil {
			if errors.Is(err, ErrNoSuchKey) {
				return rv, nil
			}
			return rv, err
		}
		rv = reflect.Append(rv, v)
	}

	return rv, nil
}

func bindMap(n *Node, m map[string]interface{}) error {
	if len(n.Children) == 0 {
		return nil
	}

	for k, c := range n.Children {
		if len(c.Children) == 0 {
			m[string(k)] = c.Value
		} else {
			cm := make(map[string]interface{})
			bindMap(c, cm)
			m[string(k)] = cm
		}
	}

	return nil
}

type structFieldBindOpts struct {
	key    string
	ignore bool
}

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
