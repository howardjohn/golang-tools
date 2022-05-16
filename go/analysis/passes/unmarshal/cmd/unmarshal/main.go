// The unmarshal command runs the unmarshal analyzer.
package main

import (
	"github.com/howardjohn/golang-tools/go/analysis/passes/unmarshal"
	"github.com/howardjohn/golang-tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(unmarshal.Analyzer) }
