package funccheckers

import (
	"go/ast"
	"go/token"
)

var zeroLitNode = &ast.BasicLit{
	Kind:  token.INT,
	Value: `0`,
}
