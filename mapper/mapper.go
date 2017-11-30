package mapper

import (
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
				k     = 1
				value interface{}
			)

			for k <= n {
				value, err = FromValue(v.RawGetInt(k))
				if err != nil {
					return nil, err
				}

				res[k-1] = value

				k++
			}
			return res, nil
		}
	default:
		return nil, NewErrUnknownType(v)
	}
}
