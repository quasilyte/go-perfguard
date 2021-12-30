package rulestest

import (
	"regexp"
)

func Warn() {
	var b []byte
	var b2 []byte
	var s string

	copy(b, []byte(s)) // want `copy(b, []byte(s)) => copy(b, s)`

	_ = len(string(b))      // want `len(string(b)) => len(b)`
	_ = len(string(b)) == 0 // want `len(string(b)) => len(b)`

	re := regexp.MustCompile(`\w+`)

	_ = re.Match([]byte(s))            // want `re.Match([]byte(s)) => re.MatchString(s)`
	_ = re.FindIndex([]byte(s))        // want `re.FindIndex([]byte(s)) => re.FindStringIndex(s)`
	_ = re.FindAllIndex([]byte(s), -1) // want `re.FindAllIndex([]byte(s), -1) => re.FindAllStringIndex(s, -1)`
}

func Ignore() {
	var b []byte
	var s string

	copy(b, s)

	{
		copy := func(int) {}
		copy(1)

		var s string
		re := regexp.MustCompile(`\w+`)

		_ = re.MatchString(s)

		_ = re.FindStringIndex(s)

		_ = re.FindAllStringIndex(s, -1)
	}
}
