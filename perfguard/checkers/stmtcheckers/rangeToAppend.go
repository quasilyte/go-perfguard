package stmtcheckers

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/go-toolsmith/astcast"
	"github.com/go-toolsmith/astequal"
	"github.com/quasilyte/go-perfguard/internal/goutil"
	"github.com/quasilyte/go-perfguard/internal/typeis"
	"github.com/quasilyte/go-perfguard/perfguard/checkers"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

// TODO: make it ruleguard rule again.
// But we need to handle the "same types" condition somehow.
//
// m.Match(`for _, $x := range $src { $dst = append($dst, $x) }`).
// 		Where(m["src"].Type.Is(`[]$_`) && !m["dst"].Contains(`$x`)).
// 		Suggest(`$dst = append($dst, $src...)`).
// 		Report(`for ... { ... } => $dst = append($dst, $src...)`)

func init() {
	doc := checkers.Doc{
		Name:  "rangeToAppend",
		Score: 3,
	}
	checkers.RegisterStmtChecker(doc, func() checkers.StmtChecker {
		return &rangeToAppendChecker{}
	})
}

type rangeToAppendChecker struct {
}

func (c *rangeToAppendChecker) CheckStmt(ctx *lint.Context, n ast.Stmt) error {
	rng, ok := n.(*ast.RangeStmt)
	if !ok {
		return nil
	}
	if rng.Value == nil || astcast.ToIdent(rng.Key).Name != "_" {
		return nil
	}
	if len(rng.Body.List) != 1 {
		return nil
	}
	rangeExprType := ctx.TypeOf(rng.X)
	if !typeis.Slice(rangeExprType) {
		return nil
	}
	assign, ok := rng.Body.List[0].(*ast.AssignStmt)
	if !ok || assign.Tok != token.ASSIGN || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return nil
	}
	lhs := assign.Lhs[0]
	rhs := assign.Rhs[0]
	lhsType := ctx.TypeOf(lhs)
	if !types.Identical(rangeExprType, lhsType) {
		return nil
	}
	call, ok := rhs.(*ast.CallExpr)
	if !ok {
		return nil
	}
	called, ok := call.Fun.(*ast.Ident)
	if !ok || called.Name != "append" || len(call.Args) != 2 {
		return nil
	}
	if !astequal.Expr(call.Args[0], lhs) || !astequal.Expr(call.Args[1], rng.Value) {
		return nil
	}
	if goutil.ContainsIdent(lhs, rng.Value.(*ast.Ident)) {
		return nil
	}
	ctx.SuggestNode(lint.SuggestParams{
		Message: ctx.Sprintf("for … { … } => %s = append(%s, %s...)", lhs, lhs, rng.X),
		OldNode: rng,
		NewNode: &ast.AssignStmt{
			Lhs: []ast.Expr{lhs},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun:      &ast.Ident{Name: "append"},
					Args:     []ast.Expr{lhs, rng.Value},
					Ellipsis: 1,
				},
			},
		},
	})
	return nil
}
