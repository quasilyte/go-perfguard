package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/pprof/profile"
	"github.com/quasilyte/go-perfguard/internal/imports"
	"github.com/quasilyte/go-perfguard/internal/quickfix"
	"github.com/quasilyte/go-perfguard/perfguard"
	"github.com/quasilyte/perf-heatmap/heatmap"
	"golang.org/x/tools/go/packages"
)

type arguments struct {
	heatmapFile      string
	heatmapThreshold float64
}

type statistics struct {
	pkgloadTime  int64
	analysisTime int64
}

// runner unifies both `lint` and `optimize` modes.
type runner struct {
	targets []string
	autofix bool

	debugEnabled bool

	args  arguments
	stats statistics

	heatmap          *heatmap.Index
	heatmapPackages  map[string]struct{}
	heatmapFiles     map[string]struct{}
	numFilesSkipped  int
	numFilesAnalyzed int

	wd string

	coloredOutput bool
	absFilenames  bool

	loadLintRules bool
	loadOptRules  bool

	goVersion string

	stdout io.Writer
	stderr io.Writer

	pkgWarnings []perfguard.Warning

	// We try to avoid reporting more errors than necessary.
	// There is a hard limit on how many errors we'll print.
	// There is also a filter that will exclude any repeated
	// errors from the output (errorSet).
	errorSet    map[string]struct{}
	errorsList  []string
	extraErrors int

	numLoadCalls int
}

func newRunner(stdout, stderr io.Writer) *runner {
	debugEnabled := os.Getenv("PERFGUARD_DEBUG") == "1"
	return &runner{
		stdout:       stdout,
		stderr:       stderr,
		debugEnabled: debugEnabled,

		errorSet: make(map[string]struct{}),
	}
}

func (r *runner) pushErrorf(key, formatString string, args ...interface{}) {
	if _, ok := r.errorSet[key]; ok {
		return
	}
	const maxErrorsNum = 10
	if len(r.errorSet) > maxErrorsNum {
		r.extraErrors++
		return
	}
	r.errorSet[key] = struct{}{}

	msg := fmt.Sprintf(formatString, args...)
	r.errorsList = append(r.errorsList, msg)
}

func (r *runner) printAllErrors() {
	for _, msg := range r.errorsList {
		tag := ">> error"
		if r.coloredOutput {
			tag = "\033[31;1m" + tag + "\033[0m"
		}
		msg := tag + ": " + msg + "\n"
		_, err := io.WriteString(r.stderr, msg)
		if err != nil {
			panic(err)
		}
	}
}

func (r *runner) printDebugf(formatString string, args ...interface{}) {
	if !r.debugEnabled {
		return
	}
	tag := ">> debug"
	if r.coloredOutput {
		tag = "\033[34;1m" + tag + "\033[0m"
	}
	msg := tag + ": " + fmt.Sprintf(formatString, args...) + "\n"
	_, err := io.WriteString(r.stderr, msg)
	if err != nil {
		panic(err)
	}
}

func (r *runner) Run() error {
	if len(r.targets) == 0 {
		return fmt.Errorf("no analysis targets provided")
	}

	ctx := context.Background()
	startTime := time.Now()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r.wd = wd

	if r.args.heatmapFile != "" {
		heatmapIndex, err := r.createHeatmap()
		if err != nil {
			return err
		}
		r.heatmap = heatmapIndex
		r.inspectHeatmap()
	}

	fileSet := token.NewFileSet()
	targetPackages, err := r.findPackages(ctx, fileSet, r.targets)
	if err != nil {
		return fmt.Errorf("load packages: %w", err)
	}

	analyzer, err := r.createAnalyzer()
	if err != nil {
		return fmt.Errorf("create analyzer: %w", err)
	}

	if r.heatmap != nil {
		filtered := targetPackages[:0]
		numSkipped := 0
		for _, ref := range targetPackages {
			if _, ok := r.heatmapPackages[ref.name]; ok {
				filtered = append(filtered, ref)
			} else {
				r.printDebugf("skip %s package", ref.path)
				numSkipped++
			}
		}
		if numSkipped != 0 {
			r.printDebugf("skipped %d packages", numSkipped)
		}
		targetPackages = filtered
	}

	if r.debugEnabled {
		for _, ref := range targetPackages {
			r.printDebugf("found %s package", ref.path)
		}
	}

	target := &perfguard.Target{}
	for i, ref := range targetPackages {
		r.printDebugf("loading %s package (%d/%d)", ref.path, i+1, len(targetPackages))
		pkg, err := r.loadPackage(ctx, fileSet, ref)
		if err != nil {
			return err
		}
		r.printDebugf("analyzing %s package (%d/%d)", ref.path, i+1, len(targetPackages))

		r.pkgWarnings = r.pkgWarnings[:0]
		target.Files = target.Files[:0]
		for _, f := range pkg.Syntax {
			if r.heatmapFiles != nil {
				filename := fileSet.Position(f.Pos()).Filename
				if _, ok := r.heatmapFiles[filepath.Base(filename)]; !ok {
					r.numFilesSkipped++
					continue
				}
			}
			r.numFilesAnalyzed++
			target.Files = append(target.Files, perfguard.SourceFile{
				Syntax: f,
			})
		}
		target.Fset = fileSet
		target.Sizes = pkg.TypesSizes
		target.Types = pkg.TypesInfo
		target.Pkg = pkg.Types
		start := time.Now()
		err = analyzer.CheckPackage(target)
		elapsed := time.Since(start)
		atomic.AddInt64(&r.stats.pkgloadTime, int64(elapsed))
		if err != nil {
			return fmt.Errorf("checking %s: %w", pkg.PkgPath, err)
		}
		if len(r.pkgWarnings) != 0 {
			if err := r.handleWarnings(target); err != nil {
				return err
			}
		}
	}

	timeElapsed := time.Since(startTime)

	if r.numFilesSkipped == 0 {
		r.printDebugf("analyzed %d files", r.numFilesAnalyzed)
	} else {
		r.printDebugf("analyzed %d files (%d skipped)", r.numFilesAnalyzed, r.numFilesSkipped)
	}
	r.printDebugf("packages.Load calls: %d", r.numLoadCalls)
	r.printDebugf("packages.Load time: %.2fs", time.Duration(r.stats.pkgloadTime).Seconds())
	r.printDebugf("analysis time: %.2fs", time.Duration(r.stats.analysisTime).Seconds())
	r.printDebugf("total time: %.2fs", timeElapsed.Seconds())

	if len(r.errorsList) != 0 {
		r.printAllErrors()
		if r.extraErrors != 0 {
			fmt.Fprintf(r.stderr, "+ %d more errors\n", r.extraErrors)
		}
	}

	return nil
}

func (r *runner) createAnalyzer() (*perfguard.Analyzer, error) {
	a := perfguard.NewAnalyzer()
	initConfig := &perfguard.Config{
		Heatmap: r.heatmap,

		GoVersion: r.goVersion,

		Warn: r.appendWarning,

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

func (r *runner) appendWarning(w perfguard.Warning) {
	r.pkgWarnings = append(r.pkgWarnings, w)
}

func (r *runner) reportWarning(w *perfguard.Warning) {
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
		ruleName = "\033[93m" + ruleName + "\033[0m"
		message = strings.Replace(message, " => ", " \033[35;1m=>\033[0m ", 1)
	}
	fmt.Fprintf(r.stdout, "%s:%s: %s: %s\n", filename, line, ruleName, message)
}

func (r *runner) handleWarnings(target *perfguard.Target) error {
	// TODO: don't run imports fixing for every modified file?
	// We can infer which rules may affect the imports set.

	type warningWithFix struct {
		w   *perfguard.Warning
		fix quickfix.TextEdit
	}

	needFmt := make(map[string]struct{})
	fixablePerFile := make(map[string][]warningWithFix)
	for i := range r.pkgWarnings {
		w := &r.pkgWarnings[i]
		if !r.autofix || w.Fix == nil {
			r.reportWarning(w)
			continue
		}
		pos := target.Fset.Position(w.Fix.From)
		from := pos.Offset
		filename := pos.Filename
		endPos := target.Fset.Position(w.Fix.To)
		to := endPos.Offset
		if pos.Line != endPos.Line {
			needFmt[filename] = struct{}{}
		}
		fix := quickfix.TextEdit{
			StartOffset: from,
			EndOffset:   to,
			Replacement: w.Fix.Replacement,
		}
		fixablePerFile[filename] = append(fixablePerFile[filename], warningWithFix{w: w, fix: fix})
	}

	// TODO.
	importsConfig := imports.FixConfig{}

	for filename, pairs := range fixablePerFile {
		quickfix.Sort(pairs, func(i int) quickfix.TextEdit {
			return pairs[i].fix
		})
		edits := make([]quickfix.TextEdit, len(pairs))
		for i, p := range pairs {
			edits[i] = p.fix
		}
		fileText, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		afterQuickFixes, overlapping := quickfix.Apply(fileText, edits)
		for _, pairIndex := range overlapping {
			r.reportWarning(pairs[pairIndex].w)
		}
		newText, err := imports.Fix(importsConfig, afterQuickFixes)
		if err != nil {
			return fmt.Errorf("fix imports: %w", err)
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

func (r *runner) loadPackage(ctx context.Context, fset *token.FileSet, ref packageRef) (*packages.Package, error) {
	r.numLoadCalls++

	loadMode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedSyntax |
		packages.NeedTypes |
		packages.NeedTypesInfo |
		packages.NeedTypesSizes
	config := &packages.Config{
		Mode:    loadMode,
		Tests:   false,
		Fset:    fset,
		Context: ctx,
	}
	start := time.Now()
	loaded, err := packages.Load(config, ref.path)
	elapsed := time.Since(start)
	atomic.AddInt64(&r.stats.pkgloadTime, int64(elapsed))
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}
	if len(loaded) != 1 {
		return nil, fmt.Errorf("expected 1 package for %s, got %d", ref.path, len(loaded))
	}
	pkg := loaded[0]

	if len(pkg.Errors) != 0 {
		err := pkg.Errors[0]
		r.pushErrorf(err.Msg, "load %s package: %v", pkg.Name, err)
	}

	return pkg, nil
}

type packageRef struct {
	id   string
	name string
	path string
}

// findPackages returns a list of matched packages for given target patterns.
//
// Note that these packages do not include AST files (syntax) or types info.
// We don't load them right away to avoid OOM situations for big projects.
func (r *runner) findPackages(ctx context.Context, fset *token.FileSet, targets []string) ([]packageRef, error) {
	loadMode := packages.NeedName | packages.NeedFiles
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

	if len(pkgs) == 1 {
		pkg := pkgs[0]
		if pkg.PkgPath == "command-line-arguments" && len(targets) == 1 {
			ref := packageRef{id: pkg.ID, path: targets[0]}
			return []packageRef{ref}, nil
		}
	}

	// We specify tests=false, but just in case we're still going
	// to filter the packages here.
	packageSet := make(map[string]struct{})
	var result []packageRef
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

		if len(pkg.GoFiles) == 0 {
			continue
		}

		ref := packageRef{
			id:   pkg.ID,
			name: pkg.Name,
			path: pkg.PkgPath,
		}

		if strings.HasPrefix(pkg.PkgPath, "_/") {
			absFilename := pkg.GoFiles[0]
			relFilename, err := filepath.Rel(r.wd, absFilename)
			if err != nil {
				return nil, err
			}
			ref.path = "./" + filepath.Dir(relFilename)
		}

		if _, ok := packageSet[pkg.PkgPath]; ok {
			continue
		}
		packageSet[pkg.PkgPath] = struct{}{}

		result = append(result, ref)
	}

	return result, nil
}

func (r *runner) inspectHeatmap() {
	r.heatmapPackages = make(map[string]struct{})
	r.heatmapFiles = make(map[string]struct{})
	r.heatmap.Inspect(func(l heatmap.LineStats) {
		if l.GlobalHeatLevel == 0 {
			return
		}
		r.heatmapPackages[l.Func.PkgName] = struct{}{}
		r.heatmapFiles[filepath.Base(l.Func.Filename)] = struct{}{}
	})
}

func (r *runner) createHeatmap() (*heatmap.Index, error) {
	data, err := os.ReadFile(r.args.heatmapFile)
	if err != nil {
		return nil, err
	}
	pprofProfile, err := profile.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	index := heatmap.NewIndex(heatmap.IndexConfig{
		Threshold: r.args.heatmapThreshold,
	})
	if err := index.AddProfile(pprofProfile); err != nil {
		return nil, err
	}
	return index, nil
}
