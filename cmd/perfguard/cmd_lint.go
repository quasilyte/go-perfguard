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
	_ = fs.Parse(args)

	r.targets = fs.Args()
	r.loadLintRules = true

	return r.Run()
}
