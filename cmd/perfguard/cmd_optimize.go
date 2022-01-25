package main

import (
	"errors"
	"flag"
	"io"
)

func cmdOptimize(stdout, stderr io.Writer, args []string) error {
	r := newRunner(stdout, stderr)

	fs := flag.NewFlagSet("perfguard optimize", flag.ExitOnError)
	addCommonFlags(r, fs)
	fs.StringVar(&r.args.heatmapFile, "heatmap", "",
		`a CPU profile that will be used to build a heatmap, needed for IsHot() filters`)
	fs.Float64Var(&r.args.heatmapThreshold, "heatmap-threshold", 0.5,
		`a threshold argument used to create a heatmap, see perf-heatmap docs on it`)
	noColor := fs.Bool("no-color", false, `disable colored output`)
	_ = fs.Parse(args)

	r.targets = fs.Args()
	r.loadOptRules = true
	r.coloredOutput = !*noColor

	if r.args.heatmapFile == "" {
		return errors.New("CPU profile is required, see --heatmap argument")
	}

	return r.Run()
}
