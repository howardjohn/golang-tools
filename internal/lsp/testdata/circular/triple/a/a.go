package a

import (
	_ "github.com/howardjohn/golang-tools/internal/lsp/circular/triple/b" //@diag("_ \"github.com/howardjohn/golang-tools/internal/lsp/circular/triple/b\"", "compiler", "import cycle not allowed", "error")
)
