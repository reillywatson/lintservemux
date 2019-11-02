package main

import (
	"github.com/reillywatson/lintservemux"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(lintservemux.Analyzer) }
