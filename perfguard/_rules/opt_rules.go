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

//doc:summary Detects sprint calls that can be rewritten as a string concat
//doc:tags    o2
func sprintConcat2(m dsl.Matcher) {
	// It's impractical to implement this kind of analysis via the rules.
	// I've added a few most common patterns here just in case, but
	// we need to make a generalized form of this optimization later.

	m.Match(`fmt.Sprintf("%s=%s", $x, $y)`).
		Where(m["x"].Type.Is(`string`) && m["y"].Type.Is(`string`)).
		Suggest(`$x + "=" + $y`)

	m.Match(`fmt.Sprintf("%s.%s", $x, $y)`).
		Where(m["x"].Type.Is(`string`) && m["y"].Type.Is(`string`)).
		Suggest(`$x + "." + $y`)

	m.Match(`fmt.Sprintf("%s/%s", $x, $y)`).
		Where(m["x"].Type.Is(`string`) && m["y"].Type.Is(`string`)).
		Suggest(`$x + "/" + $y`)

	m.Match(`fmt.Sprintf("%s:%s", $x, $y)`).
		Where(m["x"].Type.Is(`string`) && m["y"].Type.Is(`string`)).
		Suggest(`$x + ":" + $y`)

	m.Match(`fmt.Sprintf("%s: %s", $x, $y)`).
		Where(m["x"].Type.Is(`string`) && m["y"].Type.Is(`string`)).
		Suggest(`$x + ": " + $y`)
}
