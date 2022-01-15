package perfguard

import (
	"fmt"
	"go/token"
	"path/filepath"
	"strings"
	"time"

	"github.com/quasilyte/go-ruleguard/ruleguard"
	"github.com/quasilyte/go-ruleguard/ruleguard/ir"
	"github.com/quasilyte/perf-heatmap/heatmap"

	"github.com/quasilyte/go-perfguard/internal/resolve"
	"github.com/quasilyte/go-perfguard/perfguard/lint"
	"github.com/quasilyte/go-perfguard/perfguard/rulesdata"
)

//go:generate go run ./_rules/precompile/precompile.go -varname Universal -rules ./_rules/universal_rules.go -o ./rulesdata/universal_rules.go
//go:generate go run ./_rules/precompile/precompile.go -varname Opt -rules ./_rules/opt_rules.go -o ./rulesdata/opt_rules.go

type analyzer struct {
	rulesEngine *ruleguard.Engine

	checkers []*targetChecker

	goVersion ruleguard.GoVersion
	config    *Config
}

func newAnalyzer() *analyzer {
	return &analyzer{}
}

func (a *analyzer) Init(config *Config) error {
	a.config = config
	a.checkers = createCheckers(config)
	return a.initRulesEngine()
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

func (a *analyzer) CheckPackage(target *lint.Target) error {
	if err := a.runRules(target); err != nil {
		return err
	}
	for _, c := range a.checkers {
		if err := c.CheckTarget(target); err != nil {
			return err
		}
	}
	return nil
}

func (a *analyzer) hasReformatTag(info *ruleguard.GoRuleInfo) bool {
	for _, tag := range info.Group.DocTags {
		if tag == "reformat" {
			return true
		}
	}
	return false
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

func (a *analyzer) runRules(target *lint.Target) error {
	ruleguardContext := ruleguard.RunContext{
		Pkg:         target.Pkg,
		Types:       target.Types,
		Sizes:       target.Sizes,
		Fset:        target.Fset,
		GoVersion:   a.goVersion,
		TruncateLen: 100,
	}

	var currentFile *lint.SourceFile

	ruleguardContext.Report = func(data *ruleguard.ReportData) {
		startPos := target.Fset.Position(data.Node.Pos())

		samplesTime := time.Duration(0)
		if a.config.Heatmap != nil {
			minLevel := a.minHeatLevel(&data.RuleInfo)
			if minLevel != 0 {
				endPos := target.Fset.Position(data.Node.End())
				lineFrom := startPos.Line
				lineTo := endPos.Line
				isHot := false
				typeName, funcName := resolve.SplitFuncName(data.Func)
				key := heatmap.Key{
					TypeName: typeName,
					FuncName: funcName,
					Filename: filepath.Base(startPos.Filename),
					PkgName:  target.Pkg.Name(),
				}
				totalValue := int64(0)
				a.config.Heatmap.QueryLineRange(key, lineFrom, lineTo, func(l heatmap.LineStats) bool {
					if l.GlobalHeatLevel >= minLevel {
						isHot = true
					}
					totalValue += l.Value
					return true
				})
				if !isHot {
					return
				}
				samplesTime = time.Duration(totalValue)
			}
		}

		var fix *lint.QuickFix
		if data.Suggestion != nil {
			s := data.Suggestion
			fix = &lint.QuickFix{
				From:        s.From,
				To:          s.To,
				Replacement: make([]byte, len(s.Replacement)),
				Reformat:    a.hasReformatTag(&data.RuleInfo),
			}
			copy(fix.Replacement, s.Replacement)
		}

		message := strings.ReplaceAll(data.Message, "\n", `\n`)
		a.config.Warn(lint.Warning{
			Filename:    startPos.Filename,
			Line:        startPos.Line,
			Tag:         data.RuleInfo.Group.Name,
			Text:        message,
			Fix:         fix,
			SamplesTime: samplesTime,
		})
	}

	for i := range target.Files {
		currentFile = &target.Files[i]
		if err := a.rulesEngine.Run(&ruleguardContext, currentFile.Syntax); err != nil {
			return err
		}
	}

	return nil
}
