package rulestest

import (
	"bytes"
	"strings"
)

func Warn() {
	var x []byte
	var y string

	_ = strings.Index(string(x), y) // want `strings.Index(string(x), y) => bytes.Index(x, []byte(y))`

	_ = strings.Index(string([]byte("12")), y) // want `strings.Index(string([]byte("12")), y) => bytes.Index([]byte("12"), []byte(y))`

	_ = strings.Index(string([]byte{'1', '2'}), y) // want `strings.Index(string([]byte{'1', '2'}), y) => bytes.Index([]byte{'1', '2'}, []byte(y))`

	_ = strings.Index(string(x), "a"+y) // want `strings.Index(string(x), "a"+y) => bytes.Index(x, []byte("a"+y))`
}

func Ignore() {
	var x []byte
	var y string

	_ = bytes.Index(x, []byte(y))
	_ = bytes.Index([]byte("12"), []byte(y))
	_ = bytes.Index([]byte{'1', '2'}, []byte(y))
	_ = bytes.Index(x, []byte("a"+y))

	_ = bytes.Index(getBytes(), []byte(y))
	_ = bytes.Index(getBytes(), []byte("a"+y))
	_ = bytes.Index(x, []byte(getString()))
	_ = bytes.Index([]byte("12"), []byte(getString()))
	_ = bytes.Index([]byte{'1', '2'}, []byte(getString()))
	_ = bytes.Index(x, []byte("a"+getString()))
	_ = strings.Index(string(getBytes()), y)
	_ = strings.Index(string(getBytes()), "a"+y)
	_ = strings.Index(string(x), getString())
	_ = strings.Index(string([]byte(getString())), getString())
	_ = strings.Index(string(x), "a"+getString())
}

func getBytes() []byte  { return nil }
func getString() string { return "" }
