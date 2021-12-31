package rulestest

import (
	"bytes"
	"fmt"
)

func Warn(s1, s2 string, buf1, buf2 *bytes.Buffer) {
	_ = fmt.Sprintf("%s%s", s1, s2)     // want `fmt.Sprintf("%s%s", s1, s2) => s1 + s2`
	_ = fmt.Sprintf("%s%s", buf1, buf2) // want `fmt.Sprintf("%s%s", buf1, buf2) => buf1.String() + buf2.String(`
}

func Ignore() {
	var s1, s2 string
	_ = s1 + s2
	_ = s1 + "_" + s2
}
