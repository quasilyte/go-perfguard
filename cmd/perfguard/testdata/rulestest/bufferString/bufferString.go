package rulestest

import (
	"bytes"
	"strings"
)

func Warn(buf *bytes.Buffer, s string, b []byte) {
	_ = strings.Contains(buf.String(), string(b))  // want `strings.Contains(buf.String(), string(b)) => bytes.Contains(buf.Bytes(), b)`
	_ = strings.HasPrefix(buf.String(), string(b)) // want `strings.HasPrefix(buf.String(), string(b)) => bytes.HasPrefix(buf.Bytes(), b)`
	_ = strings.HasSuffix(buf.String(), string(b)) // want `strings.HasSuffix(buf.String(), string(b)) => bytes.HasSuffix(buf.Bytes(), b)`
	_ = strings.Count(buf.String(), string(b))     // want `strings.Count(buf.String(), string(b)) => bytes.Count(buf.Bytes(), b)`

	_ = strings.Contains(buf.String(), s)  // want `strings.Contains(buf.String(), s) => bytes.Contains(buf.Bytes(), []byte(s))`
	_ = strings.HasPrefix(buf.String(), s) // want `strings.HasPrefix(buf.String(), s) => bytes.HasPrefix(buf.Bytes(), []byte(s))`
	_ = strings.HasSuffix(buf.String(), s) // want `strings.HasSuffix(buf.String(), s) => bytes.HasSuffix(buf.Bytes(), []byte(s))`
	_ = strings.Count(buf.String(), s)     // want `strings.Count(buf.String(), s) => bytes.Count(buf.Bytes(), []byte(s))`
}

func Ignore(buf *bytes.Buffer, s string, b []byte) {
	_ = bytes.Contains(buf.Bytes(), b)
	_ = bytes.HasPrefix(buf.Bytes(), b)
	_ = bytes.HasSuffix(buf.Bytes(), b)
	_ = bytes.Count(buf.Bytes(), b)

	_ = bytes.Contains(buf.Bytes(), []byte(s))
	_ = bytes.HasPrefix(buf.Bytes(), []byte(s))
	_ = bytes.HasSuffix(buf.Bytes(), []byte(s))
	_ = bytes.Count(buf.Bytes(), []byte(s))
}
