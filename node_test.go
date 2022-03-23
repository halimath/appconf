package appconf

import (
	"testing"
	"time"

	"github.com/go-test/deep"
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

	if diff := deep.Equal(want, n); diff != nil {
		t.Error(diff)
	}
}

func TestConvertToNode(t *testing.T) {
	in := map[string]interface{}{
		"foo":     "bar",
		"timeout": time.Second,
		"spam": map[string]interface{}{
			"eggs": "ham",
		},
		"spam.salad": "none",
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
		},
	}

	got, err := ConvertToNode(in)
	if err != nil {
		t.Fatal(err)
	}

	want.Dump(0)
	got.Dump(0)

	if diff := deep.Equal(want, got); diff != nil {
		t.Error(diff)
	}
}
