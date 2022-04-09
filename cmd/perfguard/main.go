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
			Description: "print perfguard version info",
			Do:          versionMain,
		},

		{
			Name:        "env",
			Description: "print perfguard-related env variables and their values",
			Do:          envMain,
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

func envMain(args []string) {
	varInfoList := []struct {
		name    string
		comment string
	}{
		{"PERFGUARD_DEBUG", "if set to 1, enables debug prints"},
	}

	for _, info := range varInfoList {
		fmt.Printf("%s=%q # %s\n", info.name, os.Getenv(info.name), info.comment)
	}
}

func lintMain(args []string) {
	issuesCount, err := cmdLint(os.Stdout, os.Stderr, args)
	if err != nil {
		log.Fatalf("perfguard lint: error: %+v", err)
	}

	if issuesCount > 0 {
		os.Exit(1)
	}
}

func optimizeMain(args []string) {
	if err := cmdOptimize(os.Stdout, os.Stderr, args); err != nil {
		log.Fatalf("perfguard optimize: error: %+v", err)
	}
}
