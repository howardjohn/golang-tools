// Package another has another type.
package another

import (
	other "github.com/howardjohn/golang-tools/internal/lsp/references/other"
)

func _() {
	xes := other.GetXes()
	for _, x := range xes {
		_ = x.Y //@mark(anotherXY, "Y"),refs("Y", typeXY, anotherXY, GetXesY)
	}
}
