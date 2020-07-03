package middlemux

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMiddleMux(t *testing.T) {
	mm := NewMiddleMux()

	require.NotNil(t, mm.ServeMux, "ServeMux ought to be initialized.")
	assert.NotNil(t, mm.handlerFn, "handlerFn ought to be initialized.")
	assert.Equal(t,
		reflect.ValueOf(mm.handlerFn).Pointer(),
		reflect.ValueOf(mm.ServeMux.ServeHTTP).Pointer(),
		"handlerFn ought to be ServeMux's ServeHTTP method initially.",
	)
}
