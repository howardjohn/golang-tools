// The findcall command runs the findcall analyzer.
package main

import (
	"github.com/howardjohn/golang-tools/go/analysis/passes/findcall"
	"github.com/howardjohn/golang-tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(findcall.Analyzer) }
