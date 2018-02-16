package main

import (
	"github.com/corpix/lua/examples/binding/module"
	"github.com/corpix/lua/pool"

	lua "github.com/yuin/gopher-lua"
)

func newLState() *lua.LState {
	var (
		l = lua.NewState()
	)

	l.PreloadModule("module", module.Loader)

	return l
}

func main() {
	var (
		p   = pool.New(newLState)
		l   = p.Get()
		err error
	)
	defer p.Close()
	defer p.Put(l)

	err = l.DoFile("script.lua")
	if err != nil {
		panic(err)
	}
}
