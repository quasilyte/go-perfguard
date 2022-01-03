![Build Status](https://github.com/quasilyte/go-perfguard/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/quasilyte/go-perfguard)](https://goreportcard.com/report/github.com/quasilyte/go-perfguard)

# perfguard

> This tool is a work in progress.
> It's not production-ready yet.

perfguard is a static analyzer with emphasis on performance.

There are two main modes: optimize and lint. Optimization mode uses CPU profile information to improve the analysis precision and avoid suggestions in the cold execution paths. Lint mode reports all potential performance issues.

perfguard features:

* Profile-guided analysis in "optimize" mode
* Most found issues are auto-fixable (quickfixes)
* Easy to extend with custom rules
