package rulestest

import (
	"strings"
)

func Warn(s string) {
	_ = string([]byte(s)) // want `string([]byte(s)) => strings.Clone(s)`
}

func Ignore(s string, b []byte) {
	_ = strings.Clone(s)
	_ = string(b)
	_ = string([]byte(b))
	_ = strings([]byte("literal"))
}
