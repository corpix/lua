package mapper

import (
	"github.com/corpix/reflect"
	lua "github.com/yuin/gopher-lua"
)

func ToGFunction(gfn interface{}) (lua.LGFunction, error) {
	var (
		v        = reflect.IndirectValue(reflect.ValueOf(gfn))
		t        = v.Type()
		numIn    = t.NumIn()
		variadic = t.IsVariadic()
	)

	if k := t.Kind(); k != reflect.Func {
		return nil, reflect.NewErrWrongKind(
			reflect.Func,
			k,
		)
	}

	return lua.LGFunction(
		func(l *lua.LState) int {
			var (
				n          = l.GetTop()
				fixedNumIn = numIn
				buf        interface{}
				args       []reflect.Value
				tt         reflect.Type
				rt         lua.LValue
				rts        []reflect.Value
				err        error
			)

			if variadic {
				fixedNumIn--

				if n < fixedNumIn {
					l.ArgError(
						n,
						reflect.NewErrWrongArgumentsQuantity(fixedNumIn, n).Error(),
					)
				}
			} else {
				if n != fixedNumIn {
					l.ArgError(
						n,
						reflect.NewErrWrongArgumentsQuantity(fixedNumIn, n).Error(),
					)
					return 0
				}
			}

			args = make([]reflect.Value, n)

			for k := n; k > 0; k-- {
				buf, err = FromValue(l.Get(k))
				if err != nil {
					l.ArgError(k, err.Error())
					return 0
				}

				switch {
				case k <= fixedNumIn:
					tt = t.In(k - 1)
				case variadic && k == n:
					tt = t.In(fixedNumIn).Elem()
				}

				buf, err = reflect.ConvertToType(
					buf,
					tt,
				)
				if err != nil {
					l.ArgError(k, err.Error())
					return 0
				}

				args[k-1] = reflect.ValueOf(buf)
			}

			rts = v.Call(args)

			for _, r := range rts {
				rt, err = ToValue(r.Interface())
				if err != nil {
					l.RaiseError("%s", err)
					return len(rts)
				}
				l.Push(rt)
			}

			return len(rts)
		},
	), nil
}

func MustToGFunction(fn interface{}) lua.LGFunction {
	var (
		res lua.LGFunction
		err error
	)

	res, err = ToGFunction(fn)
	if err != nil {
		panic(err)
	}
	return res
}
