package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

// Universal rules are shared in both `lint` and `optimize` modes.
//
// By default, all rules trigger on every successful match.
// For most optimization rules, it's better to set one of the
// tags to limit its scope of application.
//
// There are two special tags for this: o1 and o2.
// o1 requires that heat level for this line is not zero
// o2 requires that heat level for this line is 5 (max level)
//
// Use o2 for rules that should be applied carefully.
// This is usually the case when optimized code is more verbose
// or generally less pretty.
//
// Lint mode ignores o1 and o2 tags completely.

//doc:summary Detects unoptimal strings/bytes case-insensitive comparison
//doc:tags    o1
//doc:before  strings.ToLower(x) == strings.ToLower(y)
//doc:after   strings.EqualFold(x, y)
func equalFold(m dsl.Matcher) {
	// string == patterns
	m.Match(
		`strings.ToLower($x) == $y`,
		`strings.ToLower($x) == strings.ToLower($y)`,
		`$x == strings.ToLower($y)`,
		`strings.ToUpper($x) == $y`,
		`strings.ToUpper($x) == strings.ToUpper($y)`,
		`$x == strings.ToUpper($y)`).
		Where(m["x"].Pure && m["y"].Pure && m["x"].Text != m["y"].Text).
		Suggest(`strings.EqualFold($x, $y)`)

	// string != patterns
	m.Match(
		`strings.ToLower($x) != $y`,
		`strings.ToLower($x) != strings.ToLower($y)`,
		`$x != strings.ToLower($y)`,
		`strings.ToUpper($x) != $y`,
		`strings.ToUpper($x) != strings.ToUpper($y)`,
		`$x != strings.ToUpper($y)`).
		Where(m["x"].Pure && m["y"].Pure && m["x"].Text != m["y"].Text).
		Suggest(`!strings.EqualFold($x, $y)`)

	// bytes.Equal patterns
	m.Match(
		`bytes.Equal(bytes.ToLower($x), $y)`,
		`bytes.Equal(bytes.ToLower($x), bytes.ToLower($y))`,
		`bytes.Equal($x, bytes.ToLower($y))`,
		`bytes.Equal(bytes.ToUpper($x), $y)`,
		`bytes.Equal(bytes.ToUpper($x), bytes.ToUpper($y))`,
		`bytes.Equal($x, bytes.ToUpper($y))`).
		Where(m["x"].Pure && m["y"].Pure && m["x"].Text != m["y"].Text).
		Suggest(`bytes.EqualFold($x, $y)`)
}

//doc:summary Detects redundant fmt.Sprint calls
//doc:tags    o1
func redundantSprint(m dsl.Matcher) {
	m.Match(`fmt.Sprint($x)`, `fmt.Sprintf("%s", $x)`, `fmt.Sprintf("%v", $x)`).
		Where(m["x"].Type.Implements(`fmt.Stringer`)).
		Suggest(`$x.String()`)

	m.Match(`fmt.Sprint($x)`, `fmt.Sprintf("%s", $x)`, `fmt.Sprintf("%v", $x)`).
		Where(m["x"].Type.Implements(`error`)).
		Suggest(`$x.Error()`)

	m.Match(`fmt.Sprint($x)`, `fmt.Sprintf("%s", $x)`, `fmt.Sprintf("%v", $x)`).
		Where(m["x"].Type.Is(`string`)).
		Suggest(`$x`)
}

//doc:summary Detect strings.Join usages that can be rewritten as a string concat
//doc:tags    o1
func stringsJoinConcat(m dsl.Matcher) {
	m.Match(`strings.Join([]string{$x, $y}, "")`).Suggest(`$x + $y`)
	m.Match(`strings.Join([]string{$x, $y, $z}, "")`).Suggest(`$x + $y + $z`)

	m.Match(`strings.Join([]string{$x, $y}, $glue)`).Suggest(`$x + $glue + $y`)

	m.Match(`strings.Join([]string{$x, $y, $z}, $glue)`).
		Where(m["glue"].Pure).
		Suggest(`$x + $glue + $y + $glue + $z`)
}

//doc:summary Detects redundant conversions between string and []byte
//doc:tags    o1
//doc:before  copy(b, []byte(s))
//doc:after   copy(b, s)
func stringCopyElim(m dsl.Matcher) {
	m.Match(`copy($b, []byte($s))`).
		Where(m["s"].Type.Is(`string`)).
		Suggest(`copy($b, $s)`)

	m.Match(`len(string($b))`).Where(m["b"].Type.Is(`[]byte`)).Suggest(`len($b)`)

	m.Match(`$re.Match([]byte($s))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["s"].Type.Is(`string`)).
		Suggest(`$re.MatchString($s)`)

	m.Match(`$re.FindIndex([]byte($s))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["s"].Type.Is(`string`)).
		Suggest(`$re.FindStringIndex($s)`)

	m.Match(`$re.FindAllIndex([]byte($s), $n)`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["s"].Type.Is(`string`)).
		Suggest(`$re.FindAllStringIndex($s, $n)`)
}

//doc:summary Detects strings.Index calls that may cause unwanted allocs
//doc:tags    o1
//doc:before  strings.Index(string(x), y)
//doc:after   bytes.Index(x, []byte(y))
//doc:note    See Go issue for details: https://github.com/golang/go/issues/25864
func indexAlloc(m dsl.Matcher) {
	m.Match(`strings.Index(string($x), $y)`).
		Where(m["x"].Pure && m["y"].Pure && m.File().Imports(`bytes`)).
		Suggest(`bytes.Index($x, []byte($y))`)
}

//doc:summary Detects WriteRune calls with rune literal argument that is single byte and reports to use WriteByte instead
//doc:tags    o1
//doc:before  w.WriteRune('\n')
//doc:after   w.WriteByte('\n')
func writeByte(m dsl.Matcher) {
	// utf8.RuneSelf:
	// characters below RuneSelf are represented as themselves in a single byte.
	const runeSelf = 0x80
	m.Match(`$w.WriteRune($c)`).
		Where(m["w"].Type.Implements("io.ByteWriter") && (m["c"].Const && m["c"].Value.Int() < runeSelf)).
		Suggest(`$w.WriteByte($c)`)
}

//doc:summary Detects slice clear loops, suggests an idiom that is recognized by the Go compiler
//doc:tags    o1
//doc:before  for i := 0; i < len(buf); i++ { buf[i] = 0 }
//doc:after   for i := range buf { buf[i] = 0 }
func sliceClear(m dsl.Matcher) {
	m.Match(`for $i := 0; $i < len($xs); $i++ { $xs[$i] = $zero }`).
		Where(m["zero"].Value.Int() == 0).
		Suggest(`for $i := range $xs { $xs[$i] = $zero }`).
		Report(`for ... { ... } => for $i := range $xs { $xs[$i] = $zero }`)
}

//doc:summary Detects expressions like []rune(s)[0] that may cause unwanted rune slice allocation
//doc:tags    o1
//doc:before  r := []rune(s)[0]
//doc:after   r, _ := utf8.DecodeRuneInString(s)
//doc:note    See Go issue for details: https://github.com/golang/go/issues/45260
func utf8DecodeRune(m dsl.Matcher) {
	// TODO: instead of File().Imports("utf8") filter we
	// want to have a way to import "utf8" package if it's not yet imported.
	// See https://github.com/quasilyte/go-ruleguard/issues/329
	// Or maybe we can run goimports (as a library?) for these cases.
	// goimports may add more diff noise though (like imports order, etc).

	m.Match(`$ch := []rune($s)[0]`).
		Where(m["s"].Type.Is(`string`) && m.File().Imports(`unicode/utf8`)).
		Suggest(`$ch, _ := utf8.DecodeRuneInString($ch)`)

	m.Match(`$ch = []rune($s)[0]`).
		Where(m["s"].Type.Is(`string`) && m.File().Imports(`unicode/utf8`)).
		Suggest(`$ch, _ = utf8.DecodeRuneInString($ch)`)

	// Without !Imports this rule will result in duplicated messages
	// for a single slice conversion.
	m.Match(`[]rune($s)[0]`).
		Where(m["s"].Type.Is(`string`) && !m.File().Imports(`unicode/utf8`)).
		Report(`use utf8.DecodeRuneInString($s) here`)
}

//doc:summary Detects fmt.Sprint(f/ln) calls which can be replaced with fmt.Fprint(f/ln)
//doc:tags    o1
//doc:before  w.Write([]byte(fmt.Sprintf("%x", 10)))
//doc:after   fmt.Fprintf(w, "%x", 10)
func fprint(m dsl.Matcher) {
	m.Match(`$w.Write([]byte(fmt.Sprint($*args)))`).
		Where(m["w"].Type.Implements("io.Writer")).
		Suggest(`fmt.Fprint($w, $args)`)

	m.Match(`$w.Write([]byte(fmt.Sprintf($*args)))`).
		Where(m["w"].Type.Implements("io.Writer")).
		Suggest(`fmt.Fprintf($w, $args)`)

	m.Match(`$w.Write([]byte(fmt.Sprintln($*args)))`).
		Where(m["w"].Type.Implements("io.Writer")).
		Suggest(`fmt.Fprintln($w, $args)`)

	m.Match(`io.WriteString($w, fmt.Sprint($*args))`).
		Suggest(`fmt.Fprint($w, $args)`)

	m.Match(`io.WriteString($w, fmt.Sprintf($*args))`).
		Suggest(`fmt.Fprintf($w, $args)`)

	m.Match(`io.WriteString($w, fmt.Sprintln($*args))`).
		Suggest(`fmt.Fprintln($w, $args)`)
}
