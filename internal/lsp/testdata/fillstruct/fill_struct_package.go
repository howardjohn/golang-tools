package fillstruct

import (
	h2 "net/http"

	"github.com/howardjohn/golang-tools/internal/lsp/fillstruct/data"
)

func unexported() {
	a := data.B{}   //@suggestedfix("}", "refactor.rewrite")
	_ = h2.Client{} //@suggestedfix("}", "refactor.rewrite")
}
