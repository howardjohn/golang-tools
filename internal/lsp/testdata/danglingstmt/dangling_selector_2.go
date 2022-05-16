package danglingstmt

import "github.com/howardjohn/golang-tools/internal/lsp/foo"

func _() {
	foo. //@rank(" //", Foo)
	var _ = []string{foo.} //@rank("}", Foo)
}
