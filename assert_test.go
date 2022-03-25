package appconf

import (
	"testing"

	"github.com/go-test/deep"
)

func assertEqual[T any](t *testing.T, want, got T) {
	if diff := deep.Equal(want, got); diff != nil {
		t.Error(diff)
	}
}
