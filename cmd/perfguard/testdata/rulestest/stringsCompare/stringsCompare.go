package rulestest

import (
	"strings"
)

func Warn(s1, s2 string) {
	_ = strings.Compare(s1, s2) == 0 // want `strings.Compare(s1, s2) == 0 => s1 == s2`
	_ = strings.Compare(s1, s2) != 0 // want `strings.Compare(s1, s2) != 0 => s1 != s2`

	_ = strings.Compare(s1, s2) >= 0  // want `strings.Compare(s1, s2) >= 0 => s1 >= s2`
	_ = strings.Compare(s1, s2) != -1 // want `strings.Compare(s1, s2) != -1 => s1 >= s2`

	_ = strings.Compare(s1, s2) <= 0 // want `strings.Compare(s1, s2) <= 0 => s1 <= s2`
	_ = strings.Compare(s1, s2) != 1 // want `strings.Compare(s1, s2) != 1 => s1 <= s2`

	_ = strings.Compare(s1, s2) == -1 // want `strings.Compare(s1, s2) == -1 => s1 < s2`
	_ = strings.Compare(s1, s2) < 0   // want `strings.Compare(s1, s2) < 0 => s1 < s2`

	_ = strings.Compare(s1, s2) == 1 // want `strings.Compare(s1, s2) == 1 => s1 > s2`
	_ = strings.Compare(s1, s2) > 0  // want `strings.Compare(s1, s2) > 0 => s1 > s2`
}

func Ignore(s1, s2 string) {
	_ = s1 == s2
	_ = s1 != s2
	_ = s1 >= s2
	_ = s1 <= s2
	_ = s1 > s2
	_ = s1 < s2
}
