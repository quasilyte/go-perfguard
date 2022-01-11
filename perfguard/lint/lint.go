package lint

import (
	"go/ast"
	"go/token"
	"go/types"
)

type Warning struct {
	Filename string
	Line     int
	Tag      string
	Text     string
	Fix      *QuickFix
}

type QuickFix struct {
	From        token.Pos
	To          token.Pos
	Replacement []byte
}

type SourceFile struct {
	Syntax *ast.File
}

type Target struct {
	Pkg   *types.Package
	Fset  *token.FileSet
	Types *types.Info
	Sizes types.Sizes
	Files []SourceFile
}
