package rulestest

import (
	"bytes"
	"strings"
)

func Warn() {
	var b []byte
	var s string

	_ = strings.Index(string(b), s) // want `strings.Index(string(b), s) => bytes.Index(b, []byte(s))`

	_ = strings.Index(string([]byte("12")), s) // want `strings.Index(string([]byte("12")), s) => bytes.Index([]byte("12"), []byte(s))`

	_ = strings.Index(string([]byte{'1', '2'}), s) // want `strings.Index(string([]byte{'1', '2'}), s) => bytes.Index([]byte{'1', '2'}, []byte(s))`

	_ = strings.Index(string(b), "a"+s) // want `strings.Index(string(b), "a"+s) => bytes.Index(b, []byte("a"+s))`

	_ = strings.Contains(string(b), s)  // want `strings.Contains(string(b), s) => bytes.Contains(b, []byte(s))`
	_ = bytes.Index([]byte(s), b)       // want `bytes.Index([]byte(s), b) => strings.Index(s, string(b))`
	_ = bytes.Contains([]byte(s), b)    // want `bytes.Contains([]byte(s), b) => strings.Contains(s, string(b))`
	_ = strings.HasPrefix(string(b), s) // want `strings.HasPrefix(string(b), s) => bytes.HasPrefix(b, []byte(s))`
	_ = strings.HasSuffix(string(b), s) // want `strings.HasSuffix(string(b), s) => bytes.HasSuffix(b, []byte(s))`
	_ = bytes.HasPrefix([]byte(s), b)   // want `bytes.HasPrefix([]byte(s), b) => strings.HasPrefix(s, string(b))`
	_ = bytes.HasSuffix([]byte(s), b)   // want `bytes.HasSuffix([]byte(s), b) => strings.HasSuffix(s, string(b))`
}

func Ignore() {
	var b1 []byte
	var b2 []byte
	var s1 string
	var s2 string

	_ = bytes.Index(b1, b2)
	_ = bytes.Index([]byte(b1), b2)
	_ = bytes.Index(b1, []byte(b2))
	_ = bytes.Index([]byte(b1), []byte(b2))

	_ = bytes.Index(b1, []byte(s1))
	_ = bytes.Index([]byte("12"), []byte(s1))
	_ = bytes.Index([]byte{'1', '2'}, []byte(s1))
	_ = bytes.Index(b1, []byte("a"+s2))

	_ = bytes.Index(getBytes(), []byte(s1))
	_ = bytes.Index(getBytes(), []byte("a"+s1))
	_ = bytes.Index(b1, []byte(getString()))
	_ = bytes.Index([]byte("12"), []byte(getString()))
	_ = bytes.Index([]byte{'1', '2'}, []byte(getString()))
	_ = bytes.Index(b1, []byte("a"+getString()))
	_ = strings.Index(string(getBytes()), s1)
	_ = strings.Index(string(getBytes()), "a"+s1)
	_ = strings.Index(string(b1), getString())
	_ = strings.Index(string([]byte(getString())), getString())
	_ = strings.Index(string(b1), "a"+getString())

	{
		var b []byte
		var s string
		_ = bytes.Contains(b, []byte(s))
		_ = strings.Index(s, string(b))
		_ = strings.Contains(s, string(b))
		_ = bytes.HasPrefix(b, []byte(s))
		_ = bytes.HasSuffix(b, []byte(s))
		_ = strings.HasPrefix(s, string(b))
		_ = strings.HasSuffix(s, string(b))
	}
}

func getBytes() []byte  { return nil }
func getString() string { return "" }
