package lint

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/quasilyte/perf-heatmap/heatmap"
)

type Context struct {
	*SharedContext

	tag          string
	minHeatLevel int
}

func NewContext(tag string, minHeatLevel int) Context {
	return Context{
		tag:          tag,
		minHeatLevel: minHeatLevel,
	}
}

type SuggestParams struct {
	OldNode ast.Node
	NewNode ast.Node

	HotNodes []ast.Node
}

func (ctx *Context) SuggestNode(params SuggestParams) {
	oldNode := params.OldNode
	newNode := params.NewNode

	if len(params.HotNodes) == 0 {
		if !ctx.matchesHeatmap(oldNode) {
			return
		}
	} else {
		matches := false
		for _, heatNode := range params.HotNodes {
			if ctx.matchesHeatmap(heatNode) {
				matches = true
				break
			}
		}
		if !matches {
			return
		}
	}

	startPos := ctx.Target.Fset.Position(oldNode.Pos())

	var b strings.Builder
	b.Write(ctx.NodeText(oldNode))
	b.WriteString(" => ")
	replacement := ctx.NodeText(newNode)
	b.Write(replacement)
	message := strings.ReplaceAll(b.String(), "\n", `\n`)

	ctx.Warn(Warning{
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

func (ctx *Context) Report(n ast.Node, format string, args ...interface{}) {
	startPos := ctx.Target.Fset.Position(n.Pos())

	var message string
	if len(args) != 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}
	message = strings.ReplaceAll(message, "\n", `\n`)

	ctx.Warn(Warning{
		Filename: startPos.Filename,
		Line:     startPos.Line,
		Tag:      ctx.tag,
		Text:     message,
	})
}

func (ctx *Context) matchesHeatmap(n ast.Node) bool {
	if ctx.Heatmap == nil {
		return true
	}
	minLevel := ctx.minHeatLevel
	if minLevel == 0 {
		return true
	}
	startPos := ctx.Target.Fset.Position(n.Pos())
	endPos := ctx.Target.Fset.Position(n.End())
	lineFrom := startPos.Line
	lineTo := endPos.Line
	isHot := false
	key := heatmap.Key{
		TypeName: ctx.TypeName,
		FuncName: ctx.FuncName,
		Filename: filepath.Base(startPos.Filename),
		PkgName:  ctx.Target.Pkg.Name(),
	}
	ctx.Heatmap.QueryLineRange(key, lineFrom, lineTo, func(line int, level heatmap.HeatLevel) bool {
		if level.Global >= minLevel {
			isHot = true
			return false
		}
		return true
	})
	return isHot
}
