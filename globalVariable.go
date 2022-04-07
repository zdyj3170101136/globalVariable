package globalVariable

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

const Doc = `find global variable.

find the global variable that modified not in init func.`

var Analyzer = &analysis.Analyzer{
	Name:     "globalVariable",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	allGlobalVariableMap := make(map[*ast.Object]*ast.Ident)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		if strings.Contains(pass.Fset.Position(n.Pos()).Filename, "_test.go") {
			// filter testing file
			return
		}
		if scope, ok := pass.TypesInfo.Scopes[n]; ok {
			ast.Inspect(n, func(node ast.Node) bool {
				// *ast.AssignStmt only load the statement in func
				g, ok := node.(*ast.GenDecl)
				if !ok || g.Tok != token.VAR {
					return true
				}
				if innerMost := scope.Innermost(g.Pos()); innerMost == scope {
					// declaration in the file scope is global variable
					for _, spec := range g.Specs {
						s := spec.(*ast.ValueSpec)
						for _, name := range s.Names {
							allGlobalVariableMap[name.Obj] = name
						}
					}
				}
				return true
			})
		}
	})

	modifiedGlobalVariable := make(map[*ast.Ident][]*ast.Ident, 0)
	nodeFilter = []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	// find the global variable that had changed in another place
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		stmt := node.(*ast.AssignStmt)
		for _, lhs := range stmt.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				if ident2 := allGlobalVariableMap[ident.Obj]; ident2 != nil && ident2 != ident {
					modifiedGlobalVariable[ident2] = append(modifiedGlobalVariable[ident2], ident)
				}
			}

		}
	})

	// filter variable in init func
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.FuncDecl); ok && decl.Name.Name == "init" {
				funcScope := pass.TypesInfo.Scopes[decl.Type]

				for k, modifies := range modifiedGlobalVariable {
					newModifies := make([]*ast.Ident, 0, len(modifies))
					for _, modify := range modifies {
						if !funcScope.Contains(modify.Pos()) {
							newModifies = append(newModifies, modify)
						} else {
							pass.ReportRangef(modify, "in init func")
						}
					}
					if len(newModifies) > 0 {
						modifiedGlobalVariable[k] = newModifies
					} else {
						delete(modifiedGlobalVariable, k)
					}
				}
			}
		}
	}

	for k, modifies := range modifiedGlobalVariable {
		pass.ReportRangef(k, "global variable")
		for _, modify := range modifies {
			pass.ReportRangef(modify, "modify")
		}
	}

	return nil, nil
}
