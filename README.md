![Build Status](https://github.com/quasilyte/go-perfguard/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/quasilyte/go-perfguard)](https://goreportcard.com/report/github.com/quasilyte/go-perfguard)

# perfguard

> This tool is a work in progress.
> It's not fully production-ready yet, but you can try it out.

## Overview

perfguard is a Go static analyzer with an emphasis on performance.

It supports two run modes:

1. `perfguard lint` finds potential issues, works like traditional static analysis
2. `perfguard optimize` uses CPU profiles to improve the analysis precision

perfguard key features:

* Profile-guided analysis in `perfguard optimize` mode
* Most found issues are auto-fixable with `--fix` argument (quickfixes)
* Easy to extend with custom rules (no recompilation needed)
* Can analyze big projects* even if they have some compilation errors

> (*) It doesn't try to load analysis targets into memory all at once.

Here are some examples of what it can do for you:

* Remove redundant data copying or make it faster
* Reduce the amounts of heap allocations
* Suggest more optimized functions or types from stdlib
* Recognize expensive operations in hot paths that can be lifted

## Installation

Installing from source:

```bash
# Installs a `perfguard` binary under your `$(go env GOPATH)/bin`
$ go install -v github.com/quasilyte/go-perfguard/cmd/perfguard@latest
```

## Using perfguard

It's recommended that you collect CPU profiles on realistic workflows.

For a short-lived CLI app it could be a full run. For a long-living app you may want to turn the profiling on for a minute or more, then save it to a file.

Profiles that are obtained from benchmarks are not representative and may lead to suboptimal results.

Hot spots in the profile may appear in three main places:

1. Standard Go library and the runtime. We can't apply fixes to that
2. Your app (or library) own code
3. Your code dependencies (direct or indirect)

Optimizing your own code is straightforward. Run perfguard on the root of your project:

```bash
$ perfguard optimize --heatmap cpu.out ./...
```

This will only suggest fixes to the `(2)` category.

To optimize the code from `(3)` we have several choices.

1. Optimize the library itself
2. Optimize the whole code base with an explicit vendor

The first option is preferable. You can use the same CPU profile to optimize the library. Run the perfguard on the library source code root just like you did with your application.

The second option can work for the cases when you want to deploy an optimized binary while not having a way to fix dependencies using the first option. Follow these steps:

```bash
# Make dependencies easily available for perfguard.
$ go mod vendor
# Run the analysis over the vendor.
# We use --fix argument to immediately apply the suggested changes.
$ perfguard optimize --heatmap cpu.out --fix ./vendor/...
# Build the optimized binary.
$ go build -o bin/app ./cmd/myapp
```

Then you can revert the changes to the `./vendor` or remove it if you're not using vendoring.
