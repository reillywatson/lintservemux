package lintservemux

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const Doc = `check for uses of http.DefaultServeMux`

var Analyzer = &analysis.Analyzer{
	Doc:      Doc,
	Name:     "lintservemux",
	Run:      lintservemuxCheck,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func lintservemuxCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch fn := n.(type) {
			case *ast.CallExpr:
				switch fnName := fn.Fun.(type) {
				case *ast.SelectorExpr:
					if ident, ok := fnName.X.(*ast.Ident); ok {
						if ident.Name != "http" {
							return true
						}
					}
					if fnName.Sel != nil {
						switch fnName.Sel.Name {
						case "Handle", "HandleFunc":
							reportNodef(pass, n, "http.%s uses DefaultServeMux", fnName.Sel.Name)
						case "ListenAndServe", "ServeTLS", "Serve":
							if len(fn.Args) >= 2 {
								if ident, ok := fn.Args[1].(*ast.Ident); ok {
									if ident.Name == "nil" {
										reportNodef(pass, n, "http.%s should pass an http.Handler", fnName.Sel.Name)
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func reportNodef(pass *analysis.Pass, node ast.Node, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	pass.Report(analysis.Diagnostic{Pos: node.Pos(), End: node.End(), Message: msg})
}
