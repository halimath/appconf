package appconf

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
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

func ConvertToNode(m map[string]interface{}) (*Node, error) {
	n := &Node{
		Children: make(map[Key]*Node),
	}

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
		return &Node{Value: fmt.Sprint(val)}, nil
	default:
		return nil, fmt.Errorf("%w: unsupported value type: %t", ErrUnsupportedValue, val)
	}
}
