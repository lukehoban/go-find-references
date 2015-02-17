package main

import (
	"code.google.com/p/rog-go/exp/go/ast"
	"code.google.com/p/rog-go/exp/go/parser"
	"code.google.com/p/rog-go/exp/go/token"
	"code.google.com/p/rog-go/exp/go/types"
)

var fileSet = types.FileSet
var pkgScope = ast.NewScope(parser.Universe)

// Parses the given file and returns its AST.
func parseAST(filename string) *ast.File {
	f, err := parser.ParseFile(fileSet, filename, nil, 0, pkgScope)
	if err != nil {
		panic(err)
	}
	return f
}

// Looks up the identifier in the given ast and offset
// and returns its name and the position of its declaration in fs.
// Returns an invalid Position when either is not found.
func findDecl(f *ast.File, offset int) (string, token.Position) {
	containsCursor := func(node ast.Node) bool {
		from := fileSet.Position(node.Pos()).Offset
		to := fileSet.Position(node.End()).Offset
		return offset >= from && offset < to
	}

	// traverse the ast tree until we find a node at the given offset position
	var ident ast.Expr
	ast.Inspect(f, func(node ast.Node) bool {
		switch expr := node.(type) {
		case *ast.SelectorExpr:
			if containsCursor(expr) && containsCursor(expr.Sel) {
				ident = expr
			}
		case *ast.Ident:
			if containsCursor(expr) {
				ident = expr
			}
		}
		return ident == nil
	})

	if ident == nil {
		return "", token.Position{}
	}

	obj, _ := types.ExprType(ident, types.DefaultImporter)
	if obj == nil {
		return "", token.Position{}
	}
	return obj.Name, fileSet.Position(types.DeclPos(obj))
}

// Scans the given ast for identifiers that are declared at the given
// declPos and returns their positions.
func findReferences(f ast.Node, declPos token.Position, identName string) chan token.Position {
	out := make(chan token.Position)

	check := func(ident ast.Expr) {
		obj, _ := types.ExprType(ident, types.DefaultImporter)
		if obj == nil {
			return
		}

		dp := fileSet.Position(types.DeclPos(obj))
		if dp == declPos {
			out <- fileSet.Position(ident.Pos())
		}
	}

	go func() {
		ast.Inspect(f, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.SelectorExpr:
				if node.(*ast.SelectorExpr).Sel.Name == identName {
					check(node.(ast.Expr))
				}
			case *ast.Ident:
				if node.(*ast.Ident).Name == identName {
					check(node.(ast.Expr))
				}
			}
			return true
		})
		close(out)
	}()

	return out
}
