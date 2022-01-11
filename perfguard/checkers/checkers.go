package checkers

import (
	"fmt"
	"go/ast"

	"github.com/quasilyte/go-perfguard/internal/resolve"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

type callCheckerInfo struct {
	doc Doc
	new func() CallChecker
}

var (
	callCheckers = make(map[string]callCheckerInfo)
)

type Doc struct {
	Name  string
	Score int
}

type CallContext struct {
	*lint.Context

	Sym resolve.CallInfo
}

type CallChecker interface {
	CheckCall(ctx *CallContext, call *ast.CallExpr) error
}

type PackageChecker interface {
	CheckPackage(ctx *lint.Context, files []lint.SourceFile) error
}

func RegisterCallChecker(doc Doc, constructor func() CallChecker) {
	if _, ok := callCheckers[doc.Name]; ok {
		panic(fmt.Sprintf("%s call checker is already registered", doc.Name))
	}
	callCheckers[doc.Name] = callCheckerInfo{
		doc: doc,
		new: constructor,
	}
}

func Create(filter func(doc Doc) bool) []PackageChecker {
	callChecker := &callcheckerWalker{}
	for _, c := range callCheckers {
		if filter(c.doc) {
			callChecker.checkers = append(callChecker.checkers, c.new())
			callChecker.tags = append(callChecker.tags, c.doc.Name)
		}
	}

	var result []PackageChecker
	result = append(result, callChecker)

	return result
}
