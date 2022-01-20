package checkers

import (
	"go/ast"

	"github.com/quasilyte/go-perfguard/internal/resolve"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

type stmtcheckerWithContext struct {
	ctx lint.Context
	obj StmtChecker
}

type stmtcheckerWalker struct {
	checkers []stmtcheckerWithContext
}

func (w *stmtcheckerWalker) CheckPackage(ctx *lint.SharedContext, files []lint.SourceFile) error {
	for i := range w.checkers {
		w.checkers[i].ctx.SharedContext = ctx
	}

	var checkError error

	for _, f := range files {
		ctx.Filename = ctx.Position(f.Syntax).Filename
		ast.Inspect(f.Syntax, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.FuncDecl:
				ctx.TypeName, ctx.FuncName = resolve.SplitFuncName(n)
			case *ast.FuncLit:
				ctx.TypeName = ""
				ctx.FuncName = ""
			case ast.Stmt:
				for _, c := range w.checkers {
					if err := c.obj.CheckStmt(&c.ctx, n); err != nil {
						checkError = err
					}
				}
			}
			return true
		})
	}

	return checkError
}
