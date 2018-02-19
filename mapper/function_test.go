package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestToGFunction(t *testing.T) {
	type testCase struct {
		name   string
		code   string
		args   []interface{}
		fn     func(*testing.T, func()) interface{}
		result interface{}
		err    error
	}

	var (
		samples = []testCase{
			{
				name: "without arguments",
				code: `require("module").fn()`,
				args: []interface{}{},
				fn: func(t *testing.T, callback func()) interface{} {
					return func() {
						callback()
					}
				},
				result: nil,
				err:    nil,
			},
			{
				name: "with arguments",
				code: `require("module").fn(1)`,
				args: []interface{}{},
				fn: func(t *testing.T, callback func()) interface{} {
					return func(a int) {
						assert.Equal(t, 1, a)
						callback()
					}
				},
				err: nil,
			},
			{
				name: "with variadic arguments",
				code: `require("module").fn(1, 2, 3, 4)`,
				args: []interface{}{},
				fn: func(t *testing.T, callback func()) interface{} {
					return func(a ...int) {
						assert.Equal(t, []int{1, 2, 3, 4}, a)
						callback()
					}
				},
				err: nil,
			},
			{
				name: "with mixed arguments",
				code: `require("module").fn(1, 2, 3, 4)`,
				args: []interface{}{},
				fn: func(t *testing.T, callback func()) interface{} {
					return func(a int, b ...int) {
						assert.Equal(t, 1, a)
						assert.Equal(t, []int{2, 3, 4}, b)
						callback()
					}
				},
				err: nil,
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					l     = lua.NewState()
					state = struct {
						called bool
					}{}
					loader func(l *lua.LState) int
					err    error
				)

				loader = func(l *lua.LState) int {
					var (
						fn      lua.LGFunction
						exports map[string]lua.LGFunction
						module  *lua.LTable
						err     error
					)

					fn, err = ToGFunction(
						sample.fn(
							t,
							func() {
								state.called = true
							},
						),
					)
					if err != nil {
						t.Error(err)
						return 0
					}

					exports = map[string]lua.LGFunction{"fn": fn}
					module = l.SetFuncs(l.NewTable(), exports)

					l.Push(module)
					return 1
				}

				l.PreloadModule("module", loader)

				for n := 0; n < 10; n++ {
					state.called = false

					err = l.DoString(sample.code)
					if err != nil {
						t.Error(err)
						return
					}

					assert.Equal(t, true, state.called)
				}
			},
		)
	}
}
