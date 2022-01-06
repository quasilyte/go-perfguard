package main

import (
	"context"
	"fmt"
	"go/format"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/quasilyte/go-perfguard/internal/imports"
	"github.com/quasilyte/go-perfguard/internal/quickfix"
	"github.com/quasilyte/go-perfguard/perfguard"
	"golang.org/x/tools/go/packages"
)

// runner unifies both `lint` and `optimize` modes.
type runner struct {
	heatmapFile      string
	heatmapThreshold float64
	targets          []string
	autofix          bool

	wd string

	coloredOutput bool
	absFilenames  bool

	loadLintRules bool
	loadOptRules  bool

	goVersion string

	stdout io.Writer
	stderr io.Writer

	pkgFixes []*perfguard.QuickFix
}

func newRunner(stdout, stderr io.Writer) *runner {
	return &runner{stdout: stdout, stderr: stderr}
}

func (r *runner) Run() error {
	if len(r.targets) == 0 {
		return fmt.Errorf("no analysis targets provided")
	}

	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r.wd = wd

	fileSet := token.NewFileSet()
	loadedPackages, err := r.loadPackages(ctx, fileSet, r.targets)
	if err != nil {
		return fmt.Errorf("load packages: %w", err)
	}

	analyzer, err := r.createAnalyzer()
	if err != nil {
		return fmt.Errorf("create analyzer: %w", err)
	}

	target := &perfguard.Target{}
	for _, pkg := range loadedPackages {
		r.pkgFixes = r.pkgFixes[:0]
		target.Files = target.Files[:0]
		for _, f := range pkg.Syntax {
			target.Files = append(target.Files, perfguard.SourceFile{
				Syntax: f,
			})
		}
		target.Fset = fileSet
		target.Sizes = pkg.TypesSizes
		target.Types = pkg.TypesInfo
		target.Pkg = pkg.Types
		if err := analyzer.CheckPackage(target); err != nil {
			return fmt.Errorf("checking %s: %w", pkg.PkgPath, err)
		}
		if len(r.pkgFixes) != 0 {
			if err := r.applyFixes(target); err != nil {
				return fmt.Errorf("apply fixes: %w", err)
			}
		}
	}

	return nil
}

func (r *runner) createAnalyzer() (*perfguard.Analyzer, error) {
	a := perfguard.NewAnalyzer()
	initConfig := &perfguard.Config{
		HeatmapFile:      r.heatmapFile,
		HeatmapThreshold: r.heatmapThreshold,

		GoVersion: r.goVersion,

		Warn: r.handleWarning,

		LoadUniversalRules: true,
		LoadOptRules:       r.loadOptRules,
		LoadLintRules:      r.loadLintRules,
	}
	perfguard.NewAnalyzer()
	if err := a.Init(initConfig); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *runner) handleWarning(w perfguard.Warning) {
	if r.autofix && w.Fix != nil {
		r.pkgFixes = append(r.pkgFixes, w.Fix)
	} else {
		filename := w.Filename
		line := strconv.Itoa(w.Line)
		ruleName := w.Tag
		message := w.Text
		if !r.absFilenames {
			rel, err := filepath.Rel(r.wd, filename)
			if err != nil {
				panic(err)
			}
			filename = rel
		}
		if r.coloredOutput {
			filename = "\033[35m" + filename + "\033[0m"
			line = "\033[32m" + line + "\033[0m"
			ruleName = "\033[31m" + ruleName + "\033[0m"
			message = strings.Replace(message, " => ", " \033[35;1m=>\033[0m ", 1)
		}
		fmt.Fprintf(r.stdout, "%s:%s: %s: %s\n", filename, line, ruleName, message)
	}
}

func (r *runner) applyFixes(target *perfguard.Target) error {
	// TODO: don't run imports fixing for every modified file?
	// We can infer which rules may affect the imports set.

	needFmt := make(map[string]struct{})
	editsPerFile := make(map[string][]quickfix.TextEdit)
	for _, fix := range r.pkgFixes {
		pos := target.Fset.Position(fix.From)
		from := pos.Offset
		filename := pos.Filename
		endPos := target.Fset.Position(fix.To)
		to := endPos.Offset
		if pos.Line != endPos.Line {
			needFmt[filename] = struct{}{}
		}
		editsPerFile[filename] = append(editsPerFile[filename], quickfix.TextEdit{
			StartOffset: from,
			EndOffset:   to,
			Replacement: fix.Replacement,
		})
	}

	// TODO.
	importsConfig := imports.FixConfig{}

	for filename, edits := range editsPerFile {
		fileText, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		afterQuickFixes := quickfix.Apply(fileText, edits)
		newText, err := imports.Fix(importsConfig, afterQuickFixes)
		if err != nil {
			return err
		}
		if _, ok := needFmt[filename]; ok {
			newText, err = format.Source(newText)
			if err != nil {
				return fmt.Errorf("gofmt: %w", err)
			}
		}
		if err := os.WriteFile(filename, newText, 0o600); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) loadPackages(ctx context.Context, fset *token.FileSet, targets []string) ([]*packages.Package, error) {
	loadMode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo |
		packages.NeedTypesSizes
	config := &packages.Config{
		Mode:    loadMode,
		Tests:   false,
		Fset:    fset,
		Context: ctx,
	}
	pkgs, err := packages.Load(config, targets...)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	// We specify tests=false, but just in case we're still going
	// to filter the packages here.
	packageSet := make(map[string]struct{})
	var result []*packages.Package
	for _, pkg := range pkgs {
		if pkg.Name == "" {
			// Empty or invalid package: not interesting.
			continue
		}
		// Skip any test-like package.
		if pkg.Name == "main" && strings.HasSuffix(pkg.PkgPath, ".test") {
			// Implicit main package for tests.
			continue
		}
		if strings.HasSuffix(pkg.Name, "_test") {
			// External test package, like strings_test.
			continue
		}
		if strings.Contains(pkg.ID, ".test]") {
			// Test version of the package.
			continue
		}

		if len(pkg.Errors) != 0 {
			extra := ""
			err := pkg.Errors[0]
			if len(pkg.Errors) > 1 {
				extra = fmt.Sprintf(" (and %d more errors)", len(pkg.Errors)-1)
			}
			fmt.Fprintf(r.stderr, "load %s package: %v%s\n", pkg.Name, err, extra)
		}

		if _, ok := packageSet[pkg.PkgPath]; ok {
			continue
		}
		packageSet[pkg.PkgPath] = struct{}{}

		result = append(result, pkg)
	}

	return result, nil
}
