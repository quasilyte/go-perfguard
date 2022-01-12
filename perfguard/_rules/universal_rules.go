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
//
// When several changes can be applied to the same code spot,
// we need to choose how to rewrite the source code, which replacement wins.
// To achieve that, we use the scoring system.
// There are 5 score classes: from 1 to 5 inclusively.
// We use [score1, score2, score3, score4, score5] tags for this.
// If several rewrites have the same score, we pick the first one
// by using some sorting algorithm to make the decision stable.
//
// Usually, cutting extra allocations is really good, so it
// deserves to be a score3 or above (depending on how much it saves).
//
// Simple CPU optimizations that save a few nanoseconds are score1.
//
// score5 is something that can make the code several times faster
// or make it zero allocations (as opposed to the replaced form).

//doc:summary Detects use cases for strings.Cut
//doc:tags    o1 score3
//doc:before  email := strings.Split(s, "@")[0]
//doc:after   email, _, _ := strings.Cut(s, "@")
func stringsCut(m dsl.Matcher) {
	m.Match(`$dst := strings.Split($s, $sep)[0]`).
		Where(m.GoVersion().GreaterEqThan("1.18")).
		Suggest(`$dst, _, _ := strings.Cut($s, $sep)`)
	m.Match(`$dst = strings.Split($s, $sep)[0]`).
		Where(m.GoVersion().GreaterEqThan("1.18")).
		Suggest(`$dst, _, _ = strings.Cut($s, $sep)`)
}

//doc:summary Detects use cases for bytes.Cut
//doc:tags    o1 score3
//doc:before  email := bytes.Split(b, "@")[0]
//doc:after   email, _, _ := bytes.Cut(b, []byte("@"))
func bytesCut(m dsl.Matcher) {
	m.Match(`$dst := bytes.Split($b, $sep)[0]`).
		Where(m.GoVersion().GreaterEqThan("1.18")).
		Suggest(`$dst, _, _ := bytes.Cut($b, $sep)`)
	m.Match(`$dst = bytes.Split($b, $sep)[0]`).
		Where(m.GoVersion().GreaterEqThan("1.18")).
		Suggest(`$dst, _, _ = bytes.Cut($b, $sep)`)
}

//doc:summary Detects use cases for strings.Clone
//doc:tags    o1 score3
//doc:before  s2 := string([]byte(s1))
//doc:after   s2 := strings.Clone(s1)
func stringsClone(m dsl.Matcher) {
	m.Match(`string([]byte($s))`).
		Where(m["s"].Type.Is(`string`) &&
			!m["s"].Const &&
			m.GoVersion().GreaterEqThan("1.18")).
		Suggest(`strings.Clone($s)`)
}

//doc:summary Detects unoptimal strings/bytes case-insensitive comparison
//doc:tags    o1 score2
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

	// Strings prefix/suffix patterns.
	m.Match(
		`strings.HasPrefix(strings.ToLower($x), $y)`,
		`strings.HasPrefix(strings.ToUpper($x), $y)`).
		Where(m["x"].Pure && m["y"].Pure && m["x"].Text != m["y"].Text).
		Suggest(`(len($x) >= len($y) && strings.EqualFold($x[:len($y)], $y))`)
	m.Match(
		`strings.HasSuffix(strings.ToLower($x), $y)`,
		`strings.HasSuffix(strings.ToUpper($x), $y)`).
		Where(m["x"].Pure && m["y"].Pure && m["x"].Text != m["y"].Text).
		Suggest(`(len($x) >= len($y) && strings.EqualFold($x[len($x)-len($y):], $y))`)

	// Bytes prefix/suffix patterns.
	m.Match(
		`bytes.HasPrefix(bytes.ToLower($x), $y)`,
		`bytes.HasPrefix(bytes.ToUpper($x), $y)`).
		Where(m["x"].Pure && m["y"].Pure && m["x"].Text != m["y"].Text).
		Suggest(`(len($x) >= len($y) && bytes.EqualFold($x[:len($y)], $y))`)
	m.Match(
		`bytes.HasSuffix(bytes.ToLower($x), $y)`,
		`bytes.HasSuffix(bytes.ToUpper($x), $y)`).
		Where(m["x"].Pure && m["y"].Pure && m["x"].Text != m["y"].Text).
		Suggest(`(len($x) >= len($y) && bytes.EqualFold($x[len($x)-len($y):], $y))`)
}

//doc:summary Detects redundant fmt.Sprint calls
//doc:tags    o1 score3
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

	m.Match(`fmt.Sprint($x)`, `fmt.Sprintf("%s", $x)`, `fmt.Sprintf("%v", $x)`).
		Where(m["x"].Type.ConvertibleTo(`string`) && !m["x"].Type.OfKind("numeric")).
		Suggest(`string($x)`)
}

//doc:summary Detects redundant fmt.Fprint calls
//doc:tags    o1 score3
//doc:before  fmt.Fprintf(w, "%s", data)
//doc:after   w.WriteString(data.String())
func redundantFprint(m dsl.Matcher) {
	m.Match(`fmt.Fprint($w, $x)`, `fmt.Fprintf($w, "%s", $x)`, `fmt.Fprintf($w, "%v", $x)`).
		Where(m["x"].Type.Implements(`fmt.Stringer`) && m["w"].Type.Implements(`io.StringWriter`)).
		Suggest(`$w.WriteString($x.String())`)

	m.Match(`fmt.Fprint($w, $x)`, `fmt.Fprintf($w, "%s", $x)`, `fmt.Fprintf($w, "%v", $x)`).
		Where(m["x"].Type.Implements(`error`) && m["w"].Type.Implements(`io.StringWriter`)).
		Suggest(`$w.WriteString($x.Error())`)

	m.Match(`fmt.Fprint($w, $x)`, `fmt.Fprintf($w, "%s", $x)`, `fmt.Fprintf($w, "%v", $x)`).
		Where(m["x"].Type.Is(`string`) && m["w"].Type.Implements(`io.StringWriter`)).
		Suggest(`$w.WriteString($x)`)

	m.Match(`fmt.Fprint($w, $x)`, `fmt.Fprintf($w, "%s", $x)`, `fmt.Fprintf($w, "%v", $x)`).
		Where(m["x"].Type.Is(`[]byte`)).
		Suggest(`$w.Write($x)`)
}

//doc:summary Detects slice copying patterns that can be optimized
//doc:tags    o2 score2
//doc:before  dst := append([]int(nil), src...)
//doc:after   dst := make([]int, len(src)); copy(dst, src)
func sliceClone(m dsl.Matcher) {
	m.Match(`$dst = append([]$elem(nil), $src...)`, `$dst = append([]$elem{}, $src...)`).
		Where(!m["elem"].Type.HasPointers()).
		Suggest(`$dst = make([]$elem, len($src)); copy($dst, $src)`)
	m.Match(`$dst := append([]$elem(nil), $src...)`, `$dst := append([]$elem{}, $src...)`).
		Where(!m["elem"].Type.HasPointers()).
		Suggest(`$dst := make([]$elem, len($src)); copy($dst, $src)`)
}

//doc:summary Detect strings.Join usages that can be rewritten as a string concat
//doc:tags    o1 score3
func stringsJoinConcat(m dsl.Matcher) {
	m.Match(`strings.Join([]string{$x, $y}, "")`).
		Where(!m["x"].Const && !m["y"].Const).
		Suggest(`$x + $y`)
	m.Match(`strings.Join([]string{$x, $y, $z}, "")`).
		Where(!m["x"].Const && !m["y"].Const && !m["z"].Const).
		Suggest(`$x + $y + $z`)

	m.Match(`strings.Join([]string{$x, $y}, $glue)`).
		Where(!m["x"].Const && !m["y"].Const).
		Suggest(`$x + $glue + $y`)

	m.Match(`strings.Join([]string{$x, $y, $z}, $glue)`).
		Where(m["glue"].Const && !m["x"].Const && !m["y"].Const && !m["z"].Const).
		Suggest(`$x + $glue + $y + $glue + $z`)
}

//doc:summary Detects sprint calls that can be rewritten as a string concat
//doc:tags    o1 score3
//doc:before  fmt.Sprintf("%s%s", x, y)
//doc:after   x + y
func sprintConcat(m dsl.Matcher) {
	m.Match(`fmt.Sprintf("%s%s", $x, $y)`).
		Where(m["x"].Type.Is(`string`) && m["y"].Type.Is(`string`)).
		Suggest(`$x + $y`)

	m.Match(`fmt.Sprintf("%s%s", $x, $y)`).
		Where(m["x"].Type.Implements(`fmt.Stringer`) && m["y"].Type.Implements(`fmt.Stringer`)).
		Suggest(`$x.String() + $y.String()`)
}

//doc:summary Detects fmt uses that can be replaced with strconv
//doc:tags    o1 score2
//doc:before  fmt.Sprintf("%d", i)
//doc:after   strconv.Itoa(i)
func strconv(m dsl.Matcher) {
	// Sprint(x) is basically Sprintf("%v", x), so we treat it identically.

	// The most simple cases that can be converted to Itoa.
	m.Match(`fmt.Sprintf("%d", $x)`, `fmt.Sprintf("%v", $x)`, `fmt.Sprint($x)`).
		Where(m["x"].Type.Is(`int`)).Suggest(`strconv.Itoa($x)`)

	// Patterns for int64 and uint64 go first,
	// so we don't insert unnecessary conversions by the rules below.
	m.Match(`fmt.Sprintf("%d", $x)`, `fmt.Sprintf("%v", $x)`, `fmt.Sprint($x)`).
		Where(m["x"].Type.Is(`int64`)).Suggest(`strconv.FormatInt($x, 10)`)
	m.Match(`fmt.Sprintf("%x", $x)`).
		Where(m["x"].Type.Is(`int64`)).Suggest(`strconv.FormatInt($x, 16)`)
	m.Match(`fmt.Sprintf("%d", $x)`, `fmt.Sprintf("%v", $x)`, `fmt.Sprint($x)`).
		Where(m["x"].Type.Is(`uint64`)).Suggest(`strconv.FormatUint($x, 10)`)
	m.Match(`fmt.Sprintf("%x", $x)`).
		Where(m["x"].Type.Is(`uint64`)).Suggest(`strconv.FormatUint($x, 16)`)

	m.Match(`fmt.Sprintf("%d", $x)`, `fmt.Sprintf("%v", $x)`, `fmt.Sprint($x)`).
		Where(m["x"].Type.OfKind(`int`)).Suggest(`strconv.FormatInt(int64($x), 10)`)
	m.Match(`fmt.Sprintf("%x", $x)`).
		Where(m["x"].Type.OfKind(`int`)).Suggest(`strconv.FormatInt(int64($x), 16)`)

	m.Match(`fmt.Sprintf("%d", $x)`, `fmt.Sprintf("%v", $x)`, `fmt.Sprint($x)`).
		Where(m["x"].Type.OfKind(`uint`)).Suggest(`strconv.FormatUint(uint64($x), 10)`)
	m.Match(`fmt.Sprintf("%x", $x)`).
		Where(m["x"].Type.OfKind(`uint`)).Suggest(`strconv.FormatUint(uint64($x), 16)`)
}

//doc:summary Detects cases that can benefit from append-friendly APIs
//doc:tags    o1 score4
//doc:before  b = append(b, strconv.Itoa(v)...)
//doc:after   b = strconv.AppendInt(b, v, 10)
func appendAPI(m dsl.Matcher) {
	// append functions are generally much better than alternatives,
	// but we can only go so far with the rules.
	// Maybe it's worthwhile to implement more thorough analysis
	// that detects where append-style APIs can be used.

	// Not checking the fmt.Sprint cases and alike as they
	// should be handled by other rule.
	m.Match(`$b = append($b, strconv.Itoa($x)...)`).
		Suggest(`$b = strconv.AppendInt($b, int64($x), 10)`)
	m.Match(`$b = append($b, strconv.FormatInt($x, $base)...)`).
		Suggest(`$b = strconv.AppendInt($b, $x, $base)`)
	m.Match(`$b = append($b, strconv.FormatUint($x, $base)...)`).
		Suggest(`$b = strconv.AppendUint($b, $x, $base)`)

	m.Match(`$b = append($b, $t.Format($layout)...)`).
		Where(m["t"].Type.Is(`time.Time`) || m["t"].Type.Is(`*time.Time`)).
		Suggest(`$b = $t.AppendFormat($b, $layout)`)

	m.Match(`$b = append($b, $v.String()...)`).
		Where(m["v"].Type.Is(`big.Float`) || m["v"].Type.Is(`*big.Float`)).
		Suggest(`$b = $v.Append($b, 'g', 10)`)
	m.Match(`$b = append($b, $v.Text($format, $prec)...)`).
		Where(m["v"].Type.Is(`big.Float`) || m["v"].Type.Is(`*big.Float`)).
		Suggest(`$b = $v.Append($b, $format, $prec)`)

	m.Match(`$b = append($b, $v.String()...)`).
		Where(m["v"].Type.Is(`big.Int`) || m["v"].Type.Is(`*big.Int`)).
		Suggest(`$b = $v.Append($b, 10)`)
	m.Match(`$b = append($b, $v.Text($base)...)`).
		Where(m["v"].Type.Is(`big.Int`) || m["v"].Type.Is(`*big.Int`)).
		Suggest(`$b = $v.Append($b, $base)`)
}

//doc:summary Detects patterns that can be reordered to make the code faster
//doc:tags    o1 score2
//doc:before  strings.TrimSpace(string(b))
//doc:after   string(bytes.TrimSpace(b))
func convReorder(m dsl.Matcher) {
	// When we trim some string/bytes part, the result data is smaller
	// than it was before (or it's the same if nothing was trimmed).
	// This means that it's beneficial to apply a copying (allocationg)
	// conversion over the trim result, so we allocate and copy less data.

	m.Match(`strings.TrimSpace(string($b))`).
		Where(m["b"].Type.Is(`[]byte`)).
		Suggest(`string(bytes.TrimSpace($b))`)

	m.Match(`bytes.TrimSpace([]byte($s))`).
		Where(m["s"].Type.Is(`string`)).
		Suggest(`[]byte(strings.TrimSpace($s))`)

	m.Match(`strings.TrimPrefix(string($b1), string($b2))`).
		Where(m["b1"].Type.Is(`[]byte`) && m["b2"].Type.Is(`[]byte`)).
		Suggest(`string(bytes.TrimPrefix($b1, $b2))`)

	m.Match(`bytes.TrimPrefix([]byte($s1), []byte($s2))`).
		Where(m["s1"].Type.Is(`string`) && m["s2"].Type.Is(`string`)).
		Suggest(`[]byte(strings.TrimPrefix($s1, $s2))`)
}

//doc:summary Detects sliced slice copying that can be optimized
//doc:tags    o1 score3
//doc:before  string(b)[:n]
//doc:after   string(b[:n])
func slicedConv(m dsl.Matcher) {
	m.Match(`string($b)[:$n]`).
		Where(m["b"].Type.Is(`[]byte`)).
		Suggest(`string($b[:$n])`)

	m.Match(`[]byte($s)[:$n]`).
		Where(m["s"].Type.Is(`string`)).
		Suggest(`[]byte($s[:$n])`)
}

//doc:summary Detects redundant conversions between string and []byte
//doc:tags    o1 score4
//doc:before  copy(b, []byte(s))
//doc:after   copy(b, s)
func stringCopyElim(m dsl.Matcher) {
	m.Match(`copy($b, []byte($s))`).
		Where(m["s"].Type.Is(`string`)).
		Suggest(`copy($b, $s)`)

	m.Match(`append($b, []byte($s)...)`).
		Where(m["s"].Type.Is(`string`)).
		Suggest(`append($b, $s...)`)

	m.Match(`len(string($b))`).Where(m["b"].Type.Is(`[]byte`)).Suggest(`len($b)`)

	m.Match(`[]byte(strings.$f(string($b)))`).
		Where(m["b"].Type.Is(`[]byte`) &&
			m["f"].Text.Matches(`ToUpper|ToLower|TrimSpace`)).
		Suggest(`bytes.$f($b)`)

	m.Match(`[]byte(strings.$f(string($b), $s2))`).
		Where(m["b"].Type.Is(`[]byte`) &&
			m["f"].Text.Matches(`TrimPrefix|TrimSuffix`)).
		Suggest(`bytes.$f($b, []byte($s2))`)
}

//doc:summary Detects inefficient regexp usage in regard to string/[]byte conversions
//doc:tags    o1 score3
//doc:before  regexp.ReplaceAll([]byte(s), []byte("foo"))
//doc:after   regexp.ReplaceAllString(s, "foo")
func regexpStringCopyElim(m dsl.Matcher) {
	// Cases where []byte(s) is used.

	m.Match(`$re.Match([]byte($s))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["s"].Type.Is(`string`)).
		Suggest(`$re.MatchString($s)`)

	m.Match(`$re.FindIndex([]byte($s))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["s"].Type.Is(`string`)).
		Suggest(`$re.FindStringIndex($s)`)

	m.Match(`$re.FindAllIndex([]byte($s), $n)`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["s"].Type.Is(`string`)).
		Suggest(`$re.FindAllStringIndex($s, $n)`)

	m.Match(`string($re.ReplaceAll([]byte($s), []byte($s2)))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["s"].Type.Is(`string`) && m["s2"].Type.Is(`string`)).
		Suggest(`$re.ReplaceAllString($s, $s2)`)

	m.Match(`string($re.ReplaceAll([]byte($s), $b))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["s"].Type.Is(`string`) && m["b"].Type.Is(`[]byte`)).
		Suggest(`$re.ReplaceAllString($s, string($b))`)

	// Cases where string(b) is used.

	m.Match(`$re.MatchString(string($b))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["b"].Type.Is(`[]byte`)).
		Suggest(`$re.Match($b)`)

	m.Match(`$re.FindStringIndex(string($b))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["b"].Type.Is(`[]byte`)).
		Suggest(`$re.FindIndex($b)`)

	m.Match(`$re.FindAllStringIndex(string($b), $n)`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["b"].Type.Is(`[]byte`)).
		Suggest(`$re.FindAllIndex($b, $n)`)

	m.Match(`[]byte($re.ReplaceAllString(string($b), string($b2)))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["b"].Type.Is(`[]byte`) && m["b2"].Type.Is(`[]byte`)).
		Suggest(`$re.ReplaceAll($b, $b2)`)

	m.Match(`[]byte($re.ReplaceAllString(string($b), $s))`).
		Where(m["re"].Type.Is(`*regexp.Regexp`) && m["b"].Type.Is(`[]byte`) && m["s"].Type.Is(`string`)).
		Suggest(`$re.ReplaceAll($b, []byte($s))`)
}

//doc:summary Detects strings.Index()-like calls that may allocate more than they should
//doc:tags    o1 score3
//doc:before  strings.Index(string(x), y)
//doc:after   bytes.Index(x, []byte(y))
//doc:note    See Go issue for details: https://github.com/golang/go/issues/25864
func indexAlloc(m dsl.Matcher) {
	// These rules work on the observation that substr/search item
	// is usually smaller than the containing string.

	m.Match(`strings.$f(string($b1), string($b2))`).
		Where(m["f"].Text.Matches(`Compare|Contains|HasPrefix|HasSuffix|EqualFold`) &&
			m["b1"].Type.Is(`[]byte`) && m["b2"].Type.Is(`[]byte`)).
		Suggest(`bytes.$f($b1, $b2)`)

	m.Match(`bytes.$f([]byte($s1), []byte($s2))`).
		Where(m["f"].Text.Matches(`Compare|Contains|HasPrefix|HasSuffix|EqualFold`) &&
			m["s1"].Type.Is(`string`) && m["s2"].Type.Is(`string`)).
		Suggest(`strings.$f($s1, $s2)`)

	canOptimizeStrings := func(m dsl.Matcher) bool {
		return m["x"].Pure && m["y"].Pure &&
			!m["y"].Node.Is(`CallExpr`) &&
			m["x"].Type.Is(`[]byte`)
	}

	m.Match(`strings.Index(string($x), $y)`).Where(canOptimizeStrings(m)).Suggest(`bytes.Index($x, []byte($y))`)
	m.Match(`strings.Contains(string($x), $y)`).Where(canOptimizeStrings(m)).Suggest(`bytes.Contains($x, []byte($y))`)
	m.Match(`strings.HasPrefix(string($x), $y)`).Where(canOptimizeStrings(m)).Suggest(`bytes.HasPrefix($x, []byte($y))`)
	m.Match(`strings.HasSuffix(string($x), $y)`).Where(canOptimizeStrings(m)).Suggest(`bytes.HasSuffix($x, []byte($y))`)

	canOptimizeBytes := func(m dsl.Matcher) bool {
		return m["x"].Pure && m["y"].Pure &&
			!m["y"].Node.Is(`CallExpr`) &&
			m["x"].Type.Is(`string`)
	}

	m.Match(`bytes.Index([]byte($x), $y)`).Where(canOptimizeBytes(m)).Suggest(`strings.Index($x, string($y))`)
	m.Match(`bytes.Contains([]byte($x), $y)`).Where(canOptimizeBytes(m)).Suggest(`strings.Contains($x, string($y))`)
	m.Match(`bytes.HasPrefix([]byte($x), $y)`).Where(canOptimizeBytes(m)).Suggest(`strings.HasPrefix($x, string($y))`)
	m.Match(`bytes.HasSuffix([]byte($x), $y)`).Where(canOptimizeBytes(m)).Suggest(`strings.HasSuffix($x, string($y))`)
}

//doc:summary Detects WriteRune calls with rune literal argument that is single byte and reports to use WriteByte instead
//doc:tags    o1 score1
//doc:before  w.WriteRune('\n')
//doc:after   w.WriteByte('\n')
func writeByte(m dsl.Matcher) {
	// utf8.RuneSelf:
	// characters below RuneSelf are represented as themselves in a single byte.
	const runeSelf = 0x80
	m.Match(`$w.WriteRune($c)`).
		Where(m["w"].Type.HasMethod(`io.ByteWriter.WriteByte`) && (m["c"].Const && m["c"].Value.Int() < runeSelf)).
		Suggest(`$w.WriteByte($c)`)
}

//doc:summary Detects slice clear loops, suggests an idiom that is recognized by the Go compiler
//doc:tags    o1 score2
//doc:before  for i := 0; i < len(buf); i++ { buf[i] = 0 }
//doc:after   for i := range buf { buf[i] = 0 }
func sliceClear(m dsl.Matcher) {
	m.Match(`for $i := 0; $i < len($xs); $i++ { $xs[$i] = $zero }`).
		Where(m["zero"].Value.Int() == 0).
		Suggest(`for $i := range $xs { $xs[$i] = $zero }`).
		Report(`for ... { ... } => for $i := range $xs { $xs[$i] = $zero }`)
}

//doc:summary Detects expressions like []rune(s)[0] that may cause unwanted rune slice allocation
//doc:tags    o1 score4
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
//doc:tags    o1 score2
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

//doc:summary Detects w.Write calls which can be replaced with w.WriteString
//doc:tags    o1 score4
//doc:before  w.Write([]byte("foo"))
//doc:after   w.WriteString("foo")
func writeString(m dsl.Matcher) {
	m.Match(`$w.Write([]byte($s))`).
		Where(m["w"].Type.HasMethod("io.StringWriter.WriteString") && m["s"].Type.Is(`string`)).
		Suggest("$w.WriteString($s)")
}

//doc:summary Detects w.WriteString calls which can be replaced with w.Write
//doc:tags    o1 score4
//doc:before  w.WriteString(buf.String())
//doc:after   w.Write(buf.Bytes())
func writeBytes(m dsl.Matcher) {
	isBuffer := func(v dsl.Var) bool {
		return v.Type.Is(`bytes.Buffer`) || v.Type.Is(`*bytes.Buffer`)
	}

	m.Match(`io.WriteString($w, $buf.String())`).
		Where(isBuffer(m["buf"])).
		Suggest(`$w.Write($buf.Bytes())`)

	m.Match(`io.WriteString($w, string($buf.Bytes()))`).
		Where(isBuffer(m["buf"])).
		Suggest(`$w.Write($buf.Bytes())`)

	m.Match(`$w.WriteString($buf.String())`).
		Where(m["w"].Type.HasMethod("io.Writer.Write") && isBuffer(m["buf"])).
		Suggest(`$w.Write($buf.Bytes())`)

	m.Match(`$w.WriteString(string($b))`).
		Where(m["w"].Type.HasMethod("io.Writer.Write") && m["b"].Type.Is(`[]byte`)).
		Suggest("$w.Write($b)")
}

//doc:summary Detects bytes.Buffer String() calls where Bytes() could be used instead
//doc:tags    o1 score4
//doc:before  strings.Contains(buf.String(), string(b))
//doc:after   bytes.Contains(buf.Bytes(), b)
func bufferString(m dsl.Matcher) {
	isBuffer := func(v dsl.Var) bool {
		return v.Type.Is(`bytes.Buffer`) || v.Type.Is(`*bytes.Buffer`)
	}

	m.Match(`strings.$f($buf1.String(), $buf2.String())`).
		Where(
			isBuffer(m["buf1"]) && isBuffer(m["buf2"]) &&
				m["f"].Text.Matches(`Compare|Contains|HasPrefix|HasSuffix|EqualFold`),
		).
		Suggest(`bytes.$f($buf1.Bytes(), $buf2.Bytes())`)

	m.Match(`strings.Contains($buf.String(), string($b))`).
		Where(isBuffer(m["buf"]) && m["b"].Type.Is(`[]byte`)).
		Suggest(`bytes.Contains($buf.Bytes(), $b)`)
	m.Match(`strings.HasPrefix($buf.String(), string($b))`).
		Where(isBuffer(m["buf"]) && m["b"].Type.Is(`[]byte`)).
		Suggest(`bytes.HasPrefix($buf.Bytes(), $b)`)
	m.Match(`strings.HasSuffix($buf.String(), string($b))`).
		Where(isBuffer(m["buf"]) && m["b"].Type.Is(`[]byte`)).
		Suggest(`bytes.HasSuffix($buf.Bytes(), $b)`)
	m.Match(`strings.Count($buf.String(), string($b))`).
		Where(isBuffer(m["buf"]) && m["b"].Type.Is(`[]byte`)).
		Suggest(`bytes.Count($buf.Bytes(), $b)`)
	m.Match(`strings.Index($buf.String(), string($b))`).
		Where(isBuffer(m["buf"]) && m["b"].Type.Is(`[]byte`)).
		Suggest(`bytes.Index($buf.Bytes(), $b)`)
	m.Match(`strings.EqualFold($buf.String(), string($b))`).
		Where(isBuffer(m["buf"]) && m["b"].Type.Is(`[]byte`)).
		Suggest(`bytes.EqualFold($buf.Bytes(), $b)`)

	m.Match(`strings.Contains($buf.String(), $s)`).
		Where(isBuffer(m["buf"]) && m["s"].Type.Is(`string`)).
		Suggest(`bytes.Contains($buf.Bytes(), []byte($s))`)
	m.Match(`strings.HasPrefix($buf.String(), $s)`).
		Where(isBuffer(m["buf"]) && m["s"].Type.Is(`string`)).
		Suggest(`bytes.HasPrefix($buf.Bytes(), []byte($s))`)
	m.Match(`strings.HasSuffix($buf.String(), $s)`).
		Where(isBuffer(m["buf"]) && m["s"].Type.Is(`string`)).
		Suggest(`bytes.HasSuffix($buf.Bytes(), []byte($s))`)
	m.Match(`strings.Count($buf.String(), $s)`).
		Where(isBuffer(m["buf"]) && m["s"].Type.Is(`string`)).
		Suggest(`bytes.Count($buf.Bytes(), []byte($s))`)
	m.Match(`strings.Index($buf.String(), $s)`).
		Where(isBuffer(m["buf"]) && m["s"].Type.Is(`string`)).
		Suggest(`bytes.Index($buf.Bytes(), []byte($s))`)
	m.Match(`strings.EqualFold($buf.String(), $s)`).
		Where(isBuffer(m["buf"]) && m["s"].Type.Is(`string`)).
		Suggest(`bytes.EqualFold($buf.Bytes(), []byte($s))`)

	m.Match(`[]byte($buf.String())`).Where(isBuffer(m["buf"])).Suggest(`$buf.Bytes()`)

	m.Match(`fmt.Fprint($w, $buf.String())`, `fmt.Fprintf($w, "%s", $buf.String())`, `fmt.Fprintf($w, "%v", $buf.String())`).
		Where(isBuffer(m["buf"])).
		Suggest(`$w.Write($buf.Bytes())`)
}

//doc:summary Detects array range loops that result in an excessive full data copy
//doc:tags    o1 score2
func rangeExprCopy(m dsl.Matcher) {
	m.Match(`for $_, $_ := range $e`, `for $_, $_ = range $e`).
		Where(m["e"].Addressable && m["e"].Type.Is(`[$_]$_`) && m["e"].Type.Size > 2048).
		Suggest(`&$e`).
		At(m["e"])

	// Same rule, but without Addressable requirement.
	// We can't suggest a simple fix, but we'll give a warning anyway.
	m.Match(`for $_, $_ := range $e`, `for $_, $_ = range $e`).
		Where(m["e"].Type.Is(`[$_]$_`) && m["e"].Type.Size > 2048).
		Report(`range over big array value expression is ineffective`).
		At(m["e"])
}

//doc:summary Detects range loops that can be turned into a single append call
//doc:tags    o1 score3
func rangeToAppend(m dsl.Matcher) {
	m.Match(`for $_, $x := range $src { $dst = append($dst, $x) }`).
		Where(m["src"].Type.Is(`[]$_`)).
		Suggest(`$dst = append($dst, $src...)`).
		Report(`for ... { ... } => $dst = append($dst, $src...)`)
}

//doc:summary Detects a range over []rune(string) where copying to a new slice is redundant
//doc:tags    o1 score3
func rangeRuneSlice(m dsl.Matcher) {
	m.Match(`for _, $r := range []rune($s)`).
		Where(m["s"].Type.Underlying().Is(`string`)).
		Suggest(`for _, $r := range $s`)

	m.Match(`for _, $r = range []rune($s)`).
		Where(m["s"].Type.Underlying().Is(`string`)).
		Suggest(`for _, $r = range $s`)

	m.Match(`for range []rune($s)`).
		Where(m["s"].Type.Underlying().Is(`string`)).
		Suggest(`for range $s`)

	m.Match(`for _, $r := range string($runes)`).
		Where(m["runes"].Type.Underlying().Is(`[]rune`)).
		Suggest(`for _, $r := range $runes`)

	m.Match(`for _, $r = range string($runes)`).
		Where(m["runes"].Type.Underlying().Is(`[]rune`)).
		Suggest(`for _, $r = range $runes`)
}

//doc:summary Detects usages of reflect.DeepEqual that can be rewritten
//doc:tags    o1 score2
func reflectDeepEqual(m dsl.Matcher) {
	m.Match(`reflect.DeepEqual($x, $y)`).
		Where(m["x"].Type.Is(`[]byte`) && m["y"].Type.Is(`[]byte`)).
		Suggest(`bytes.Equal($x, $y)`)

	// We insert extra () around $x == $y to ensure that we don't break
	// the code in case it was actually `!reflect.DeepEqual(...)`.
	// Without parens we would end up in `!$x == $y`, which makes no sense.
	// TODO: figure out to do it better?
	m.Match(`reflect.DeepEqual($x, $y)`).
		Where((m["x"].Type.Is(`string`) && m["y"].Type.Is(`string`)) ||
			(m["x"].Type.OfKind(`numeric`) && m["y"].Type.OfKind(`numeric`))).
		Suggest(`($x == $y)`)
}

//doc:summary Detects reflect Type() related patterns that can be optimized
//doc:tags    o1 score1
func reflectType(m dsl.Matcher) {
	m.Match(`reflect.ValueOf($x).Type()`).Suggest(`reflect.TypeOf($x)`)

	m.Match(`reflect.TypeOf($x.Interface())`).
		Where(m["x"].Type.Is(`reflect.Value`)).
		Suggest(`$x.Type()`)

	m.Match(`fmt.Sprintf("%T", $x.Interface())`).
		Where(m["x"].Type.Is(`reflect.Value`)).
		Suggest(`$x.Type().String()`)
	m.Match(`fmt.Sprintf("%T", $x)`).
		Suggest(`reflect.TypeOf($x).String()`)
}

//doc:summary Detects binary.Write uses that can be optimized
//doc:tags    o1 score3
func binaryWrite(m dsl.Matcher) {
	m.Match(`$err := binary.Write($w, $_, $b)`).
		Where(m["b"].Type.Is(`[]byte`)).
		Suggest(`_, $err := $w.Write($b)`)

	m.Match(`binary.Write($w, $_, $b)`).
		Where(m["$$"].Node.Parent().Is(`ExprStmt`) && m["b"].Type.Is(`[]byte`)).
		Suggest(`$w.Write($b)`)
}
