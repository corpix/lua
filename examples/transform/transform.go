package main

import (
	"github.com/davecgh/go-spew/spew"
	lua "github.com/yuin/gopher-lua"

	"github.com/corpix/lua/mapper"
	"github.com/corpix/lua/pool"
)

func newLState() *lua.LState {
	var (
		l   = lua.NewState()
		err error
	)

	err = l.DoString(
		`
            function transform(v)
                v["hello"] = "world"
                return v
            end
        `,
	)
	if err != nil {
		panic(err)
	}

	return l
}

func main() {
	var (
		input = map[string]interface{}{
			"how you": "doin?",
		}
		output      interface{}
		p           = pool.New(newLState)
		l           = p.Get()
		inputValue  lua.LValue
		outputValue lua.LValue
		err         error
	)
	defer p.Close()
	defer p.Put(l)

	inputValue, err = mapper.ToValue(input)
	if err != nil {
		panic(err)
	}

	err = l.CallByParam(
		lua.P{
			Fn:      l.GetGlobal("transform"),
			NRet:    1,
			Protect: true,
		},
		inputValue,
	)
	if err != nil {
		panic(err)
	}

	outputValue = l.Get(-1)
	l.Pop(1)

	output, err = mapper.FromValue(outputValue)
	if err != nil {
		panic(err)
	}

	spew.Dump(
		input,
		output,
	)
}
