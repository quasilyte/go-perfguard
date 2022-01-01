package main

import (
	"errors"
	"flag"
	"io"
)

func cmdOptimize(stdout, stderr io.Writer, args []string) error {
	r := newRunner(stdout, stderr)

	fs := flag.NewFlagSet("perfguard optimize", flag.ExitOnError)
	fs.StringVar(&r.heatmapFile, "heatmap", "",
		`a CPU profile that will be used to build a heatmap, needed for IsHot() filters`)
	fs.Float64Var(&r.heatmapThreshold, "heatmap-threshold", 0.25,
		`a threshold argument used to create a heatmap, see perf-heatmap docs on it`)
	fs.BoolVar(&r.autofix, "fix", false,
		`apply the suggested fixes automatically, where possible`)
	fs.StringVar(&r.goVersion, "go", "",
		`select the Go version to target; leave as string for the latest`)
	_ = fs.Parse(args)

	r.targets = fs.Args()
	r.loadOptRules = true

	if r.heatmapFile == "" {
		return errors.New("CPU profile is required, see --heatmap argument")
	}

	return r.Run()
}
