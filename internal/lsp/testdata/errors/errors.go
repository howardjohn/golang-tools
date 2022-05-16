package errors

import (
	"github.com/howardjohn/golang-tools/internal/lsp/types"
)

func _() {
	bob.Bob() //@complete(".")
	types.b   //@complete(" //", Bob_interface)
}
