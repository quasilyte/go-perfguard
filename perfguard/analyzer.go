package perfguard

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/pprof/profile"
	"github.com/quasilyte/go-perfguard/perfguard/rulesdata"
	"github.com/quasilyte/go-ruleguard/ruleguard"
	"github.com/quasilyte/go-ruleguard/ruleguard/ir"
	"github.com/quasilyte/perf-heatmap/heatmap"
)

//go:generate go run ./_rules/precompile/precompile.go -varname Universal -rules ./_rules/universal_rules.go -o ./rulesdata/universal_rules.go
//go:generate go run ./_rules/precompile/precompile.go -varname Opt -rules ./_rules/opt_rules.go -o ./rulesdata/opt_rules.go

type analyzer struct {
	rulesEngine *ruleguard.Engine
	goVersion   ruleguard.GoVersion
	heatmap     *heatmap.Index
	config      *Config
}

func newAnalyzer() *analyzer {
	return &analyzer{}
}

func (a *analyzer) Init(config *Config) error {
	a.config = config
	if err := a.initRulesEngine(); err != nil {
		return err
	}
	if err := a.initHeatmap(config); err != nil {
		return err
	}
	return nil
}

func (a *analyzer) initHeatmap(config *Config) error {
	if config.HeatmapFile == "" {
		return nil
	}
	data, err := os.ReadFile(config.HeatmapFile)
	if err != nil {
		return err
	}
	pprofProfile, err := profile.Parse(bytes.NewReader(data))
	if err != nil {
		return err
	}
	index := heatmap.NewIndex(heatmap.IndexConfig{
		Threshold: config.HeatmapThreshold,
	})
	if err := index.AddProfile(pprofProfile); err != nil {
		return err
	}
	a.heatmap = index
	return nil
}

func (a *analyzer) initRulesEngine() error {
	goVersion, err := ruleguard.ParseGoVersion(a.config.GoVersion)
	if err != nil {
		return fmt.Errorf("parse target Go version: %w", err)
	}
	a.goVersion = goVersion

	rulesEngine := ruleguard.NewEngine()

	fset := token.NewFileSet()
	loadContext := ruleguard.LoadContext{
		Fset: fset,
	}

	toLoad := []struct {
		filename string
		ir       *ir.File
		enabled  bool
	}{
		{"universal_rules.go", rulesdata.Universal, a.config.LoadUniversalRules},
		{"opt_rules.go", rulesdata.Opt, a.config.LoadOptRules},
		// {"lint_rules.go", rulesdata.Lint, a.config.LoadLintRules},
	}
	for _, x := range toLoad {
		if !x.enabled {
			continue
		}
		if err := rulesEngine.LoadFromIR(&loadContext, x.filename, x.ir); err != nil {
			return err
		}
	}

	a.rulesEngine = rulesEngine
	return nil
}

func (a *analyzer) CheckPackage(target *Target) error {
	if err := a.runRules(target); err != nil {
		return err
	}
	return nil
}

func (a *analyzer) getTypeName(typeExpr ast.Expr) string {
	switch typ := typeExpr.(type) {
	case *ast.Ident:
		return typ.Name
	case *ast.StarExpr:
		return a.getTypeName(typ.X)
	case *ast.ParenExpr:
		return a.getTypeName(typ.X)

	default:
		return ""
	}
}

func (a *analyzer) splitFuncName(fn *ast.FuncDecl) (typeName, funcName string) {
	if fn == nil {
		return "", ""
	}
	funcName = fn.Name.Name
	if fn.Recv != nil && len(fn.Recv.List) != 0 {
		typeName = a.getTypeName(fn.Recv.List[0].Type)
	}
	return typeName, funcName
}

func (a *analyzer) minHeatLevel(info *ruleguard.GoRuleInfo) int {
	for _, tag := range info.Group.DocTags {
		switch tag {
		case "o1":
			return 1
		case "o2":
			return 5
		}
	}
	return 0
}

func (a *analyzer) runRules(target *Target) error {
	runContext := ruleguard.RunContext{
		Pkg:       target.Pkg,
		Types:     target.Types,
		Sizes:     target.Sizes,
		Fset:      target.Fset,
		GoVersion: a.goVersion,
	}

	var currentFile *SourceFile

	runContext.Report = func(data *ruleguard.ReportData) {
		startPos := target.Fset.Position(data.Node.Pos())

		if a.heatmap != nil {
			minLevel := a.minHeatLevel(&data.RuleInfo)
			if minLevel != 0 {
				endPos := target.Fset.Position(data.Node.End())
				lineFrom := startPos.Line
				lineTo := endPos.Line
				isHot := false
				typeName, funcName := a.splitFuncName(data.Func)
				key := heatmap.Key{
					TypeName: typeName,
					FuncName: funcName,
					Filename: filepath.Base(startPos.Filename),
					PkgName:  target.Pkg.Name(),
				}
				a.heatmap.QueryLineRange(key, lineFrom, lineTo, func(line int, level heatmap.HeatLevel) bool {
					if level.Global >= minLevel {
						isHot = true
						return false
					}
					return true
				})
				if !isHot {
					return
				}
			}
		}

		message := strings.ReplaceAll(data.Message, "\n", `\n`)
		a.config.Warn(Warning{
			Filename: startPos.Filename,
			Line:     startPos.Line,
			Tag:      data.RuleInfo.Group.Name,
			Text:     message,
		})
	}

	for i := range target.Files {
		currentFile = &target.Files[i]
		if err := a.rulesEngine.Run(&runContext, currentFile.Syntax); err != nil {
			return err
		}
	}

	return nil
}
