package main

import (
	"context"
	"fmt"
	"go/token"
	"io"
	"strings"

	"github.com/quasilyte/go-perfguard/perfguard"
	"golang.org/x/tools/go/packages"
)

// runner unifies both `lint` and `optimize` modes.
type runner struct {
	heatmapFile      string
	heatmapThreshold float64
	targets          []string
	autofix          bool

	loadLintRules bool
	loadOptRules  bool

	stdout io.Writer
	stderr io.Writer
}

func newRunner(stdout, stderr io.Writer) *runner {
	return &runner{stdout: stdout, stderr: stderr}
}

func (r *runner) Run() error {
	if len(r.targets) == 0 {
		return fmt.Errorf("no analysis targets provided")
	}

	ctx := context.Background()

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
	}

	return nil
}

func (r *runner) createAnalyzer() (*perfguard.Analyzer, error) {
	a := perfguard.NewAnalyzer()
	initConfig := &perfguard.Config{
		HeatmapFile:      r.heatmapFile,
		HeatmapThreshold: r.heatmapThreshold,

		Autofix: r.autofix,

		Warn: r.reportWarning,

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

func (r *runner) reportWarning(w perfguard.Warning) {
	fmt.Fprintf(r.stdout, "%s:%d: %s: %s\n", w.Filename, w.Line, w.Tag, w.Text)
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

		for _, err := range pkg.Errors {
			fmt.Fprintf(r.stderr, "load %s package: %v\n", pkg.Name, err)
		}

		if _, ok := packageSet[pkg.PkgPath]; ok {
			continue
		}
		packageSet[pkg.PkgPath] = struct{}{}

		result = append(result, pkg)
	}

	return result, nil
}
