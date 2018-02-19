package mapper

import (
	"errors"
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

func TestToGFunctionInteroperability(t *testing.T) {
	type testCase struct {
		name    string
		code    string
		args    map[string][]interface{}
		exports map[string]func(*testing.T, func(...interface{})) interface{}
		err     error
	}

	type state struct {
		called bool
		args   []interface{}
	}

	var (
		samples = []testCase{
			{
				name: "single return value",
				code: `local m = require("module"); m.fn1(m.fn2())`,
				args: map[string][]interface{}{
					"fn1": []interface{}{1},
					"fn2": nil,
				},
				exports: map[string]func(*testing.T, func(...interface{})) interface{}{
					"fn1": func(t *testing.T, callback func(...interface{})) interface{} {
						return func(x int) {
							callback(x)
						}
					},
					"fn2": func(t *testing.T, callback func(...interface{})) interface{} {
						return func() int {
							callback()
							return 1
						}
					},
				},
				err: nil,
			},
			// FIXME: This one is broken, we should not zip out args with in args automagicaly
			// {
			// 	name: "multiple return value",
			// 	code: `local m = require("module"); m.fn1(m.fn2(), m.fn3())`,
			// 	args: map[string][]interface{}{
			// 		"fn1": []interface{}{1, "hello", 1, "hello"},
			// 		"fn2": nil,
			// 	},
			// 	exports: map[string]func(*testing.T, func(...interface{})) interface{}{
			// 		"fn1": func(t *testing.T, callback func(...interface{})) interface{} {
			// 			return func(x int, y string, xx int, yy string) {
			// 				callback(x, y, xx, yy)
			// 			}
			// 		},
			// 		"fn2": func(t *testing.T, callback func(...interface{})) interface{} {
			// 			return func() (int, string) {
			// 				callback()
			// 				return 1, "hello"
			// 			}
			// 		},
			// 		"fn3": func(t *testing.T, callback func(...interface{})) interface{} {
			// 			return func() (int, string) {
			// 				callback()
			// 				return 1, "hello"
			// 			}
			// 		},
			// 	},
			// 	err: nil,
			// },
			{
				name: "errors",
				code: `local m = require("module"); m.fn1(m.fn2())`,
				args: map[string][]interface{}{
					"fn1": []interface{}{
						errors.New("we bring you an interface{} so you could interface{} while interface{}"),
					},
					"fn2": nil,
				},
				exports: map[string]func(*testing.T, func(...interface{})) interface{}{
					"fn1": func(t *testing.T, callback func(...interface{})) interface{} {
						return func(x error) {
							callback(x)
						}
					},
					"fn2": func(t *testing.T, callback func(...interface{})) interface{} {
						return func() error {
							callback()
							return errors.New("we bring you an interface{} so you could interface{} while interface{}")
						}
					},
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
					l      = lua.NewState()
					s      = map[string]*state{}
					loader func(l *lua.LState) int
					err    error
				)

				loader = func(l *lua.LState) int {
					var (
						exports = map[string]lua.LGFunction{}
						fn      lua.LGFunction
						module  *lua.LTable
						err     error
					)

					for k, v := range sample.exports {
						fn, err = ToGFunction(
							v(
								t,
								func(name string) func(...interface{}) {
									return func(args ...interface{}) {
										s[name] = &state{
											called: true,
											args:   args,
										}
									}
								}(k),
							),
						)
						if err != nil {
							t.Error(err)
							return 0
						}

						exports[k] = fn
					}

					module = l.SetFuncs(l.NewTable(), exports)

					l.Push(module)
					return 1
				}

				l.PreloadModule("module", loader)

				for n := 0; n < 1; n++ {
					for _, v := range s {
						v.called = false
						v.args = []interface{}{}
					}

					err = l.DoString(sample.code)
					if err != nil {
						t.Error(err)
						return
					}

					for k, v := range s {
						assert.Equal(t, true, v.called)
						assert.Equal(t, sample.args[k], v.args)
					}
				}
			},
		)
	}
}
