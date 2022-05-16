package c

import "github.com/howardjohn/golang-tools/internal/lsp/rename/b"

func _() {
	b.Hello() //@rename("Hello", "Goodbye")
}
