package resolve

import (
	"go/ast"
	"go/types"
)

type CallInfo struct {
	PkgName  string
	PkgPath  string
	FuncName string
}

func Call(typesInfo *types.Info, call *ast.CallExpr) CallInfo {
	var info CallInfo

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return info
	}
	pkgIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return info
	}
	obj := typesInfo.ObjectOf(pkgIdent)
	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		return info
	}

	info.PkgName = pkgIdent.Name
	info.FuncName = selector.Sel.String()
	info.PkgPath = pkgName.Imported().Path()
	return info
}

type ConvInfo struct {
	DstType types.Type
	Arg     ast.Expr
}

func ConvExpr(typesInfo *types.Info, e ast.Expr) ConvInfo {
	var info ConvInfo

	call, ok := e.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return info
	}
	typ := typesInfo.TypeOf(e)
	if typ == nil {
		return info
	}
	info.DstType = typ
	info.Arg = call.Args[0]
	return info
}

func SplitFuncName(fn *ast.FuncDecl) (typeName, funcName string) {
	if fn == nil {
		return "", ""
	}
	funcName = fn.Name.Name
	if fn.Recv != nil && len(fn.Recv.List) != 0 {
		typeName = getTypeName(fn.Recv.List[0].Type)
	}
	return typeName, funcName
}

func getTypeName(typeExpr ast.Expr) string {
	switch typ := typeExpr.(type) {
	case *ast.Ident:
		return typ.Name
	case *ast.StarExpr:
		return getTypeName(typ.X)
	case *ast.ParenExpr:
		return getTypeName(typ.X)

	default:
		return ""
	}
}
