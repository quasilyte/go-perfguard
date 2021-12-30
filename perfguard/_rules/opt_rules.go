package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

//doc:summary Detects regexp compilation on hot execution paths
//doc:tags    o1
func regexpCompile(m dsl.Matcher) {
	// TODO: for constant string patterns we can move the regexp compilation
	// to a global scope and use compiled var on the original call site.
	// But this can't be done by the rules.

	m.Match(
		// Explicit compilation.
		`regexp.Compile($*_)`,
		`regexp.MustCompile($*_)`,
		`regexp.CompilePOSIX($*_)`,
		`regexp.MustCompilePOSIX($*_)`,
		// Implicit compilation - these calls do a compile per call without any cache.
		`regexp.Match($*_)`,
		`regexp.MatchString($*_)`,
		`regexp.MatchReader($*_)`,
	).Report(`regexp compilation should be avoided on the hot paths`)
}
