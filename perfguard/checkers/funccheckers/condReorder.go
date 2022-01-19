package callcheckers

import (
	"go/ast"
	"go/token"

	"github.com/go-toolsmith/astcopy"
	"github.com/quasilyte/go-perfguard/perfguard/checkers"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

func init() {
	doc := checkers.Doc{
		Name:  "condReorder",
		Score: 2,
	}
	checkers.RegisterFuncChecker(doc, func() checkers.FuncChecker {
		return &condReorderChecker{}
	})
}

type condReorderChecker struct {
	ctx *lint.Context

	order []int
}

type condReorderExprInfo struct {
	hasCall bool
}

func (c *condReorderChecker) CheckFunc(ctx *lint.Context, body *ast.BlockStmt) error {
	c.ctx = ctx
	c.order = c.order[:0]

	ast.Inspect(body, c.walk)

	return nil
}

func (c *condReorderChecker) walk(n ast.Node) bool {
	// We only analyze the most common places for conditions that may matter.

	switch n := n.(type) {
	case *ast.IfStmt:
		c.walkCond(n.Cond)

	case *ast.ForStmt:
		if n.Cond != nil {
			c.walkCond(n.Cond)
		}

	case *ast.ReturnStmt:
		for _, x := range n.Results {
			c.walkCond(x)
		}

	case *ast.AssignStmt:
		for _, x := range n.Rhs {
			c.walkCond(x)
		}
	}

	return true
}

func (c *condReorderChecker) walkCond(cond ast.Expr) bool {
	e, ok := cond.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	if e.Op != token.LAND && e.Op != token.LOR {
		return false
	}

	stop := false
	if lhs, ok := e.X.(*ast.BinaryExpr); ok && (lhs.Op == token.LAND || lhs.Op == token.LOR) {
		stop = true
		if c.walkCond(lhs) {
			return true
		}
	}
	if rhs, ok := e.X.(*ast.BinaryExpr); ok && (rhs.Op == token.LAND || rhs.Op == token.LOR) {
		stop = true
		if c.walkCond(rhs) {
			return true
		}
	}
	if stop {
		return false
	}

	var lhsInfo condReorderExprInfo
	var rhsInfo condReorderExprInfo
	if !c.exprCost(&lhsInfo, e.X) || !c.exprCost(&rhsInfo, e.Y) {
		return false
	}

	// expensive() && simple => simple && expensive()
	// expensive() || simple => simple || expensive()
	if lhsInfo.hasCall && !rhsInfo.hasCall {
		// Now check that lhs doesn't depend on rhs (and vice versa).
		// This is a rare branch, so we can afford a slow algorithm here.
		if c.independent(e.X, e.Y) {
			reordered := astcopy.BinaryExpr(e)
			reordered.X, reordered.Y = reordered.Y, reordered.X
			c.ctx.SuggestNode(lint.SuggestParams{
				OldNode: e,
				NewNode: reordered,
			})
			return true
		}
	}

	return false
}

func (c *condReorderChecker) independent(x, y ast.Expr) bool {
	var xVars map[string]struct{}
	ast.Inspect(x, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if xVars == nil {
				xVars = make(map[string]struct{})
			}
			xVars[id.Name] = struct{}{}
		}
		return true
	})
	found := false
	ast.Inspect(y, func(n ast.Node) bool {
		if found {
			return false
		}
		if id, ok := n.(*ast.Ident); ok {
			if _, ok := xVars[id.Name]; ok {
				found = true
			}
		}
		return true
	})
	return !found
}

func (c *condReorderChecker) exprCost(info *condReorderExprInfo, e ast.Expr) bool {
	switch e := e.(type) {
	case *ast.Ident, *ast.BasicLit:
		// Trivial, can analyze.
		return true

	case *ast.SelectorExpr:
		return c.exprCost(info, e.X)

	case *ast.StarExpr:
		return c.exprCost(info, e.X)

	case *ast.CallExpr:
		called, ok := e.Fun.(*ast.Ident)
		if ok {
			switch called.Name {
			case "len", "cap":
				if len(e.Args) != 1 {
					return false
				}
				return c.exprCost(info, e.Args[0])
			}
		}
		info.hasCall = true
		return true

	case *ast.BinaryExpr:
		switch e.Op {
		case token.ADD, token.SUB:
			// OK.
		default:
			return false
		}
		return c.exprCost(info, e.X) && c.exprCost(info, e.Y)

	case *ast.ParenExpr:
		return c.exprCost(info, e.X)

	case *ast.UnaryExpr:
		switch e.Op {
		case token.ADD, token.SUB, token.NOT:
			return c.exprCost(info, e.X)
		}

	case *ast.IndexExpr:
		return c.exprCost(info, e.X) && c.exprCost(info, e.Index)
	}

	return false
}
