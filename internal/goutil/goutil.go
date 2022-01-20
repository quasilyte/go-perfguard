package goutil

import (
	"go/ast"
)

func ContainsIdent(root ast.Node, x *ast.Ident) bool {
	return Contains(root, func(n ast.Node) bool {
		if y, ok := n.(*ast.Ident); ok {
			return x.Name == y.Name
		}
		return false
	})
}

func Contains(root ast.Node, fn func(n ast.Node) bool) bool {
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if found {
			return false
		}
		if fn(n) {
			found = true
			return false
		}
		return true
	})
	return found
}
