package mapper

import (
	"time"

	"github.com/corpix/reflect"
	lua "github.com/yuin/gopher-lua"
)

func table(v *lua.LTable) map[lua.LValue]lua.LValue {
	var (
		res = map[lua.LValue]lua.LValue{}
	)

	v.ForEach(
		func(key lua.LValue, value lua.LValue) {
			res[key] = value
		},
	)

	return res
}

func ToValue(gv interface{}) (lua.LValue, error) {
	var (
		err error
	)

	switch v := gv.(type) {
	case nil:
		return lua.LNil, nil
	case bool:
		return lua.LBool(v), nil
	case string:
		return lua.LString(v), nil
	case uint:
		return lua.LNumber(v), nil
	case uint8:
		return lua.LNumber(v), nil
	case uint16:
		return lua.LNumber(v), nil
	case uint32:
		return lua.LNumber(v), nil
	case uint64:
		return lua.LNumber(v), nil
	case int:
		return lua.LNumber(v), nil
	case int8:
		return lua.LNumber(v), nil
	case int16:
		return lua.LNumber(v), nil
	case int32:
		return lua.LNumber(v), nil
	case int64:
		return lua.LNumber(v), nil
	case float32:
		return lua.LNumber(v), nil
	case float64:
		return lua.LNumber(v), nil

	case time.Duration:
		return lua.LNumber(v), nil
	case error:
		return &lua.LUserData{
			Value:     v,
			Env:       nil,
			Metatable: &lua.LTable{},
		}, nil
	}

	// XXX: If you are looking for a way to map reflect.Func into lua.LFunction
	// there is no such mapping here, in this file,
	// because we could convert go func to lua function
	// but not vice versa.
	// (information about arguments got destroyed in the process of creation)
	// So we have separate function to convert go func into lua function.

	var (
		v     = reflect.ValueOf(gv)
		t     = &lua.LTable{}
		value lua.LValue
	)

	switch reflect.TypeOf(gv).Kind() {
	case reflect.Slice, reflect.Array:
		var (
			key = 0
		)

		for key < v.Len() {
			value, err = ToValue(v.Index(key).Interface())
			if err != nil {
				return nil, err
			}

			key++

			t.RawSetInt(key, value)
		}
	case reflect.Map:
		var (
			key lua.LValue
		)

		for _, k := range v.MapKeys() {
			key, err = ToValue(k.Interface())
			if err != nil {
				return nil, err
			}

			value, err = ToValue(v.MapIndex(k).Interface())
			if err != nil {
				return nil, err
			}

			t.RawSetH(key, value)
		}
	default:
		return nil, reflect.NewErrUnknownType(gv)
	}

	return t, nil
}

func FromValue(lv lua.LValue) (interface{}, error) {
	var (
		err error
	)

	switch v := lv.(type) {
	case *lua.LNilType:
		return nil, nil
	case lua.LBool:
		return bool(v), nil
	case lua.LString:
		return string(v), nil
	case lua.LNumber:
		return float64(v), nil
	case *lua.LUserData:
		return v.Value, nil
	case *lua.LTable:
		var (
			n = v.MaxN()
		)

		switch n {
		case 0:
			var (
				res   = make(map[interface{}]interface{})
				key   interface{}
				value interface{}
			)

			for k, v := range table(v) {
				key, err = FromValue(k)
				if err != nil {
					return nil, err
				}

				value, err = FromValue(v)
				if err != nil {
					return nil, err
				}

				res[key] = value
			}

			return res, nil
		default:
			var (
				res   = make([]interface{}, n)
				k     = 0
				value interface{}
			)

			for k < n {
				value, err = FromValue(v.RawGetInt(k + 1))
				if err != nil {
					return nil, err
				}

				res[k] = value

				k++
			}
			return res, nil
		}
	default:
		return nil, reflect.NewErrUnknownType(lv)
	}
}

func MustToValue(v interface{}) lua.LValue {
	var (
		res lua.LValue
		err error
	)

	res, err = ToValue(v)
	if err != nil {
		panic(err)
	}
	return res
}

func MustFromValue(v lua.LValue) interface{} {
	var (
		res interface{}
		err error
	)

	res, err = FromValue(v)
	if err != nil {
		panic(err)
	}
	return res
}
