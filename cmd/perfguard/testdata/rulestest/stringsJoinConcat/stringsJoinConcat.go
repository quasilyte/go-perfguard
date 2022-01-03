package rulestest

import "strings"

func Warn(x, y, z, s1, s2 string) {
	_ = strings.Join([]string{x, y}, "_") // want `strings.Join([]string{x, y}, "_") => x + "_" + y`

	_ = strings.Join([]string{y, x}, "")    // want `strings.Join([]string{y, x}, "") => y + x`
	_ = strings.Join([]string{x, y, z}, "") // want `strings.Join([]string{x, y, z}, "") => x + y + z`

	_ = strings.Join([]string{s1, s2}, concat("<", ">")) // want `strings.Join([]string{s1, s2}, concat("<", ">")) => s1 + concat("<", ">") + s2`
}

func Ignore(s1, s2 string) {
	_ = s1 + s2
	_ = s1 + "_" + s2
	_ = strings.Join([]string{s1, s2, s1}, concat("<", ">"))
	_ = strings.Join([]string{"x", "y"}, "")
	_ = strings.Join([]string{"x", "y"}, "|")
	_ = strings.Join([]string{"x", "y", "z"}, "|")
}

func concat(x, y string) string {
	return x + y
}
