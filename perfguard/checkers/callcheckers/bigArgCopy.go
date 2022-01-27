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
	c.checkRecv(ctx, call)

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

func (c *bigArgCopyChecker) checkRecv(ctx *lint.Context, call *ast.CallExpr) {
	methodExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	arg, ok := methodExpr.X.(*ast.Ident)
	if !ok {
		return
	}
	obj := ctx.VarOf(arg)
	typ := obj.Type()
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}
	named, ok := typ.(*types.Named)
	if !ok {
		return
	}
	// TODO: is there a better way to find the method?
	// types.LookupFieldOrMethod does linear search.
	result, _, _ := types.LookupFieldOrMethod(named, true, ctx.Target.Pkg, methodExpr.Sel.Name)
	methodObject, ok := result.(*types.Func)
	if !ok {
		return
	}
	methodType := methodObject.Type().(*types.Signature)
	recv := methodType.Recv()
	if recv == nil {
		return
	}
	if !c.isBig(ctx, recv.Type()) {
		return
	}
	ctx.Report(lint.ReportParams{
		PosNode: arg,
		Message: fmt.Sprintf("expensive %s receiver copy (%d bytes), consider passing it by pointer",
			arg.Name, ctx.Target.Sizes.Sizeof(recv.Type())),
		UseFlatSamples: true,
	})
}

func (c *bigArgCopyChecker) isBig(ctx *lint.Context, typ types.Type) bool {
	wordSize := ctx.Target.Sizes.Sizeof(types.Typ[types.Uint])
	size := ctx.Target.Sizes.Sizeof(typ)
	numWords := size / wordSize
	if goutil.TypeHasPointers(typ) {
		return numWords > 32
	}
	return numWords > 64
}
