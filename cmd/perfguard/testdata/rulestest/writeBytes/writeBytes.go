package rulestest

import (
	"bytes"
	"io"
)

type weirdWriter struct{}

func (w weirdWriter) Write([]byte) error       { return nil }
func (w weirdWriter) WriteString(string) error { return nil }

func Warn(src *bytes.Buffer, w io.Writer, dst *bytes.Buffer) {
	io.WriteString(w, src.String())        // want `io.WriteString(w, src.String()) => w.Write(src.Bytes())`
	io.WriteString(w, string(src.Bytes())) // want `io.WriteString(w, string(src.Bytes())) => w.Write(src.Bytes())`
	dst.WriteString(src.String())          // want `dst.WriteString(src.String()) => dst.Write(src.Bytes())`

	{
		var buf bytes.Buffer
		buf.WriteString(src.String()) // want `buf.WriteString(src.String()) => buf.Write(src.Bytes())`
	}

	{
		var b []byte
		var buf bytes.Buffer
		bufPtr := &buf
		buf.WriteString(string(b))    // want `buf.WriteString(string(b)) => buf.Write(b)`
		bufPtr.WriteString(string(b)) // want `bufPtr.WriteString(string(b)) => bufPtr.Write(b)`
	}
}

func Ignore(src *bytes.Buffer, w io.Writer, dst *bytes.Buffer) {
	w.Write(src.Bytes())
	w.Write(src.Bytes())
	dst.Write(src.Bytes())

	{
		var ww weirdWriter
		var b []byte
		ww.WriteString(string(src.Bytes()))
		ww.WriteString(string(b))
	}
}
