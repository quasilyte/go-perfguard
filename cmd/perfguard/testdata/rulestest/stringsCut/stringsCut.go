package rulestest

import (
	"bytes"
	"strings"
)

func Warn(s, delim string, b, bDelim []byte) {
	var sOut string

	{
		sOut = strings.Split(s, delim)[0]   // want `sOut = strings.Split(s, delim)[0] => sOut, _, _ = strings.Cut(s, delim)`
		sOut2 := strings.Split(s, delim)[0] // want `sOut2 := strings.Split(s, delim)[0] => sOut2, _, _ := strings.Cut(s, delim)`
		_ = sOut2
	}

	_ = sOut
}

func Ignore(s, delim string, b, bDelim []byte) {
	var sOut string
	var bOut []byte

	{
		sOut, _, _ = strings.Cut(s, delim)
		sOut2, _, _ := strings.Cut(s, delim)
		_ = sOut2
	}
	{
		bOut, _, _ = bytes.Cut(b, bDelim)
		bOut2, _, _ := bytes.Cut(b, bDelim)
		_ = bOut2
	}

	_ = sOut
	_ = bOut
}
