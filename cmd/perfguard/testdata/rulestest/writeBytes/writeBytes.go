package rulestest

import (
	"bytes"
	"io"
)

func Warn(src *bytes.Buffer, w io.Writer, dst *bytes.Buffer) {
	io.WriteString(w, src.String())        // want `io.WriteString(w, src.String()) => w.Write(src.Bytes())`
	io.WriteString(w, string(src.Bytes())) // want `io.WriteString(w, string(src.Bytes())) => w.Write(src.Bytes())`
	dst.WriteString(src.String())          // want `dst.WriteString(src.String()) => dst.Write(src.Bytes())`
}

func Ignore(src *bytes.Buffer, w io.Writer, dst *bytes.Buffer) {
	w.Write(src.Bytes())
	w.Write(src.Bytes())
	dst.Write(src.Bytes())
}
