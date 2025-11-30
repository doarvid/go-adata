package utils

import (
	js "github.com/doarvid/go-adata/common/js"
	"github.com/dop251/goja"
)

func ThsCookie() string {
	vm := goja.New()
	if _, err := vm.RunString(js.ThsJS); err != nil {
		return "v="
	}
	fn, ok := goja.AssertFunction(vm.Get("v"))
	if !ok {
		return "v="
	}
	ret, err := fn(goja.Undefined())
	if err != nil {
		return "v="
	}
	return "v=" + ret.String() + ";"
}
