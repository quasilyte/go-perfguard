package imports

// Fix tries to fix imports from a Go file source code.
//
// It only works with src that can be parsed by a Go parser,
// otherwise a parsing error will be returned.
//
// We perform these operations:
// 1. Remove unused imports from a file.
// 2. Add missing imports.
//
// It does not introduce any formatting changes, unless strictly necessary.
func Fix(config FixConfig, src []byte) ([]byte, error) {
	f := newFixer(config)
	return f.Fix(src)
}

type FixConfig struct {
	StdlibPackages map[string]string
	Packages       map[string]string
}
