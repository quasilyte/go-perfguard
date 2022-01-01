package rulestest

import (
	"strings"
)

func Warn() {
	_ = []rune("abc")[0]        // want `use utf8.DecodeRuneInString("abc") here`
	_ = []rune(makeString())[0] // want `use utf8.DecodeRuneInString(makeString()) here`
}

func Ignore() {
	{
		r := []uint64{10}
		_ = r[0]
		r2 := []rune{10, 12, 34}
		_ = r2[0]
	}

	{
		// OK: 'runes' is not string-typed.
		var runes []rune
		_ = []rune(runes)[0]
	}

	{
		// OK: not a 0 index.
		var s string
		_ = []rune(s)[1]
	}

	{
		// OK: let's allow using int32 for now?
		var s string
		_ = []int32(s)[0]
	}
}

func getBytes() []byte  { return nil }
func getString() string { return "" }

func makeString() string {
	return strings.Repeat("abc", 3)
}
