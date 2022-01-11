package perfguard

import (
	"github.com/quasilyte/go-perfguard/perfguard/lint"
	"github.com/quasilyte/go-ruleguard/dsl"
	"github.com/quasilyte/perf-heatmap/heatmap"
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

type Config struct {
	Heatmap *heatmap.Index

	GoVersion string

	LoadOptRules       bool
	LoadLintRules      bool
	LoadUniversalRules bool

	Warn func(lint.Warning)
}

func (a *Analyzer) Init(config *Config) error {
	return a.impl.Init(config)
}

func (a *Analyzer) CheckPackage(target *lint.Target) error {
	return a.impl.CheckPackage(target)
}
