package funccheckers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/go-toolsmith/astequal"
	"github.com/go-toolsmith/typep"
	"github.com/quasilyte/go-perfguard/internal/typeis"
	"github.com/quasilyte/go-perfguard/perfguard/checkers"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
)

func init() {
	doc := checkers.Doc{
		Name:  "prealloc",
		Score: 4,
	}
	checkers.RegisterFuncChecker(doc, func() checkers.FuncChecker {
		return &preallocChecker{allowedUsages: make(map[*ast.Ident]struct{})}
	})
}

type preallocChecker struct {
	ctx *lint.Context

	allowedUsages map[*ast.Ident]struct{}

	nestedFunc bool

	depth int

	candidates []preallocCandidate
}

type preallocCandidate struct {
	obj *types.Var

	declNode ast.Node
	typeExpr ast.Expr

	depth int

	bound ast.Expr

	hotNode ast.Node
}

func (c *preallocChecker) CheckFunc(ctx *lint.Context, body *ast.BlockStmt) error {
	c.ctx = ctx
	c.candidates = c.candidates[:0]
	c.nestedFunc = false
	c.depth = 0
	for k := range c.allowedUsages {
		delete(c.allowedUsages, k)
	}

	ast.Inspect(body, c.walk)

	for _, candidate := range c.candidates {
		if candidate.bound == nil {
			continue
		}
		declValueSpec, _ := candidate.declNode.(*ast.ValueSpec)
		switch {
		case declValueSpec == nil && typeis.Map(candidate.obj.Type()):
			// Fixing from `x := make(map[K]V)` to `x := make(map[K]V, len(bound))`.
			makeCall := &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					candidate.typeExpr,
					&ast.CallExpr{
						Fun:  &ast.Ident{Name: "len"},
						Args: []ast.Expr{candidate.bound},
					},
				},
			}
			ctx.SuggestNode(lint.SuggestParams{
				OldNode:  candidate.declNode,
				NewNode:  makeCall,
				HotNodes: []ast.Node{candidate.hotNode},
			})

		case declValueSpec == nil && typeis.Slice(candidate.obj.Type()):
			// Fixing from `x := []T{}` to `x := make([]T, 0, len(bound))`.
			makeCall := &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					candidate.typeExpr,
					zeroLitNode,
					&ast.CallExpr{
						Fun:  &ast.Ident{Name: "len"},
						Args: []ast.Expr{candidate.bound},
					},
				},
			}
			ctx.SuggestNode(lint.SuggestParams{
				OldNode:  candidate.declNode,
				NewNode:  makeCall,
				HotNodes: []ast.Node{candidate.hotNode},
			})

		case declValueSpec != nil && typeis.Map(candidate.obj.Type()):
			// Fixing from `var x = make(map[K]V)` to `var x = make(map[K]V, len(bound))`.
			makeCall := &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					candidate.typeExpr,
					&ast.CallExpr{
						Fun:  &ast.Ident{Name: "len"},
						Args: []ast.Expr{candidate.bound},
					},
				},
			}
			fixedValueSpec := &ast.ValueSpec{
				Type:   declValueSpec.Type,
				Names:  declValueSpec.Names,
				Values: []ast.Expr{makeCall},
			}
			ctx.SuggestNode(lint.SuggestParams{
				OldNode:  declValueSpec,
				NewNode:  fixedValueSpec,
				HotNodes: []ast.Node{candidate.hotNode},
			})

		case declValueSpec != nil && typeis.Slice(candidate.obj.Type()):
			if declValueSpec.Values == nil {
				// Report that `var x []T` can be changed to `var x = make([]T, 0, len(bound))`,
				// but don'e apply an auto fix here.
				expr := c.ctx.NodeText(candidate.bound)
				ctx.Report(lint.ReportParams{
					PosNode:  declValueSpec,
					Message:  fmt.Sprintf("can use len(%s) as make size hint for %s", expr, candidate.obj.Name()),
					HotNodes: []ast.Node{candidate.hotNode},
				})
			} else {
				makeCall := &ast.CallExpr{
					Fun: &ast.Ident{Name: "make"},
					Args: []ast.Expr{
						candidate.typeExpr,
						zeroLitNode,
						&ast.CallExpr{
							Fun:  &ast.Ident{Name: "len"},
							Args: []ast.Expr{candidate.bound},
						},
					},
				}
				fixedValueSpec := &ast.ValueSpec{
					Type:   declValueSpec.Type,
					Names:  declValueSpec.Names,
					Values: []ast.Expr{makeCall},
				}
				ctx.SuggestNode(lint.SuggestParams{
					OldNode:  declValueSpec,
					NewNode:  fixedValueSpec,
					HotNodes: []ast.Node{candidate.hotNode},
				})
			}
		}

	}

	return nil
}

func (c *preallocChecker) visitDecl(decl ast.Node, declTypeExpr ast.Expr, id *ast.Ident, init ast.Expr) bool {
	obj := c.ctx.VarOf(id)

	switch {
	case init == nil:
		if typeis.Slice(obj.Type()) {
			// Track `var $x []T`.
			candidate := c.track(obj)
			candidate.declNode = decl
			candidate.typeExpr = declTypeExpr
			return false
		}

	case init != nil:
		var typeExpr ast.Expr
		switch e := init.(type) {
		case *ast.CallExpr:
			called, ok := e.Fun.(*ast.Ident)
			if !ok || called.Name != "make" {
				return true
			}
			switch obj.Type().(type) {
			case *types.Slice:
				if len(e.Args) != 2 {
					return true
				}
				lengthArg, ok := e.Args[1].(*ast.BasicLit)
				if !ok || lengthArg.Kind != token.INT || lengthArg.Value != `0` {
					return true
				}
			case *types.Map:
				if len(e.Args) != 1 {
					return true
				}
			}
			typeExpr = e.Args[0]

		case *ast.CompositeLit:
			if len(e.Elts) != 0 {
				return true
			}
			typeExpr = e.Type
		}

		if typeExpr == nil {
			return true
		}

		switch obj.Type().(type) {
		case *types.Map, *types.Slice:
			candidate := c.track(obj)
			candidate.declNode = decl
			candidate.typeExpr = typeExpr
			return false
		}
	}

	return true
}

func (c *preallocChecker) walk(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.FuncLit:
		nestedFunc := c.nestedFunc
		c.nestedFunc = true
		ast.Inspect(n.Body, c.walk)
		c.nestedFunc = nestedFunc
		return false

	case *ast.BlockStmt:
		c.depth++
		for _, stmt := range n.List {
			ast.Inspect(stmt, c.walk)
		}
		c.depth--
		return false

	case *ast.RangeStmt:
		// Even if we have candidates, maybe their depths are incompatible.
		// This code also handles the 0 candidates case.
		skip := true
		for i := range c.candidates {
			if c.candidates[i].depth == c.depth {
				skip = false // Have at least 1 potential candidate
				break
			}
		}
		if skip {
			return true
		}
		rangeKey, _ := n.Key.(*ast.Ident)
		var rangeExpr = n.X
		if !typep.SideEffectFree(c.ctx.Target.Types, rangeExpr) {
			return true
		}
		// TODO: handle more complicated for bodies?
		for _, stmt := range n.Body.List {
			if !c.isSimpleStmt(stmt) {
				return true
			}
		}
		for _, stmt := range n.Body.List {
			assign, ok := stmt.(*ast.AssignStmt)
			if !ok || assign.Tok != token.ASSIGN || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
				continue
			}
			// For slice candidates, we're interested in append calls.
			// For map candidates, we're interested in IndexExpr assignments.
			lhs := assign.Lhs[0]
			rhs := assign.Rhs[0]
			switch lhs := lhs.(type) {
			case *ast.IndexExpr:
				if rangeKey == nil || rangeKey.Name == "_" {
					continue
				}
				if !astequal.Expr(lhs.Index, rangeKey) {
					continue
				}
				indexed, ok := lhs.X.(*ast.Ident)
				if !ok {
					continue
				}
				obj := c.ctx.VarOf(indexed)
				if !typeis.Map(obj.Type()) {
					continue
				}
				candidate := c.find(obj)
				if candidate != nil && candidate.bound == nil {
					candidate.bound = rangeExpr
					candidate.hotNode = assign
					c.allowedUsages[indexed] = struct{}{}
				}

			case *ast.Ident:
				call, ok := rhs.(*ast.CallExpr)
				if !ok {
					continue
				}
				called, ok := call.Fun.(*ast.Ident)
				if !ok || called.Name != "append" {
					continue
				}
				if call.Ellipsis.IsValid() || len(call.Args) != 2 {
					continue
				}
				if !astequal.Expr(lhs, call.Args[0]) {
					continue
				}
				obj := c.ctx.VarOf(lhs)
				candidate := c.find(obj)
				if candidate != nil && candidate.bound == nil {
					candidate.bound = rangeExpr
					candidate.hotNode = assign
					c.allowedUsages[lhs] = struct{}{}
					c.allowedUsages[call.Args[0].(*ast.Ident)] = struct{}{}
				}
			}
		}

	case *ast.AssignStmt:
		if c.nestedFunc {
			return true
		}
		if len(n.Lhs) != 1 || len(n.Rhs) != 1 || n.Tok != token.DEFINE {
			return true
		}
		lhs, ok := n.Lhs[0].(*ast.Ident)
		if !ok {
			return true
		}
		return c.visitDecl(n.Rhs[0], nil, lhs, n.Rhs[0])

	case *ast.ValueSpec:
		if c.nestedFunc {
			return true
		}
		if len(n.Names) != 1 {
			return true
		}
		var init ast.Expr
		switch len(n.Values) {
		case 0:
			// init remains nil.
		case 1:
			init = n.Values[0]
		default:
			return true
		}
		return c.visitDecl(n, n.Type, n.Names[0], init)

	case *ast.Ident:
		if len(c.candidates) == 0 {
			return true
		}
		obj := c.ctx.VarOf(n)
		switch obj.Type().(type) {
		case *types.Slice, *types.Map:
			if _, ok := c.allowedUsages[n]; !ok {
				candidate, index := c.findWithIndex(obj)
				if candidate != nil && candidate.bound == nil {
					c.untrackByIndex(index, "uncategorized usage")
				}
			}
		}

	}

	return true
}

func (c *preallocChecker) isSimpleStmt(n ast.Stmt) bool {
	switch n := n.(type) {
	case *ast.AssignStmt, *ast.ExprStmt, *ast.DeclStmt, *ast.IncDecStmt:
		return true

	case *ast.IfStmt:
		if len(n.Body.List) > 4 {
			return false
		}
		for _, stmt := range n.Body.List {
			if !c.isSimpleStmt(stmt) {
				return false
			}
		}
		return true

	default:
		return false
	}
}

func (c *preallocChecker) track(v *types.Var) *preallocCandidate {
	c.candidates = append(c.candidates, preallocCandidate{
		obj:   v,
		depth: c.depth,
	})
	return &c.candidates[len(c.candidates)-1]
}

func (c *preallocChecker) findWithIndex(v *types.Var) (candidate *preallocCandidate, index int) {
	for i := range c.candidates {
		if c.candidates[i].obj == v {
			return &c.candidates[i], i
		}
	}
	return nil, -1
}

func (c *preallocChecker) find(v *types.Var) *preallocCandidate {
	candidate, _ := c.findWithIndex(v)
	return candidate
}

func (c *preallocChecker) untrackByIndex(index int, reason string) {
	c.candidates[index] = c.candidates[len(c.candidates)-1]
	c.candidates = c.candidates[:len(c.candidates)-1]
}
