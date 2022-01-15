package lint

import (
	"go/ast"
	"go/token"
	"go/types"
	"time"
)

type Warning struct {
	Filename string
	Line     int
	Tag      string
	Text     string
	Fix      *QuickFix

	SamplesTime time.Duration
}

type QuickFix struct {
	From        token.Pos
	To          token.Pos
	Replacement []byte
	Reformat    bool
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
