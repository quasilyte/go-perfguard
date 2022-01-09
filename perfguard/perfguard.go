package perfguard

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/quasilyte/go-ruleguard/dsl"
)

// This is a temporary kludge to make DSL a direct dependency.
var _ = dsl.Matcher{}

type Analyzer struct {
	impl *analyzer
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		impl: newAnalyzer(),
	}
}

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

type Config struct {
	HeatmapFile      string
	HeatmapThreshold float64

	GoVersion string

	LoadOptRules       bool
	LoadLintRules      bool
	LoadUniversalRules bool

	Warn func(Warning)
}

func (a *Analyzer) Init(config *Config) error {
	return a.impl.Init(config)
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

func (a *Analyzer) CheckPackage(target *Target) error {
	return a.impl.CheckPackage(target)
}
