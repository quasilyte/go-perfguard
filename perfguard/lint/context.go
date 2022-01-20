package lint

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"time"

	"github.com/quasilyte/perf-heatmap/heatmap"
)

type Context struct {
	*SharedContext

	tag          string
	minHeatLevel int
}

type NodeReplacement struct {
	Text   []byte
	Syntax ast.Node
}

func NewContext(tag string, minHeatLevel int) Context {
	return Context{
		tag:          tag,
		minHeatLevel: minHeatLevel,
	}
}

type MultiChangeSuggestParams struct {
	ReportPos     token.Pos
	ReportMessage string

	OldNodes []ast.Node
	NewNodes []NodeReplacement
	HotNodes []ast.Node
}

func (ctx *Context) MultiChangeSuggest(params MultiChangeSuggestParams) {
	var hotNodes = params.HotNodes
	if len(hotNodes) == 0 {
		hotNodes = params.OldNodes
	}
	samplesValue, matched := ctx.listMatchesHeatmap(hotNodes)
	if !matched {
		return
	}

	reportPos := ctx.Target.Fset.Position(params.ReportPos)
	message := params.ReportMessage

	textEdits := make([]TextEdit, 0, len(params.OldNodes))
	for i, oldNode := range params.OldNodes {
		newNode := params.NewNodes[i]
		replacement := newNode.Text
		if newNode.Syntax != nil {
			replacement = ctx.NodeText(newNode.Syntax)
		}
		textEdits = append(textEdits, TextEdit{
			From:        oldNode.Pos(),
			To:          oldNode.End(),
			Replacement: replacement,
		})
	}

	ctx.Warn(Warning{
		Filename:    reportPos.Filename,
		Line:        reportPos.Line,
		Tag:         ctx.tag,
		Text:        message,
		Fixes:       textEdits,
		SamplesTime: time.Duration(samplesValue),
	})
}

type SuggestParams struct {
	OldNode ast.Node
	NewNode ast.Node

	Message string

	HotNodes []ast.Node
}

func (ctx *Context) SuggestNode(params SuggestParams) {
	oldNode := params.OldNode
	newNode := params.NewNode

	var hotNodes = params.HotNodes
	if len(hotNodes) == 0 {
		hotNodes = []ast.Node{oldNode}
	}
	samplesValue, matched := ctx.listMatchesHeatmap(hotNodes)
	if !matched {
		return
	}

	startPos := ctx.Target.Fset.Position(oldNode.Pos())

	message := params.Message
	replacement := ctx.NodeText(newNode)
	if message == "" {
		var b strings.Builder
		b.Write(ctx.NodeText(oldNode))
		b.WriteString(" => ")
		b.Write(replacement)
		message = strings.ReplaceAll(b.String(), "\n", `\n`)
	}

	textEdit := TextEdit{
		From:        oldNode.Pos(),
		To:          oldNode.End(),
		Replacement: replacement,
	}
	ctx.Warn(Warning{
		Filename:    startPos.Filename,
		Line:        startPos.Line,
		Tag:         ctx.tag,
		Text:        message,
		Fixes:       []TextEdit{textEdit},
		SamplesTime: time.Duration(samplesValue),
	})
}

type ReportParams struct {
	PosNode ast.Node

	Message string

	HotNodes []ast.Node
}

func (ctx *Context) Report(params ReportParams) {
	startPos := ctx.Target.Fset.Position(params.PosNode.Pos())

	var hotNodes = params.HotNodes
	if len(hotNodes) == 0 {
		hotNodes = []ast.Node{params.PosNode}
	}
	samplesValue, matched := ctx.listMatchesHeatmap(hotNodes)
	if !matched {
		return
	}

	message := strings.ReplaceAll(params.Message, "\n", `\n`)

	ctx.Warn(Warning{
		Filename:    startPos.Filename,
		Line:        startPos.Line,
		Tag:         ctx.tag,
		Text:        message,
		SamplesTime: time.Duration(samplesValue),
	})
}

func (ctx *Context) listMatchesHeatmap(nodes []ast.Node) (int64, bool) {
	samplesValue := int64(0)
	matched := false
	for _, heatNode := range nodes {
		if ctx.matchesHeatmap(heatNode, &samplesValue) {
			matched = true
		}
	}
	return samplesValue, matched
}

func (ctx *Context) matchesHeatmap(n ast.Node, samplesValue *int64) bool {
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
	totalValue := int64(0)
	ctx.Heatmap.QueryLineRange(key, lineFrom, lineTo, func(l heatmap.LineStats) bool {
		if l.GlobalHeatLevel >= minLevel {
			isHot = true
		}
		totalValue += l.Value
		return true
	})
	if isHot {
		*samplesValue += totalValue
	}
	return isHot
}
