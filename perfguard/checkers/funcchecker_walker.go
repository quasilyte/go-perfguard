package checkers

import (
	"go/ast"

	"github.com/quasilyte/go-perfguard/internal/resolve"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

type funccheckerWithContext struct {
	ctx lint.Context
	obj FuncChecker
}

type funccheckerWalker struct {
	checkers []funccheckerWithContext
}

func (w *funccheckerWalker) CheckPackage(ctx *lint.SharedContext, files []lint.SourceFile) error {
	for i := range w.checkers {
		w.checkers[i].ctx.SharedContext = ctx
	}

	var checkError error

	for _, f := range files {
		ctx.Filename = ctx.Position(f.Syntax).Filename
		ast.Inspect(f.Syntax, func(n ast.Node) bool {
			var body *ast.BlockStmt
			switch n := n.(type) {
			case *ast.FuncDecl:
				body = n.Body
				ctx.TypeName, ctx.FuncName = resolve.SplitFuncName(n)
			case *ast.FuncLit:
				body = n.Body
				ctx.TypeName = ""
				ctx.FuncName = ""
			}
			if body != nil {
				for _, c := range w.checkers {
					if err := c.obj.CheckFunc(&c.ctx, body); err != nil {
						checkError = err
					}
				}
			}
			return true
		})
	}

	return checkError
}
