package rulestest

import (
	"bytes"
	"strings"
)

func Warn(s string, b []byte) {
	_ = strings.TrimSpace(string(b)) // want `strings.TrimSpace(string(b)) => string(bytes.TrimSpace(b))`
	_ = bytes.TrimSpace([]byte(s))   // want `bytes.TrimSpace([]byte(s)) => []byte(strings.TrimSpace(s))`
}

func Ignore(s string, b []byte) {
	_ = string(bytes.TrimSpace(b))
	_ = []byte(strings.TrimSpace(s))
}
