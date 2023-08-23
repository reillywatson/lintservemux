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
			switch n := n.(type) {
			case *ast.CallExpr:
				switch fnName := n.Fun.(type) {
				case *ast.SelectorExpr:
					if ident, ok := fnName.X.(*ast.Ident); ok {
						if ident.Name != "http" {
							return true
						}
					} else {
						return true
					}
					if fnName.Sel != nil {
						switch fnName.Sel.Name {
						case "Handle", "HandleFunc":
							reportNodef(pass, n, "http.%s uses DefaultServeMux", fnName.Sel.Name)
						case "ListenAndServe", "ServeTLS", "Serve":
							if len(n.Args) >= 2 {
								if ident, ok := n.Args[1].(*ast.Ident); ok {
									if ident.Name == "nil" {
										reportNodef(pass, n, "http.%s should pass an http.Handler", fnName.Sel.Name)
									}
								}
							}
						}
					}
				}
			case *ast.CompositeLit:
				if t, ok := n.Type.(*ast.SelectorExpr); ok {
					if ident, ok := t.X.(*ast.Ident); ok {
						if ident.Name == "http" {
							if t.Sel != nil && t.Sel.Name == "Server" {
								foundHandlerArg := false
								for _, el := range n.Elts {
									if kv, ok := el.(*ast.KeyValueExpr); ok {
										if nm, ok := kv.Key.(*ast.Ident); ok {
											if nm.Name == "Handler" {
												foundHandlerArg = true
												if val, ok := kv.Value.(*ast.Ident); ok {
													if val.Name == "nil" {
														reportNodef(pass, val, "http.Server should include a non-nil Handler")
													}
												}
											}
										}
									}
								}
								if !foundHandlerArg {
									reportNodef(pass, n, "http.Server should set a Handler")
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
