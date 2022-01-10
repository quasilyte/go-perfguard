package main

import (
	"flag"
	"io"
)

func cmdLint(stdout, stderr io.Writer, args []string) error {
	r := newRunner(stdout, stderr)

	fs := flag.NewFlagSet("perfguard optimize", flag.ExitOnError)
	addCommonFlags(r, fs)
	noColor := fs.Bool("no-color", false, `disable colored output`)
	_ = fs.Parse(args)

	r.targets = fs.Args()
	r.loadLintRules = true
	r.coloredOutput = !*noColor

	return r.Run()
}
