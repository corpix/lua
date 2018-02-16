local m = require("module")
print("FooBar in Go called from Lua:", m.foo_bar(1, 2, { a = 2}))
print("spew.Sdump in Go called from Lua:", m.sdump(m.foo_bar(1, 2, { a = 2})))
m.printf("fmt.Printf in Go called from Lua: %s -> %s -> %#v\n", m.foo_bar(1, 2, { a = 2}))
