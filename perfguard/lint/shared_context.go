package lint

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"

	"github.com/quasilyte/go-perfguard/internal/resolve"
	"github.com/quasilyte/perf-heatmap/heatmap"
)

type SharedContext struct {
	Target *Target

	Heatmap *heatmap.Index

	Filename string // Filename is a name of file that is being analyzed
	TypeName string // TypeName is a receiver name of the current func
	FuncName string // FuncName is a current func/method name

	// Sym is a currently analyzed function call symbol info.
	// Only relevant for funccall checkers.
	Sym resolve.CallInfo

	Warn func(Warning)
}

func (ctx *SharedContext) NodeText(n ast.Node) []byte {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, ctx.Target.Fset, n); err != nil {
		return nil
	}
	return buf.Bytes()
}

func (ctx *SharedContext) Position(n ast.Node) token.Position {
	return ctx.Target.Fset.Position(n.Pos())
}

func (ctx *SharedContext) EndPosition(n ast.Node) token.Position {
	return ctx.Target.Fset.Position(n.End())
}

func (ctx *SharedContext) ObjectOf(x *ast.Ident) types.Object {
	obj := ctx.Target.Types.ObjectOf(x)
	if obj != nil {
		return obj
	}
	return UnknownVar
}

func (ctx *SharedContext) VarOf(x *ast.Ident) *types.Var {
	obj := ctx.Target.Types.ObjectOf(x)
	if obj != nil {
		if v, ok := obj.(*types.Var); ok {
			return v
		}
	}
	return UnknownVar
}

// TypeOf returns the type of expression x.
//
// Unlike TypesInfo.TypeOf, it never returns nil.
// Instead, it returns the Invalid type as a sentinel UnknownType value.
func (ctx *SharedContext) TypeOf(x ast.Expr) types.Type {
	typ := ctx.Target.Types.TypeOf(x)
	if typ != nil {
		return typ
	}
	// Usually it means that some incorrect type info was loaded
	// or the analyzed package was only partially (?) correct.
	// To avoid nil pointer panics we can return a sentinel value
	// that will fail most type assertions as well as kind checks
	// (if the call side expects a *types.Basic).
	return UnknownType
}

// UnknownType is a special sentinel value that is returned from the CheckerContext.TypeOf
// method instead of the nil type.
var UnknownType types.Type = types.Typ[types.Invalid]

var UnknownVar = types.NewVar(
	token.NoPos,
	types.NewPackage("unknown", "unknown"),
	"unknown",
	UnknownType)
