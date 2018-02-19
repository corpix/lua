package mapper

import (
	"errors"
	"testing"

	"github.com/corpix/reflect"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestToValue(t *testing.T) {
	type testCase struct {
		name   string
		input  interface{}
		output lua.LValue
		err    error
	}

	var (
		samples = []testCase{
			{
				name:   "nil",
				input:  nil,
				output: lua.LNil,
				err:    nil,
			},
			{
				name:   "bool",
				input:  false,
				output: lua.LBool(false),
				err:    nil,
			},
			{
				name:   "string",
				input:  "",
				output: lua.LString(""),
				err:    nil,
			},
			{
				name:   "uint64",
				input:  uint64(0),
				output: lua.LNumber(0),
				err:    nil,
			},
			{
				name:   "int64",
				input:  int64(0),
				output: lua.LNumber(0),
				err:    nil,
			},
			{
				name:   "float64",
				input:  float64(0),
				output: lua.LNumber(0),
				err:    nil,
			},
			{
				name:  "error",
				input: errors.New("hello, this is a bad way to create error"),
				output: &lua.LUserData{
					Value:     errors.New("hello, this is a bad way to create error"),
					Env:       nil,
					Metatable: &lua.LTable{},
				},
				err: nil,
			},
			{
				name:   "table",
				input:  map[interface{}]interface{}{},
				output: &lua.LTable{},
				err:    nil,
			},
			{
				name:  "array",
				input: [2]string{"one", "two"},
				output: func() *lua.LTable {
					t := &lua.LTable{}
					t.RawSetInt(1, lua.LString("one"))
					t.RawSetInt(2, lua.LString("two"))
					return t
				}(),
				err: nil,
			},
			{
				name:  "slice",
				input: []string{"one", "two"},
				output: func() *lua.LTable {
					t := &lua.LTable{}
					t.RawSetInt(1, lua.LString("one"))
					t.RawSetInt(2, lua.LString("two"))
					return t
				}(),
				err: nil,
			},
			{
				name: "map",
				input: map[interface{}]interface{}{
					"foo":      "bar",
					float64(1): "baz",
				},
				output: func() *lua.LTable {
					t := &lua.LTable{}
					t.RawSetString("foo", lua.LString("bar"))
					t.RawSetH(lua.LNumber(1), lua.LString("baz"))
					return t
				}(),
				err: nil,
			},
			func() testCase {
				var (
					f = func() {}
				)

				return testCase{
					name:   "unknown type error",
					input:  f,
					output: nil,
					err:    reflect.NewErrUnknownType(f),
				}
			}(),
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					v   lua.LValue
					err error
				)

				v, err = ToValue(sample.input)
				assert.IsType(t, sample.err, err)
				assert.Equal(t, sample.output, v)
			},
		)
	}
}

func TestFromValue(t *testing.T) {
	type testCase struct {
		name   string
		input  lua.LValue
		output interface{}
		err    error
	}

	var (
		samples = []testCase{
			{
				name:   "nil",
				input:  lua.LNil,
				output: nil,
				err:    nil,
			},
			{
				name:   "bool",
				input:  lua.LBool(false),
				output: false,
				err:    nil,
			},
			{
				name:   "string",
				input:  lua.LString(""),
				output: "",
				err:    nil,
			},
			{
				name:   "number",
				input:  lua.LNumber(0),
				output: float64(0),
				err:    nil,
			},
			{
				name: "error",
				input: &lua.LUserData{
					Value:     errors.New("hello, this is a bad way to create error"),
					Env:       nil,
					Metatable: &lua.LTable{},
				},
				output: errors.New("hello, this is a bad way to create error"),
				err:    nil,
			},
			{
				name:   "table",
				input:  &lua.LTable{},
				output: map[interface{}]interface{}{},
				err:    nil,
			},
			{
				name: "array",
				input: func() *lua.LTable {
					t := &lua.LTable{}
					t.RawSetInt(1, lua.LString("one"))
					t.RawSetInt(2, lua.LString("two"))
					return t
				}(),
				output: []interface{}{"one", "two"},
				err:    nil,
			},
			{
				name: "map",
				input: func() *lua.LTable {
					t := &lua.LTable{}
					t.RawSetString("foo", lua.LString("bar"))
					t.RawSetH(lua.LNumber(1), lua.LString("baz"))
					return t
				}(),
				output: map[interface{}]interface{}{
					"foo":      "bar",
					float64(1): "baz",
				},
				err: nil,
			},
			func() testCase {
				var (
					f = &lua.LFunction{}
				)

				return testCase{
					name:   "unknown type error",
					input:  f,
					output: nil,
					err:    reflect.NewErrUnknownType(f),
				}
			}(),
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					v   interface{}
					err error
				)

				v, err = FromValue(sample.input)
				assert.IsType(t, sample.err, err)
				assert.Equal(t, sample.output, v)
			},
		)
	}
}
