package nodisk

import (
	"github.com/howardjohn/golang-tools/internal/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}
