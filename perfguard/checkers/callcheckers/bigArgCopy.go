package callcheckers

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/quasilyte/go-perfguard/internal/goutil"
	"github.com/quasilyte/go-perfguard/perfguard/checkers"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

func init() {
	doc := checkers.Doc{
		Name:  "bigArgCopy",
		Score: 2,
	}
	checkers.RegisterCallChecker(doc, func() checkers.CallChecker {
		return &bigArgCopyChecker{}
	})
}

type bigArgCopyChecker struct{}

func (c *bigArgCopyChecker) CheckCall(ctx *lint.Context, call *ast.CallExpr) error {
	if len(call.Args) == 0 {
		return nil
	}

	fnType, ok := ctx.TypeOf(call.Fun).(*types.Signature)
	if !ok {
		return nil
	}
	numParams := fnType.Params().Len()
	for i, arg := range call.Args {
		if i >= numParams {
			break
		}
		param := fnType.Params().At(i)
		if param.Name() == "" || param.Name() == "_" {
			continue
		}
		if !c.isBig(ctx, param.Type()) {
			continue
		}
		ctx.Report(lint.ReportParams{
			PosNode: arg,
			Message: fmt.Sprintf("expensive %s arg copy (%d bytes), consider passing it by pointer",
				param.Name(), ctx.Target.Sizes.Sizeof(param.Type())),
			UseFlatSamples: true,
		})
	}

	return nil
}

func (c *bigArgCopyChecker) isBig(ctx *lint.Context, typ types.Type) bool {
	wordSize := ctx.Target.Sizes.Sizeof(types.Typ[types.Uint])
	size := ctx.Target.Sizes.Sizeof(typ)
	numWords := size / wordSize
	if goutil.TypeHasPointers(typ) {
		return numWords > 24
	}
	return numWords > 48
}
