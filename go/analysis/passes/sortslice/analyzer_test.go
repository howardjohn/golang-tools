package sortslice_test

import (
	"testing"

	"github.com/howardjohn/golang-tools/go/analysis/analysistest"
	"github.com/howardjohn/golang-tools/go/analysis/passes/sortslice"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, sortslice.Analyzer, "a")
}
