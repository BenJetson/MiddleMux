package middlemux

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestWithMiddleware(t *testing.T) {
	mm := NewMiddleMux()

	t.Run("Nil", func(t *testing.T) {
		// Ensure that nil Middleware has no side effects.
		assert.NotPanics(t, func() { mm.Use(nil) })
		assert.NotNil(t, mm.handlerFn, "handlerFn ought to be initialized.")
		assert.Equal(t,
			reflect.ValueOf(mm.handlerFn).Pointer(),
			reflect.ValueOf(mm.ServeMux.ServeHTTP).Pointer(),
			"handlerFn ought to be ServeMux's ServeHTTP method initially.",
		)
	})

	// Setup routes for test.
	mm.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "serving root") // nolint: errcheck
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", bytes.NewReader(nil))

	// Sanity check the handler before running tests.
	mm.ServeHTTP(w, r)
	assert.Equal(t,
		"serving root\n",
		w.Body.String(),
	)

	t.Run("WithOneWrapper", func(t *testing.T) {
		w = httptest.NewRecorder()

		mm.Use(func(h http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "in the middle") // nolint: errcheck

				h(w, r)
			}
		})

		mm.ServeHTTP(w, r)

		assert.Equal(t,
			"in the middle\nserving root\n",
			w.Body.String(),
		)
	})

	t.Run("WithSecondWrapper", func(t *testing.T) {
		w = httptest.NewRecorder()

		mm.Use(func(h http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "first for now") // nolint: errcheck

				h(w, r)
			}
		})

		mm.ServeHTTP(w, r)

		assert.Equal(t,
			"first for now\nin the middle\nserving root\n",
			w.Body.String(),
		)
	})

	t.Run("WithManyWrappers", func(t *testing.T) {
		w = httptest.NewRecorder()

		expect := "first for now\nin the middle\nserving root\n"

		for i := 0; i < 39; i++ {
			extra := fmt.Sprintf("number %d\n", i)
			expect = extra + expect

			mm.Use(func(h http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, extra) // nolint: errcheck

					h(w, r)
				}
			})
		}

		mm.ServeHTTP(w, r)

		assert.Equal(t, expect, w.Body.String())
	})
}
