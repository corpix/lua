package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

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
					err:    NewErrUnknownType(f),
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
				assert.EqualValues(t, sample.err, err)
				assert.EqualValues(t, sample.output, v)
			},
		)
	}
}
