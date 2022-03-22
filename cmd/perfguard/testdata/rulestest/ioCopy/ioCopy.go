package rulestest

import (
	"bytes"
	"io"
	"strings"
)

type myWriter struct{}

func (w *myWriter) Write(b []byte) (int, error)         { return 0, nil }
func (w *myWriter) ReadFrom(r io.Reader) (int64, error) { return 0, nil }

func Warn(buf *bytes.Buffer, sb *strings.Builder, s string) {
	stringsReader := strings.NewReader("123")
	var r io.Reader = stringsReader
	var data []byte

	{
		_, _ = io.Copy(buf, stringsReader) // want `io.Copy(buf, stringsReader) => stringsReader.WriteTo(buf)`
		_, _ = io.Copy(sb, stringsReader)  // want `io.Copy(sb, stringsReader) => stringsReader.WriteTo(sb)`
	}
	{
		var b bytes.Buffer
		_, _ = io.Copy(&b, stringsReader) // want `io.Copy(&b, stringsReader) => stringsReader.WriteTo(&b)`
	}
	{
		var b strings.Builder
		_, _ = io.Copy(&b, stringsReader) // want `io.Copy(&b, stringsReader) => stringsReader.WriteTo(&b)`
	}

	{
		w := &myWriter{}
		_, _ = io.Copy(w, r) // want `io.Copy(w, r) => w.ReadFrom(r)`
	}

	{
		var w myWriter
		_, _ = io.Copy(&w, r) // want `io.Copy(&w, r) => w.ReadFrom(r)`
	}

	{
		io.Copy(buf, bytes.NewReader(data)) // want `io.Copy(buf, bytes.NewReader(data)) => buf.Write(data)`
	}

	{
		io.Copy(buf, strings.NewReader(s)) // want `io.Copy(buf, strings.NewReader(s)) => buf.WriteString(s)`
	}

	{
		_, _ = io.Copy(buf, strings.NewReader(s)) // want `io.Copy(buf, strings.NewReader(s)) => strings.NewReader(s).WriteTo(buf)`
	}
}

func Ignore(buf *bytes.Buffer, sb *strings.Builder, s string) {
	stringsReader := strings.NewReader("123")
	var r io.Reader = stringsReader
	var data []byte

	{
		_, _ = stringsReader.WriteTo(buf)
		_, _ = stringsReader.WriteTo(sb)
	}
	{
		var b bytes.Buffer
		_, _ = stringsReader.WriteTo(&b)
	}
	{
		var b strings.Builder
		_, _ = stringsReader.WriteTo(&b)
	}

	{
		w := &myWriter{}
		_, _ = w.ReadFrom(r)
	}

	{
		var w myWriter
		_, _ = w.ReadFrom(r)
	}

	{
		buf.Write(data)
	}

	{
		buf.WriteString(s)
	}

	{
		_, _ = strings.NewReader(s).WriteTo(buf)
	}
}
