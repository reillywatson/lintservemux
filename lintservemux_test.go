package lintservemux_test

import (
	"testing"

	"github.com/reillywatson/lintservemux"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, lintservemux.Analyzer, "a")
}
