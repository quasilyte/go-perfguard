package rulestest

type myString string

func Warn(s string, s2 myString, runes []rune) {
	for _, ch := range []rune(s) { // want `for _, ch := range []rune(s) => for _, ch := range s`
		println(ch)
	}

	for _, ch := range []rune(s2) { // want `for _, ch := range []rune(s2) => for _, ch := range s2`
		println(ch)
	}

	for _, ch := range []rune(string(s2)) { // want `for _, ch := range []rune(string(s2)) => for _, ch := range string(s2)`
		println(ch)
	}

	for range []rune("foo") { // want `for range []rune("foo") => for range "foo"`
	}

	{
		var ch rune
		for _, ch = range []rune(s) { // want `for _, ch = range []rune(s) => for _, ch = range s`
			println(ch)
		}
		for _, ch = range []rune(s2) { // want `for _, ch = range []rune(s2) => for _, ch = range s2`
			println(ch)
		}
		for _, ch = range []rune(string(s2)) { // want `for _, ch = range []rune(string(s2)) => for _, ch = range string(s2)`
			println(ch)
		}
	}

	{
		var ch rune
		for _, ch := range string(runes) { // want `for _, ch := range string(runes) => for _, ch := range runes`
			println(ch)
		}
		for _, ch = range string(runes) { // want `for _, ch = range string(runes) => for _, ch = range runes`
			println(ch)
		}
		for _, ch := range string([]rune(runes)) { // want `for _, ch := range string([]rune(runes)) => for _, ch := range []rune(runes)`
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
