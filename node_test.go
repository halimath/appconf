package appconf

import (
	"fmt"
	"testing"
	"time"

	"github.com/halimath/assertthat-go/assert"
	"github.com/halimath/assertthat-go/is"
)

func TestNodeOverwriteWith(t *testing.T) {
	n := &Node{
		Children: map[Key]*Node{
			"foo": {Value: "bar"},
			"spam": {
				Children: map[Key]*Node{
					"eggs": {Value: "e"},
				},
			},
		},
	}

	o := &Node{
		Children: map[Key]*Node{
			"spam": {
				Children: map[Key]*Node{
					"eggs": {Value: "ham"},
				},
			},
		},
	}

	want := &Node{
		Children: map[Key]*Node{
			"foo": {Value: "bar"},
			"spam": {
				Children: map[Key]*Node{
					"eggs": {Value: "ham"},
				},
			},
		},
	}

	n.OverwriteWith(o)

	assert.That(t, n, is.DeepEqual(want))
}

func TestConvertToNode(t *testing.T) {
	in := map[string]interface{}{
		"foo":     "bar",
		"timeout": time.Second,
		"spam": map[string]interface{}{
			"eggs": "ham",
		},
		"spam.salad": "none",
		"slice":      []interface{}{"a", "b", "c"},
	}

	want := &Node{
		Children: map[Key]*Node{
			"foo":     {Value: "bar"},
			"timeout": {Value: "1s"},
			"spam": {
				Children: map[Key]*Node{
					"eggs":  {Value: "ham"},
					"salad": {Value: "none"},
				},
			},
			"slice": {
				Children: map[Key]*Node{
					"0": {Value: "a"},
					"1": {Value: "b"},
					"2": {Value: "c"},
				},
			},
		},
	}

	got, err := ConvertToNode(in)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Want")
	want.Dump(0)
	fmt.Println("Got")
	got.Dump(0)

	assert.That(t, got, is.DeepEqual(want))
}
