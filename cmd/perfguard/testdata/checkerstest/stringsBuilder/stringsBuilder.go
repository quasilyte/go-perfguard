package checkerstest

import (
	"bytes"
	"fmt"
	"io"
)

func Warn1(b []byte) string {
	var buf bytes.Buffer // want `bytes.Buffer => strings.Builder`
	buf.Write(b)
	buf.WriteByte('x')
	buf.WriteRune('Д')
	buf.WriteString("321")
	return buf.String()
}

func Warn2(b []byte) string {
	buf := bytes.Buffer{} // want `bytes.Buffer => strings.Builder`
	buf.Write(b)
	buf.WriteByte('x')
	buf.WriteRune('Д')
	buf.WriteString("321")
	return buf.String()
}

func Warn3(b []byte) string {
	var buf bytes.Buffer // want `bytes.Buffer => strings.Builder`
	writeData(&buf)
	return buf.String()
}

func Warn4(b []byte, s string) string {
	var buf bytes.Buffer // want `bytes.Buffer => strings.Builder`
	fmt.Fprintf(&buf, "%s:%d", s, 124)
	return buf.String()
}

func Warn5(format string, args []interface{}, s string) string {
	buf := bytes.Buffer{} // want `bytes.Buffer => strings.Builder`
	buf.WriteString(fmt.Sprintf(format, args...))
	buf.WriteString(s + s[:2])
	return buf.String()
}

func Warn6() string {
	f := func() string {
		var buf bytes.Buffer // want `bytes.Buffer => strings.Builder`
		buf.WriteString("123")
		return buf.String()
	}
	return f()
}

func writeData(w io.Writer) {
	w.Write([]byte("data"))
}
