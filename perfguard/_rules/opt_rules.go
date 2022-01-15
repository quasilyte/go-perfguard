package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

//doc:summary Detects string concat in hot paths
//doc:tags    o2 score5
func stringConcatAssign(m dsl.Matcher) {
	m.Match(`$s += $_`).
		Where(m["s"].Type.Is(`string`)).
		Report(`string concat on the hot path`)
}

//doc:summary Detects regexp compilation on hot execution paths
//doc:tags    o1 score4
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
//doc:tags    o2 score2
func sprintfConcat2(m dsl.Matcher) {
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

//doc:summary Detects Write calls that should be rewritten as io.WriteString
//doc:tags    o2 score3
//doc:before  w.Write([]byte(s))
//doc:after   io.WriteString(w, s)
func writeString2(m dsl.Matcher) {
	m.Match(`$w.Write([]byte($s))`).
		Where(m["w"].Type.Is("io.Writer") && m["s"].Type.Is(`string`) && m["s"].Const).
		Suggest("io.WriteString($w, $s)")
}

//doc:summary Detects range loops that copy large value on every iteration
//doc:tags    o1 score2
func rangeValueCopy(m dsl.Matcher) {
	// TODO: move to a hand-written checker so we can provide a quickfix for this.
	m.Match(`for $_, $v := range $_ { $*_ }`, `for $_, $v = range $_ { $*_ }`).
		Where(m["v"].Type.Size > 128).
		Report(`every iteration copies a large object into $v`)
}

//doc:summary Detects errors.New that can be allocated exactly once
//doc:tags    o1 score3
func constErrorNew(m dsl.Matcher) {
	m.Match(`errors.New($x)`).
		Where(m["x"].Const).
		Report(`errors with const message can be a global var, allocated only once`)
}

//doc:summary Detects sync.Pool usage on non pointer objects
//doc:tags    o1 score3
func syncPoolNonPtr(m dsl.Matcher) {
	m.Match(`$x.Put($y)`).
		Where(m["x"].Type.Is("sync.Pool") && !m["y"].Type.Is("*$_")).
		Report(`don't use sync.Pool on non pointer objects`).
		At(m["y"])
}
