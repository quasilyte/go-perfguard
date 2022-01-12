package checkers

import (
	"fmt"
	"go/ast"

	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

type callCheckerInfo struct {
	doc Doc
	new func() CallChecker
}

type funcCheckerInfo struct {
	doc Doc
	new func() FuncChecker
}

var (
	callCheckers = make(map[string]callCheckerInfo)
	funcCheckers = make(map[string]funcCheckerInfo)
)

type Doc struct {
	Name  string
	Score int

	OptLevel int

	NeedsProfile bool
}

type CallChecker interface {
	CheckCall(ctx *lint.Context, call *ast.CallExpr) error
}

type FuncChecker interface {
	CheckFunc(ctx *lint.Context, body *ast.BlockStmt) error
}

type PackageChecker interface {
	CheckPackage(ctx *lint.SharedContext, files []lint.SourceFile) error
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

func RegisterFuncChecker(doc Doc, constructor func() FuncChecker) {
	if _, ok := callCheckers[doc.Name]; ok {
		panic(fmt.Sprintf("%s func checker is already registered", doc.Name))
	}
	funcCheckers[doc.Name] = funcCheckerInfo{
		doc: doc,
		new: constructor,
	}
}

func minHeatLevel(doc *Doc) int {
	if doc.OptLevel == 2 {
		return 5
	}
	return 1
}

func Create(filter func(doc Doc) bool) []PackageChecker {
	callChecker := &callcheckerWalker{}
	for _, c := range callCheckers {
		if filter(c.doc) {
			callChecker.checkers = append(callChecker.checkers, callcheckerWithContext{
				ctx: lint.NewContext(c.doc.Name, minHeatLevel(&c.doc)),
				obj: c.new(),
			})
		}
	}

	funcChecker := &funccheckerWalker{}
	for _, c := range funcCheckers {
		if filter(c.doc) {
			funcChecker.checkers = append(funcChecker.checkers, funccheckerWithContext{
				ctx: lint.NewContext(c.doc.Name, minHeatLevel(&c.doc)),
				obj: c.new(),
			})
		}
	}

	var result []PackageChecker
	result = append(result,
		callChecker,
		funcChecker)

	return result
}
