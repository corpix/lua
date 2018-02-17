package mapper

import (
	"github.com/corpix/reflect"
	lua "github.com/yuin/gopher-lua"
)

func ToGFunction(gfn interface{}) (lua.LGFunction, error) {
	var (
		v     = reflect.IndirectValue(reflect.ValueOf(gfn))
		t     = v.Type()
		numIn = t.NumIn()
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
				n    = l.GetTop()
				tt   reflect.Type
				arg  interface{}
				args = make(
					[]reflect.Value,
					n,
					n,
				)
				rt  lua.LValue
				rts []reflect.Value
				err error
			)

			if numIn != n {
				if !t.IsVariadic() {
					l.ArgError(
						n,
						reflect.NewErrTooFewArguments(numIn, n).Error(),
					)
					return 0
				}

				numIn--
			}

			for k := 0; k < numIn; k++ {
				tt = t.In(k)
				arg, err = FromValue(l.Get(k + 1))
				if err != nil {
					l.ArgError(k+1, err.Error())
					return 0
				}
				arg, err = reflect.ConvertToType(
					arg,
					tt,
				)
				if err != nil {
					l.ArgError(k+1, err.Error())
					return 0
				}

				args[k] = reflect.ValueOf(arg)
			}

			if t.IsVariadic() {
				tt = t.In(numIn).Elem()
				for k := numIn; k < n; k++ {
					arg, err = FromValue(l.Get(k + 1))
					if err != nil {
						l.ArgError(k+1, err.Error())
						return 0
					}
					arg, err = reflect.ConvertToType(
						arg,
						tt,
					)
					if err != nil {
						l.ArgError(k+1, err.Error())
						return 0
					}

					args[k] = reflect.ValueOf(arg)
				}
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
