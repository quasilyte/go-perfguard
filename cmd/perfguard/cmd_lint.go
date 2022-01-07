package main

import (
	"flag"
	"io"
)

func cmdLint(stdout, stderr io.Writer, args []string) error {
	r := newRunner(stdout, stderr)

	fs := flag.NewFlagSet("perfguard optimize", flag.ExitOnError)
	fs.BoolVar(&r.autofix, "fix", false,
		`apply the suggested fixes automatically, where possible`)
	fs.StringVar(&r.goVersion, "go", "",
		`select the Go version to target; leave as string for the latest`)
	fs.BoolVar(&r.absFilenames, "abs", false,
		`print absolute filenames in the output`)
	noColor := fs.Bool("no-color", false, `disable colored output`)
	_ = fs.Parse(args)

	r.targets = fs.Args()
	r.loadLintRules = true
	r.coloredOutput = !*noColor

	return r.Run()
}
