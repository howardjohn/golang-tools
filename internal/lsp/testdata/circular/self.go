package circular

import (
	_ "github.com/howardjohn/golang-tools/internal/lsp/circular" //@diag("_ \"github.com/howardjohn/golang-tools/internal/lsp/circular\"", "compiler", "import cycle not allowed", "error"),diag("\"github.com/howardjohn/golang-tools/internal/lsp/circular\"", "compiler", "could not import github.com/howardjohn/golang-tools/internal/lsp/circular (no package for import github.com/howardjohn/golang-tools/internal/lsp/circular)", "error")
)
