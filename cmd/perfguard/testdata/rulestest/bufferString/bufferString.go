package rulestest

import (
	"bytes"
	"fmt"
	"io"
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

	_ = []byte(buf.String()) // want `[]byte(buf.String()) => buf.Bytes()`

	{
		var w io.Writer
		b := &bytes.Buffer{}
		fmt.Fprint(w, b.String())        // want `fmt.Fprint(w, b.String()) => w.Write(b.Bytes())`
		fmt.Fprintf(w, "%s", b.String()) // want `fmt.Fprintf(w, "%s", b.String()) => w.Write(b.Bytes())`
		fmt.Fprintf(w, "%v", b.String()) // want `fmt.Fprintf(w, "%v", b.String()) => w.Write(b.Bytes())`
	}
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

	_ = buf.Bytes()
	_ = []byte(buf.Bytes())

	{
		var w io.Writer
		b := &bytes.Buffer{}
		w.Write(b.Bytes())
	}
}
