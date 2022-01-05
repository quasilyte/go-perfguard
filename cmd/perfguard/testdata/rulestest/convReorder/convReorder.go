package rulestest

import (
	"bytes"
	"strings"
)

func Warn(s, s2 string, b, b2 []byte) {
	_ = strings.TrimSpace(string(b)) // want `strings.TrimSpace(string(b)) => string(bytes.TrimSpace(b))`
	_ = bytes.TrimSpace([]byte(s))   // want `bytes.TrimSpace([]byte(s)) => []byte(strings.TrimSpace(s))`

	_ = strings.TrimPrefix(string(b), string(b2)) // want `strings.TrimPrefix(string(b), string(b2)) => string(bytes.TrimPrefix(b, b2))`
	_ = bytes.TrimPrefix([]byte(s), []byte(s2))   // want `bytes.TrimPrefix([]byte(s), []byte(s2)) => []byte(strings.TrimPrefix(s, s2))`
}

func Ignore(s, s2 string, b, b2 []byte) {
	_ = string(bytes.TrimSpace(b))
	_ = []byte(strings.TrimSpace(s))

	_ = string(bytes.TrimPrefix(b, b2))
	_ = []byte(strings.TrimPrefix(s, s2))
}
