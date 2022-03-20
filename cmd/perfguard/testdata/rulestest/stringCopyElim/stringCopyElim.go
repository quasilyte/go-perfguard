package rulestest

import (
	"bytes"
	"strings"
)

func Warn() {
	var b []byte
	var b2 []byte
	var s string

	copy(b, []byte(s)) // want `copy(b, []byte(s)) => copy(b, s)`

	b = append(b, string(b2)...) // want `append(b, string(b2)...) => append(b, b2...)`

	b = append(b, []byte(s)...) // want `append(b, []byte(s)...) => append(b, s...)`

	b = append(b, []byte("okay")...) // want `append(b, []byte("okay")...) => append(b, "okay"...)`

	_ = len(string(b))       // want `len(string(b)) => len(b)`
	_ = len(string(b2)) == 0 // want `len(string(b2)) => len(b2)`

	{
		_ = []byte(strings.ToUpper(string(b))) // want `[]byte(strings.ToUpper(string(b))) => bytes.ToUpper(b)`
		_ = []byte(strings.ToLower(string(b))) // want `[]byte(strings.ToLower(string(b))) => bytes.ToLower(b)`

		_ = []byte(strings.TrimSuffix(string(b), s))                          // want `[]byte(strings.TrimSuffix(string(b), s)) => bytes.TrimSuffix(b, []byte(s))`
		_ = []byte(strings.TrimSuffix(string(b), strings.TrimPrefix(s, "/"))) // want `[]byte(strings.TrimSuffix(string(b), strings.TrimPrefix(s, "/"))) => bytes.TrimSuffix(b, []byte(strings.TrimPrefix(s, "/")))`
	}

	{
		_ = bytes.NewReader([]byte(s))   // want `bytes.NewReader([]byte(s)) => strings.NewReader(s)`
		_ = strings.NewReader(string(b)) // want `strings.NewReader(string(b)) => bytes.NewReader(b)`
	}
}

func Ignore() {
	var b []byte
	var b2 []byte
	var s string

	copy(b, s)

	b = append(b, b2...)

	{
		copy := func(int) {}
		copy(1)
	}

	{
		_ = bytes.ToUpper(b)
		_ = bytes.ToLower(b)
		_ = bytes.TrimSpace(b)

		_ = bytes.TrimSuffix(b, []byte(s))
		_ = bytes.TrimSuffix(b, []byte(strings.TrimPrefix(s, "/")))
	}

	{
		_ = strings.NewReader(s)
		_ = bytes.NewReader(b)
	}
}
