package rulestest

type myString string

func Warn(s string, s2 myString, runes []rune) {
	for _, ch := range []rune(s) { // want `range []rune(s) => range s`
		println(ch)
	}

	for _, ch := range []rune(s2) { // want `range []rune(s2) => range s2`
		println(ch)
	}

	for _, ch := range []rune(string(s2)) { // want `range []rune(string(s2)) => range string(s2)`
		println(ch)
	}

	for range []rune("foo") { // want `range []rune("foo") => range "foo"`
	}

	{
		var ch rune
		for _, ch = range []rune(s) { // want `range []rune(s) => range s`
			println(ch)
		}
		for _, ch = range []rune(s2) { // want `range []rune(s2) => range s2`
			println(ch)
		}
		for _, ch = range []rune(string(s2)) { // want `range []rune(string(s2)) => range string(s2)`
			println(ch)
		}
	}

	{
		var ch rune
		for _, ch := range string(runes) { // want `range string(runes) => range runes`
			println(ch)
		}
		for _, ch = range string(runes) { // want `range string(runes) => range runes`
			println(ch)
		}
		for _, ch := range string([]rune(runes)) { // want `range string([]rune(runes)) => range []rune(runes)`
			println(ch)
		}
	}
}

func Ignore(s string, s2 myString, runes []rune) {
	for _, ch := range s {
		println(ch)
	}

	for _, ch := range s2 {
		println(ch)
	}

	for _, ch := range string(s2) {
		println(ch)
	}

	for i, ch := range []rune(s) {
		println(ch)
		println(i)
	}

	for i, ch := range []rune(s2) {
		println(ch)
		println(i)
	}

	for i, ch := range []rune(string(s2)) {
		println(ch)
		println(i)
	}

	for range "foo" {
		println("ok")
	}

	{
		var ch rune
		for _, ch := range runes {
			println(ch)
		}
		for _, ch = range runes {
			println(ch)
		}
		for _, ch := range []rune(runes) {
			println(ch)
		}
	}
}
