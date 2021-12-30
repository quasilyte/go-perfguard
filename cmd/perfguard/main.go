package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cespare/subcmd"
)

// Build* variables are initialized during the build via -ldflags.
var (
	BuildVersion string
	BuildTime    string
	BuildOSUname string
	BuildCommit  string
)

func main() {
	log.SetFlags(0)

	cmds := []subcmd.Command{
		{
			Name:        "lint",
			Description: "static analysis mode, no CPU profiles needed",
			Do:          lintMain,
		},

		{
			Name:        "optimize",
			Description: "profile-guided optimizer mode",
			Do:          optimizeMain,
		},

		{
			Name:        "version",
			Description: "print ktest version info",
			Do:          versionMain,
		},
	}

	subcmd.Run(cmds)
}

func versionMain(args []string) {
	if BuildCommit == "" {
		fmt.Printf("perfguard built without version info\n")
	} else {
		fmt.Printf("perfguard version %s\nbuilt on: %s\nos: %s\ncommit: %s\n",
			BuildVersion, BuildTime, BuildOSUname, BuildCommit)
	}
}

func lintMain(args []string) {
	if err := cmdLint(os.Stdout, os.Stderr, args); err != nil {
		log.Fatalf("perfguard lint: error: %+v", err)
	}
}

func optimizeMain(args []string) {
	if err := cmdOptimize(os.Stdout, os.Stderr, args); err != nil {
		log.Fatalf("perfguard optimize: error: %+v", err)
	}
}
