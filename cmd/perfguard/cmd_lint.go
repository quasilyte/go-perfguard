package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

var ErrIssuesFound = errors.New("found issues")

func cmdLint(stdout, stderr io.Writer, args []string) error {
	r := newRunner(stdout, stderr)

	fs := flag.NewFlagSet("perfguard optimize", flag.ExitOnError)
	addCommonFlags(r, fs)
	noColor := fs.Bool("no-color", false, `disable colored output`)
	_ = fs.Parse(args)

	r.targets = fs.Args()
	r.loadLintRules = true
	r.coloredOutput = !*noColor
	if err := r.Run(); err != nil {
		return err
	}

	if r.stats.issuesTotal != 0 {
		suffix := "auto-fixable"
		if r.autofix {
			suffix = "fixed"
		}

		return fmt.Errorf("%w: %d (%d %s)", ErrIssuesFound, r.stats.issuesTotal, r.stats.issuesFixable, suffix)
	}

	return nil
}
