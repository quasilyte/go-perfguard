package lint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/types"
	"strings"
)

type Context struct {
	Target *Target

	tag  string
	warn func(Warning)
}

func (ctx *Context) SetWarnFunc(fn func(Warning)) {
	ctx.warn = fn
}

func (ctx *Context) SetTag(tag string) {
	ctx.tag = tag
}

func (ctx *Context) SuggestNode(oldNode, newNode ast.Node) {
	startPos := ctx.Target.Fset.Position(oldNode.Pos())

	var b strings.Builder
	b.Write(ctx.NodeText(oldNode))
	b.WriteString(" => ")
	replacement := ctx.NodeText(newNode)
	b.Write(replacement)
	message := strings.ReplaceAll(b.String(), "\n", `\n`)

	ctx.warn(Warning{
		Filename: startPos.Filename,
		Line:     startPos.Line,
		Tag:      ctx.tag,
		Text:     message,
		Fix: &QuickFix{
			From:        oldNode.Pos(),
			To:          oldNode.End(),
			Replacement: replacement,
		},
	})
}

func (ctx *Context) NodeText(n ast.Node) []byte {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, ctx.Target.Fset, n); err != nil {
		return nil
	}
	return buf.Bytes()
}

func (ctx *Context) Report(n ast.Node, format string, args ...interface{}) {
	startPos := ctx.Target.Fset.Position(n.Pos())

	var message string
	if len(args) != 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}
	message = strings.ReplaceAll(message, "\n", `\n`)

	ctx.warn(Warning{
		Filename: startPos.Filename,
		Line:     startPos.Line,
		Tag:      ctx.tag,
		Text:     message,
	})
}

// TypeOf returns the type of expression x.
//
// Unlike TypesInfo.TypeOf, it never returns nil.
// Instead, it returns the Invalid type as a sentinel UnknownType value.
func (ctx *Context) TypeOf(x ast.Expr) types.Type {
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
