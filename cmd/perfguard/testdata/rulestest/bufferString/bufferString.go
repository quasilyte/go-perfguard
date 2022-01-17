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
	_ = strings.Index(buf.String(), string(b))     // want `strings.Index(buf.String(), string(b)) => bytes.Index(buf.Bytes(), b)`

	_ = strings.Contains(buf.String(), s)  // want `strings.Contains(buf.String(), s) => bytes.Contains(buf.Bytes(), []byte(s))`
	_ = strings.HasPrefix(buf.String(), s) // want `strings.HasPrefix(buf.String(), s) => bytes.HasPrefix(buf.Bytes(), []byte(s))`
	_ = strings.HasSuffix(buf.String(), s) // want `strings.HasSuffix(buf.String(), s) => bytes.HasSuffix(buf.Bytes(), []byte(s))`
	_ = strings.Count(buf.String(), s)     // want `strings.Count(buf.String(), s) => bytes.Count(buf.Bytes(), []byte(s))`
	_ = strings.Index(buf.String(), s)     // want `strings.Index(buf.String(), s) => bytes.Index(buf.Bytes(), []byte(s))`

	_ = []byte(buf.String()) // want `[]byte(buf.String()) => buf.Bytes()`

	{
		var w io.Writer
		b := bytes.NewBuffer(make([]byte, 10))
		fmt.Fprint(w, b.String())        // want `fmt.Fprint(w, b.String()) => w.Write(b.Bytes())`
		fmt.Fprintf(w, "%s", b.String()) // want `fmt.Fprintf(w, "%s", b.String()) => w.Write(b.Bytes())`
		fmt.Fprintf(w, "%v", b.String()) // want `fmt.Fprintf(w, "%v", b.String()) => w.Write(b.Bytes())`
	}

	{
		var buf1 bytes.Buffer
		var buf2 bytes.Buffer
		_ = strings.Compare(buf1.String(), buf2.String())   // want `strings.Compare(buf1.String(), buf2.String()) => bytes.Compare(buf1.Bytes(), buf2.Bytes())`
		_ = strings.Contains(buf1.String(), buf2.String())  // want `strings.Contains(buf1.String(), buf2.String()) => bytes.Contains(buf1.Bytes(), buf2.Bytes())`
		_ = strings.HasPrefix(buf1.String(), buf2.String()) // want `strings.HasPrefix(buf1.String(), buf2.String()) => bytes.HasPrefix(buf1.Bytes(), buf2.Bytes())`
		_ = strings.HasSuffix(buf1.String(), buf2.String()) // want `strings.HasSuffix(buf1.String(), buf2.String()) => bytes.HasSuffix(buf1.Bytes(), buf2.Bytes())`
		_ = strings.EqualFold(buf1.String(), buf2.String()) // want `strings.EqualFold(buf1.String(), buf2.String()) => bytes.EqualFold(buf1.Bytes(), buf2.Bytes())`
		_ = buf1.Bytes()
		_ = buf2.Bytes()
	}

	{
		var buf1 bytes.Buffer
		_ = strings.Contains(buf1.String(), "foo")  // want `strings.Contains(buf1.String(), "foo") => bytes.Contains(buf1.Bytes(), []byte("foo"))`
		_ = strings.HasPrefix(buf1.String(), "foo") // want `trings.HasPrefix(buf1.String(), "foo") => bytes.HasPrefix(buf1.Bytes(), []byte("foo"))`
		_ = strings.HasSuffix(buf1.String(), "foo") // want `strings.HasSuffix(buf1.String(), "foo") => bytes.HasSuffix(buf1.Bytes(), []byte("foo"))`
		_ = strings.EqualFold(buf1.String(), "foo") // want `strings.EqualFold(buf1.String(), "foo") => bytes.EqualFold(buf1.Bytes(), []byte("foo"))`
		_ = buf1.Bytes()
	}
}

func Ignore(buf *bytes.Buffer, s string, b []byte) {
	_ = bytes.Contains(buf.Bytes(), b)
	_ = bytes.HasPrefix(buf.Bytes(), b)
	_ = bytes.HasSuffix(buf.Bytes(), b)
	_ = bytes.Count(buf.Bytes(), b)
	_ = bytes.Index(buf.Bytes(), b)

	_ = bytes.Contains(buf.Bytes(), []byte(s))
	_ = bytes.HasPrefix(buf.Bytes(), []byte(s))
	_ = bytes.HasSuffix(buf.Bytes(), []byte(s))
	_ = bytes.Count(buf.Bytes(), []byte(s))
	_ = bytes.Index(buf.Bytes(), []byte(s))

	_ = buf.Bytes()
	_ = []byte(buf.Bytes())

	{
		var w io.Writer
		b := &bytes.Buffer{}
		w.Write(b.Bytes())
	}

	{
		var buf1 bytes.Buffer
		var buf2 bytes.Buffer
		_ = bytes.Compare(buf1.Bytes(), buf2.Bytes())
		_ = bytes.Contains(buf1.Bytes(), buf2.Bytes())
		_ = bytes.HasPrefix(buf1.Bytes(), buf2.Bytes())
		_ = bytes.HasSuffix(buf1.Bytes(), buf2.Bytes())
		_ = bytes.EqualFold(buf1.Bytes(), buf2.Bytes())
	}

	{
		var buf1 bytes.Buffer
		_ = bytes.Contains(buf1.Bytes(), []byte("foo"))
		_ = bytes.HasPrefix(buf1.Bytes(), []byte("foo"))
		_ = bytes.HasSuffix(buf1.Bytes(), []byte("foo"))
		_ = bytes.EqualFold(buf1.Bytes(), []byte("foo"))
	}
}
