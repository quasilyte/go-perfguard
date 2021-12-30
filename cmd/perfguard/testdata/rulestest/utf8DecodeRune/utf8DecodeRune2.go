package rulestest

import (
	"unicode/utf8"
)

var _ = utf8.MaxRune

func Warn2() {
	_ = []rune("abc")[0]        // want `_ = []rune("abc")[0] => _, _ = utf8.DecodeRuneInString(_)`
	_ = []rune(makeString())[0] // want ` = []rune(makeString())[0] => _, _ = utf8.DecodeRuneInString(_)`

	{
		ch := []rune("abc")[0] // want `ch := []rune("abc")[0] => ch, _ := utf8.DecodeRuneInString(ch`
		_ = ch
	}

	var ch rune
	{
		ch = []rune("abc")[0] // want `ch = []rune("abc")[0] => ch, _ = utf8.DecodeRuneInString(ch)`
	}
	_ = ch
}
