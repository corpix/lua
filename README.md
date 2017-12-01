lua
---------

[![Build Status](https://travis-ci.org/corpix/lua.svg?branch=master)](https://travis-ci.org/corpix/lua)

Helpers around https://github.com/yuin/gopher-lua

## Examples

Here is a simple example which shows how to transform data from Lua:

``` console
$ go run ./examples/transform/transform.go
(map[string]interface {}) (len=1) {
 (string) (len=7) "how you": (string) (len=5) "doin?"
}
(map[interface {}]interface {}) (len=2) {
 (string) (len=7) "how you": (string) (len=5) "doin?",
 (string) (len=5) "hello": (string) (len=5) "world"
}
```
