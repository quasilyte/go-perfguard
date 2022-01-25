package goutil

import (
	"go/ast"
	"go/types"
)

func TypeHasPointers(typ types.Type) bool {
	switch typ := typ.(type) {
	case *types.Basic:
		switch typ.Kind() {
		case types.UnsafePointer, types.String, types.UntypedNil, types.UntypedString:
			return true
		}
		return false

	case *types.Named:
		return TypeHasPointers(typ.Underlying())

	case *types.Struct:
		for i := 0; i < typ.NumFields(); i++ {
			if TypeHasPointers(typ.Field(i).Type()) {
				return true
			}
		}
		return false

	case *types.Array:
		return TypeHasPointers(typ.Elem())

	default:
		return true
	}
}

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
