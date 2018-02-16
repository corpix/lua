package module

import (
	"fmt"

	luamapper "github.com/corpix/lua/mapper"
	"github.com/davecgh/go-spew/spew"
	"github.com/yuin/gopher-lua"
)

func Loader(l *lua.LState) int {
	fooBar, err := luamapper.ToGFunction(
		func(a int, b int, c map[string]string) (string, string, []interface{}) {
			return "foo", "bar", []interface{}{a, b, c}
		},
	)
	if err != nil {
		panic(err)
	}

	dump, err := luamapper.ToGFunction(spew.Dump)
	if err != nil {
		panic(err)
	}
	sdump, err := luamapper.ToGFunction(spew.Sdump)
	if err != nil {
		panic(err)
	}
	printf, err := luamapper.ToGFunction(fmt.Printf)
	if err != nil {
		panic(err)
	}

	exports := map[string]lua.LGFunction{
		"foo_bar": fooBar,
		"sdump":   sdump,
		"dump":    dump,
		"printf":  printf,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)

	return 1
}
