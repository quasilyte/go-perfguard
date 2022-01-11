package checkers

import (
	"go/ast"

	"github.com/quasilyte/go-perfguard/internal/resolve"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

type callcheckerWalker struct {
	callContext CallContext

	checkers []CallChecker
	tags     []string
}

func (d *callcheckerWalker) CheckPackage(ctx *lint.Context, files []lint.SourceFile) error {
	d.callContext.Context = ctx

	var checkError error

	for _, f := range files {
		ast.Inspect(f.Syntax, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			d.callContext.Sym = resolve.Call(ctx.Target.Types, call)
			for i, c := range d.checkers {
				ctx.SetTag(d.tags[i])
				if err := c.CheckCall(&d.callContext, call); err != nil {
					checkError = err
				}
			}
			return true
		})
	}

	return checkError
}
