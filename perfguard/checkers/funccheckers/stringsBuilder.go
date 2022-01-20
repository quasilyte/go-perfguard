package funccheckers

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/quasilyte/go-perfguard/internal/typeis"
	"github.com/quasilyte/go-perfguard/perfguard/checkers"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

func init() {
	doc := checkers.Doc{
		Name:  "stringsBuilder",
		Score: 4,
	}
	checkers.RegisterFuncChecker(doc, func() checkers.FuncChecker {
		return &stringsBuilderChecker{
			stringsBuilderExpr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "strings"},
				Sel: &ast.Ident{Name: "Builder"},
			},
			allowedUsages: make(map[*ast.Ident]struct{}),
		}
	})
}

type stringsBuilderChecker struct {
	ctx                *lint.Context
	candidates         []stringsBuilderCandidate
	stringsBuilderExpr ast.Expr
	allowedUsages      map[*ast.Ident]struct{}
	nestedFunc         bool
}

type stringsBuilderCandidate struct {
	typeNode ast.Expr

	obj *types.Var

	toStringUsages []ast.Node
}

func (c *stringsBuilderChecker) CheckFunc(ctx *lint.Context, body *ast.BlockStmt) error {
	c.ctx = ctx
	c.candidates = c.candidates[:0]
	c.nestedFunc = false
	for k := range c.allowedUsages {
		delete(c.allowedUsages, k)
	}

	ast.Inspect(body, c.walk)

	for _, candidate := range c.candidates {
		if len(candidate.toStringUsages) == 0 {
			continue
		}
		ctx.SuggestNode(lint.SuggestParams{
			OldNode:  candidate.typeNode,
			NewNode:  c.stringsBuilderExpr,
			HotNodes: candidate.toStringUsages,
		})
	}

	return nil
}

func (c *stringsBuilderChecker) walk(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.FuncLit:
		nestedFunc := c.nestedFunc
		c.nestedFunc = true
		ast.Inspect(n.Body, c.walk)
		c.nestedFunc = nestedFunc
		return false

	case *ast.ValueSpec:
		if c.nestedFunc {
			return true
		}
		// Track `var $x bytes.Buffer` variables.
		if len(n.Names) != 1 || n.Values != nil || n.Type == nil {
			return true
		}
		name := n.Names[0]
		obj := c.ctx.VarOf(name)
		if c.isBytesBuffer(obj.Type()) {
			c.track(obj, n.Type)
			return false
		}

	case *ast.AssignStmt:
		if c.nestedFunc {
			return true
		}
		// Track `$x := bytes.Buffer{}` variables.
		// Track `$x := &bytes.Buffer{}` variables.
		if len(n.Lhs) != 1 || len(n.Rhs) != 1 || n.Tok != token.DEFINE {
			return true
		}
		var typeNode ast.Expr
		rhs := n.Rhs[0]
		if unary, ok := rhs.(*ast.UnaryExpr); ok && unary.Op == token.AND {
			rhs = unary.X
		}
		if lit, ok := rhs.(*ast.CompositeLit); ok {
			if len(lit.Elts) != 0 || lit.Type == nil {
				return true
			}
			asSelector, ok := lit.Type.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if pkgIdent, ok := asSelector.X.(*ast.Ident); !ok || pkgIdent.Name != "bytes" {
				return true
			}
			if asSelector.Sel.Name != "Buffer" {
				return true
			}
			typeNode = lit.Type
		}
		name, ok := n.Lhs[0].(*ast.Ident)
		if !ok {
			return true
		}
		obj := c.ctx.VarOf(name)
		if typeNode != nil && c.isBytesBuffer(obj.Type()) {
			c.track(obj, typeNode)
			return false
		}

	case *ast.CallExpr:
		c.markCallArgs(n)
		selector, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		x, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}
		obj := c.ctx.VarOf(x)
		if !c.isBytesBuffer(obj.Type()) {
			return true
		}
		candidate := c.find(obj)
		if candidate != nil {
			c.allowedUsages[x] = struct{}{}
			return c.checkBufferCall(n, candidate, selector.Sel.Name)
		}
		return true

	case *ast.Ident:
		// If it's a buffer usage, untrack the candidate.
		// We stop tree traversal when we recognize the usage.
		obj := c.ctx.VarOf(n)
		if c.isBytesBuffer(obj.Type()) {
			if _, ok := c.allowedUsages[n]; !ok {
				c.untrack(obj, "uncategorized usage")
			}
		}
	}

	return true
}

func (c *stringsBuilderChecker) checkBufferCall(n *ast.CallExpr, candidate *stringsBuilderCandidate, methodName string) bool {
	// "Reset" method is not allowed as it doesn't work that well in
	// strings buider. Buffer will reuse the memory while builder will not.
	switch methodName {
	case "Cap", "Len", "Grow", "Write", "WriteString", "WriteByte", "WriteRune":
		return true

	case "String":
		candidate.toStringUsages = append(candidate.toStringUsages, n.Fun)
		return true

	default:
		c.untrack(candidate.obj, "unsupported method")
		return true
	}
}

func (c *stringsBuilderChecker) markCallArgs(n *ast.CallExpr) {
	if n.Ellipsis.IsValid() {
		return
	}

	fnType, ok := c.ctx.TypeOf(n.Fun).(*types.Signature)
	if !ok {
		return
	}
	params := fnType.Params()

	for i, arg := range n.Args {
		if i >= params.Len() {
			break
		}

		x := arg
		if arg, ok := arg.(*ast.UnaryExpr); ok && arg.Op == token.AND {
			x = arg.X
		}
		ident, ok := x.(*ast.Ident)
		if !ok {
			continue
		}

		obj := c.ctx.VarOf(ident)
		if !c.isBytesBuffer(obj.Type()) {
			continue
		}
		candidate := c.find(obj)
		if candidate == nil {
			continue
		}

		paramType := params.At(i).Type()
		if typeis.Named(paramType, "io", "Writer") {
			c.allowedUsages[ident] = struct{}{}
		}
	}
}

func (c *stringsBuilderChecker) isBytesBuffer(typ types.Type) bool {
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}
	return typeis.Named(typ, "bytes", "Buffer")
}

func (c *stringsBuilderChecker) track(v *types.Var, typeNode ast.Expr) {
	c.candidates = append(c.candidates, stringsBuilderCandidate{
		obj:      v,
		typeNode: typeNode,
	})
}

func (c *stringsBuilderChecker) find(v *types.Var) *stringsBuilderCandidate {
	for i := range c.candidates {
		if c.candidates[i].obj == v {
			return &c.candidates[i]
		}
	}
	return nil
}

func (c *stringsBuilderChecker) untrack(v *types.Var, reason string) {
	index := -1
	for i := range c.candidates {
		if c.candidates[i].obj == v {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	c.candidates[index] = c.candidates[len(c.candidates)-1]
	c.candidates = c.candidates[:len(c.candidates)-1]
}
