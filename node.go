package appconf

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	KeySeparator = "."
)

var (
	ErrUnsupportedValue = errors.New("unsupported value")
	ErrNoSuchKey        = errors.New("no such key")
	ErrNotAScalar       = errors.New("not a scalar value")

	keyFilterRegexp = regexp.MustCompile(`[^0-9a-z]`)
)

type Key string

func NormalizeKey(k string) Key {
	return Key(keyFilterRegexp.ReplaceAllString(strings.ToLower(k), ""))
}

type KeyPath []Key

func (p KeyPath) Join() string {
	var b strings.Builder
	for i, e := range p {
		if i > 0 {
			b.WriteString(KeySeparator)
		}
		b.WriteString(string(e))
	}
	return b.String()
}

func ParseKeyPath(s string) KeyPath {
	parts := strings.Split(s, KeySeparator)
	path := make(KeyPath, len(parts))
	for i, p := range parts {
		path[i] = NormalizeKey(p)
	}

	return path
}

type Node struct {
	Value    string
	Children map[Key]*Node
}

func NewNode(v string) *Node {
	return &Node{
		Value:    v,
		Children: make(map[Key]*Node),
	}
}

func (n *Node) resolve(path KeyPath) *Node {
	if len(path) == 0 {
		return n
	}

	v, ok := n.Children[path[0]]
	if !ok {
		return nil
	}
	return v.resolve(path[1:])
}

func (n *Node) OverwriteWith(o *Node) {
	n.Value = o.Value
	for key, node := range o.Children {
		found, ok := n.Children[key]
		if !ok {
			n.Children[key] = node
		} else {
			found.OverwriteWith(node)
		}
	}
}

func (n *Node) Dump(indent int) {
	fmt.Printf("%v\n", n.Value)

	indent += 2

	for k, v := range n.Children {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Print(k, ": ")
		v.Dump(indent)
	}
}

func (n *Node) GetString() string {
	v, _ := n.GetStringE()
	return v
}

func (n *Node) GetStringE() (string, error) {
	if len(n.Children) != 0 {
		return "", ErrNotAScalar
	}

	return n.Value, nil
}

func (n *Node) GetInt() int {
	v, _ := n.GetIntE()
	return v
}

func (n *Node) GetIntE() (int, error) {
	v, err := n.GetInt64E()
	return int(v), err
}

func (n *Node) GetInt64() int64 {
	v, _ := n.GetInt64E()
	return v
}

func (n *Node) GetInt64E() (int64, error) {
	v, err := n.GetStringE()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(v, 10, 64)
}

func (n *Node) GetUint() uint {
	v, _ := n.GetUintE()
	return v
}

func (n *Node) GetUintE() (uint, error) {
	v, err := n.GetUint64E()
	return uint(v), err
}

func (n *Node) GetUint64() uint64 {
	v, _ := n.GetUint64E()
	return v
}

func (n *Node) GetUint64E() (uint64, error) {
	v, err := n.GetStringE()
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(v, 10, 64)
}

func (n *Node) GetFloat32() float32 {
	v, _ := n.GetFloat32E()
	return v
}

func (n *Node) GetFloat32E() (float32, error) {
	v, err := n.GetStringE()
	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(v, 32)
	return float32(f), err
}

func (n *Node) GetFloat64() float64 {
	v, _ := n.GetFloat64E()
	return v
}

func (n *Node) GetFloat64E() (float64, error) {
	v, err := n.GetStringE()
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(v, 64)
}

func (n *Node) GetComplex128() complex128 {
	v, _ := n.GetComplex128E()
	return v
}

func (n *Node) GetComplex128E() (complex128, error) {
	v, err := n.GetStringE()
	if err != nil {
		return 0, err
	}

	return strconv.ParseComplex(v, 128)
}

func (n *Node) GetBool() bool {
	v, _ := n.GetBoolE()
	return v
}

func (n *Node) GetBoolE() (bool, error) {
	v, err := n.GetStringE()
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(v)
}

func (n *Node) GetDuration() time.Duration {
	d, _ := n.GetDurationE()
	return d
}

func (n *Node) GetDurationE() (time.Duration, error) {
	v, err := n.GetStringE()
	if err != nil {
		return 0, err
	}

	return time.ParseDuration(v)
}

func ConvertToNode(m map[string]interface{}) (*Node, error) {
	n := NewNode("")

	for k, val := range m {
		path := ParseKeyPath(k)

		var key Key
		var valueNode *Node
		var err error

		if len(path) > 1 {
			key = path[0]
			valueNode, err = ConvertToNode(map[string]interface{}{
				path[1:].Join(): val, // TODO: This is not very efficient!
			})
		} else {
			valueNode, err = createNodeFromValue(val)
			key = path[0]
		}

		if err != nil {
			return nil, err
		}

		existing, ok := n.Children[key]
		if ok {
			existing.OverwriteWith(valueNode)
		} else {
			n.Children[key] = valueNode
		}
	}

	return n, nil
}

func createNodeFromValue(val interface{}) (*Node, error) {
	t := reflect.TypeOf(val)
	switch t.Kind() {
	case reflect.Slice:
		s, ok := val.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: nested slice invalid: %v", ErrUnsupportedValue, val)
		}
		return createNodeFromSlice(s)
	case reflect.Map:
		nested, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: nested map is not a ConfigMap: %v", ErrUnsupportedValue, val)
		}
		return ConvertToNode(nested)
	case reflect.Bool,
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
		reflect.Float64,
		reflect.String:
		return NewNode(fmt.Sprint(val)), nil
	default:
		return nil, fmt.Errorf("%w: unsupported value type: %T", ErrUnsupportedValue, val)
	}
}

func createNodeFromSlice(val []interface{}) (*Node, error) {
	r := NewNode("")

	for idx, v := range val {
		n, err := createNodeFromValue(v)
		if err != nil {
			return nil, err
		}
		r.Children[Key(strconv.Itoa(idx))] = n
	}

	return r, nil
}
