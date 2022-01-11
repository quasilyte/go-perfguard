package callcheckers

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/quasilyte/go-perfguard/internal/resolve"
	"github.com/quasilyte/go-perfguard/internal/typeis"
	"github.com/quasilyte/go-perfguard/perfguard/checkers"
)

func init() {
	doc := checkers.Doc{
		Name:  "bytesToStringFmt",
		Score: 2,
	}
	checkers.RegisterCallChecker(doc, func() checkers.CallChecker {
		return &BytesToStringFmtChecker{}
	})
}

type BytesToStringFmtChecker struct{}

func (c *BytesToStringFmtChecker) CheckCall(ctx *checkers.CallContext, call *ast.CallExpr) error {
	if call.Ellipsis.IsValid() {
		return nil // Skip variadic calls
	}
	if ctx.Sym.PkgPath != "fmt" {
		return nil
	}

	formatArgNum := 0
	switch ctx.Sym.FuncName {
	case "Fprintf":
		formatArgNum = 1
	case "Sprintf", "Printf":
		// OK, argNum is 0
	default:
		return nil
	}
	if len(call.Args) <= formatArgNum+1 {
		return nil
	}

	formatArg, ok := call.Args[formatArgNum].(*ast.BasicLit)
	if !ok || formatArg.Kind != token.STRING {
		return nil
	}
	formatString, err := strconv.Unquote(formatArg.Value)
	if err != nil {
		return nil
	}
	if formatString == "%s" {
		// TODO: #9
		return nil // It's redundantSprint case
	}
	formatInfo, ok := resolve.FmtString(formatString)
	if !ok {
		return nil
	}

	for _, argInfo := range formatInfo.Args {
		if argInfo.Verb != 's' && argInfo.Verb != 'q' {
			continue
		}
		argIndex := argInfo.ArgNum + formatArgNum + 1
		if argIndex >= len(call.Args) {
			break
		}
		convExpr := call.Args[argIndex]
		convInfo := resolve.ConvExpr(ctx.Target.Types, convExpr)
		if convInfo.Arg == nil || !typeis.String(convInfo.DstType) {
			continue
		}
		srcType := ctx.TypeOf(convInfo.Arg)
		if !typeis.ByteSlice(srcType.Underlying()) {
			continue
		}
		ctx.SuggestNode(convExpr, convInfo.Arg)
	}

	return nil
}
