package callcheckers

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/quasilyte/go-perfguard/perfguard/checkers"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

func init() {
	doc := checkers.Doc{
		Name:  "boolValuedMap",
		Score: 3,
	}
	checkers.RegisterFuncChecker(doc, func() checkers.FuncChecker {
		return &boolValuedMapChecker{
			allowedUsages: make(map[*ast.Ident]struct{}),
		}
	})
}

type boolValuedMapChecker struct {
	ctx           *lint.Context
	candidates    []boolValuedMapCandidate
	allowedUsages map[*ast.Ident]struct{}
	nestedFunc    bool
}

type boolValuedMapCandidate struct {
	// A part of a map declaration related to the value.
	// So, in `map[$t]bool` boolValueNode holds the `bool` part.
	boolValueNode ast.Expr

	obj *types.Var

	// Holds `true` expression nodes to be changed to `struct{}{}`.
	assignTrueUsages []ast.Expr

	// Holds `m[k]` expression occurring inside if condition,
	// where we can change it to `_, ok := m[k]; ok`.
	ifCondUsages []ast.Expr

	// Holds `!m[k]` expressions in cond.
	// Like ifCondUsages, but the condition is negated with not operator.
	ifNotCondUsages []*ast.UnaryExpr
}

func (c *boolValuedMapChecker) CheckFunc(ctx *lint.Context, body *ast.BlockStmt) error {
	c.ctx = ctx
	c.candidates = c.candidates[:0]
	c.nestedFunc = false
	for k := range c.allowedUsages {
		delete(c.allowedUsages, k)
	}

	ast.Inspect(body, c.walk)

	for _, candidate := range c.candidates {
		if len(candidate.assignTrueUsages) == 0 {
			continue
		}
		oldNodes := make([]ast.Node, 0, 1+len(candidate.assignTrueUsages))
		newNodes := make([]lint.NodeReplacement, 0, cap(oldNodes))
		oldNodes = append(oldNodes, candidate.boolValueNode)
		newNodes = append(newNodes, lint.NodeReplacement{Text: []byte("struct{}")})
		for _, u := range candidate.assignTrueUsages {
			oldNodes = append(oldNodes, u)
			newNodes = append(newNodes, lint.NodeReplacement{Text: []byte("struct{}{}")})
		}
		for _, u := range candidate.ifCondUsages {
			oldNodes = append(oldNodes, u)
			replacement := make([]byte, 0, 32)
			replacement = append(replacement, "_, ok := "...)
			replacement = append(replacement, ctx.NodeText(u)...)
			replacement = append(replacement, "; ok"...)
			newNodes = append(newNodes, lint.NodeReplacement{Text: replacement})
		}
		for _, u := range candidate.ifNotCondUsages {
			oldNodes = append(oldNodes, u)
			replacement := make([]byte, 0, 32)
			replacement = append(replacement, "_, ok := "...)
			replacement = append(replacement, ctx.NodeText(u.X)...)
			replacement = append(replacement, "; !ok"...)
			newNodes = append(newNodes, lint.NodeReplacement{Text: replacement})
		}
		ctx.MultiChangeSuggest(lint.MultiChangeSuggestParams{
			ReportMessage: "change map[T]bool to map[T]struct{}",
			ReportPos:     candidate.obj.Pos(),
			OldNodes:      oldNodes,
			NewNodes:      newNodes,
		})
	}

	return nil
}

func (c *boolValuedMapChecker) handleMapDefine(n *ast.AssignStmt) bool {
	if c.nestedFunc {
		return false
	}

	if len(n.Lhs) != 1 || len(n.Rhs) != 1 || n.Tok != token.DEFINE {
		return false
	}
	var boolTypeNode *ast.Ident
	var mapTypeNode *ast.MapType
	rhs := n.Rhs[0]
	switch rhs := rhs.(type) {
	case *ast.CallExpr:
		fn, ok := rhs.Fun.(*ast.Ident)
		if !ok || fn.Name != "make" || len(rhs.Args) == 0 {
			return false
		}
		mapTypeNode, _ = rhs.Args[0].(*ast.MapType)

	case *ast.CompositeLit:
		if len(rhs.Elts) != 0 || rhs.Type == nil {
			return false
		}
		mapTypeNode, _ = rhs.Type.(*ast.MapType)
	}
	if mapTypeNode == nil {
		return false
	}

	valueTypeIdent, ok := mapTypeNode.Value.(*ast.Ident)
	if !ok || valueTypeIdent.Name != "bool" {
		return false
	}
	boolTypeNode = valueTypeIdent

	name, ok := n.Lhs[0].(*ast.Ident)
	if !ok {
		return false
	}
	if boolTypeNode == nil {
		return false
	}
	obj := c.ctx.VarOf(name)
	c.track(name, obj, boolTypeNode)
	return true
}

func (c *boolValuedMapChecker) walk(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.FuncLit:
		nestedFunc := c.nestedFunc
		c.nestedFunc = true
		ast.Inspect(n.Body, c.walk)
		c.nestedFunc = nestedFunc
		return false

	case *ast.IfStmt:
		// Check if it's `if m[k] { ... }` usage.
		// Check if it's `if !m[k] { ... }` usage.
		var cond ast.Expr
		var notExpr *ast.UnaryExpr
		unaryExpr, ok := n.Cond.(*ast.UnaryExpr)
		if ok && unaryExpr.Op == token.NOT {
			notExpr = unaryExpr
			cond = notExpr.X
		} else {
			cond = n.Cond
		}

		indexing, ok := cond.(*ast.IndexExpr)
		if !ok {
			return true
		}
		name, ok := indexing.X.(*ast.Ident)
		if !ok {
			return true
		}
		obj := c.ctx.VarOf(name)
		if !c.isBoolValuedMap(obj.Type()) {
			return true
		}
		candidate := c.find(obj)
		if candidate != nil {
			c.allowedUsages[name] = struct{}{}
			if notExpr != nil {
				candidate.ifNotCondUsages = append(candidate.ifNotCondUsages, notExpr)
			} else {
				candidate.ifCondUsages = append(candidate.ifCondUsages, cond)
			}
		}
		return true

	case *ast.RangeStmt:
		// Permit ranging over a tracked map keys.
		if n.Value != nil {
			return true
		}
		name, ok := n.X.(*ast.Ident)
		if !ok {
			return true
		}
		obj := c.ctx.VarOf(name)
		if c.isBoolValuedMap(obj.Type()) {
			c.allowedUsages[name] = struct{}{}
		}

	case *ast.AssignStmt:
		if c.handleMapDefine(n) {
			return false
		}
		// Recognize `m[k] = true` usages.
		if len(n.Lhs) != 1 || len(n.Rhs) != 1 || n.Tok != token.ASSIGN {
			return true
		}
		trueNode, ok := n.Rhs[0].(*ast.Ident)
		if !ok || trueNode.Name != "true" {
			return true
		}
		indexing, ok := n.Lhs[0].(*ast.IndexExpr)
		if !ok {
			return true
		}
		name, ok := indexing.X.(*ast.Ident)
		if !ok {
			return true
		}
		obj := c.ctx.VarOf(name)
		if !c.isBoolValuedMap(obj.Type()) {
			return true
		}
		candidate := c.find(obj)
		if candidate != nil {
			c.allowedUsages[name] = struct{}{}
			candidate.assignTrueUsages = append(candidate.assignTrueUsages, trueNode)
		}

	case *ast.CallExpr:
		// Permit `len(m)` expressions.
		if fn, ok := n.Fun.(*ast.Ident); ok && fn.Name == "len" && len(n.Args) == 1 {
			arg, ok := n.Args[0].(*ast.Ident)
			if !ok {
				return true
			}
			c.allowedUsages[arg] = struct{}{}
		}

	case *ast.Ident:
		// If it's a buffer usage, untrack the candidate.
		// We stop tree traversal when we recognize the usage.
		obj := c.ctx.VarOf(n)
		if c.isBoolValuedMap(obj.Type()) {
			if _, ok := c.allowedUsages[n]; !ok {
				c.untrack(obj, "uncategorized usage")
			}
		}
	}

	return true
}

func (c *boolValuedMapChecker) isBoolValuedMap(typ types.Type) bool {
	m, ok := typ.(*types.Map)
	if !ok {
		return false
	}
	valueType, ok := m.Elem().(*types.Basic)
	return ok && valueType.Kind() == types.Bool
}

func (c *boolValuedMapChecker) track(id *ast.Ident, v *types.Var, boolValueNode ast.Expr) {
	c.candidates = append(c.candidates, boolValuedMapCandidate{
		obj:           v,
		boolValueNode: boolValueNode,
	})
}

func (c *boolValuedMapChecker) find(v *types.Var) *boolValuedMapCandidate {
	for i := range c.candidates {
		if c.candidates[i].obj == v {
			return &c.candidates[i]
		}
	}
	return nil
}

func (c *boolValuedMapChecker) untrack(v *types.Var, reason string) {
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
