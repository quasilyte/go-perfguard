package checkerstest

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func Warn1(b []byte) string {
	var buf bytes.Buffer // want `bytes.Buffer => strings.Builder`
	buf.Write(b)
	buf.WriteByte('x')
	buf.WriteRune('Д')
	buf.WriteString("321")
	return buf.String()
}

func Fixed1(b []byte) string {
	var buf strings.Builder
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

func Fixed2(b []byte) string {
	buf := strings.Builder{}
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

func Fixed3(b []byte) string {
	var buf strings.Builder
	writeData(&buf)
	return buf.String()
}

func Warn4(b []byte, s string) string {
	var buf bytes.Buffer // want `bytes.Buffer => strings.Builder`
	fmt.Fprintf(&buf, "%s:%d", s, 124)
	return buf.String()
}

func Fixed4(b []byte, s string) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "%s:%d", s, 124)
	return buf.String()
}

func Warn5(format string, args []interface{}, s string) string {
	buf := bytes.Buffer{} // want `bytes.Buffer => strings.Builder`
	buf.WriteString(fmt.Sprintf(format, args...))
	buf.WriteString(s + s[:2])
	return buf.String()
}

func Fixed5(format string, args []interface{}, s string) string {
	buf := strings.Builder{}
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

func Fixed6() string {
	f := func() string {
		var buf strings.Builder
		buf.WriteString("123")
		return buf.String()
	}
	return f()
}

func Warn7(cond bool) string {
	buf := &bytes.Buffer{} // want `bytes.Buffer => strings.Builder`
	if cond {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	return buf.String()
}

func Fixed7(cond bool) string {
	buf := &strings.Builder{}
	if cond {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	return buf.String()
}

func Warn8(cond bool) string {
	buf := &bytes.Buffer{} // want `bytes.Buffer => strings.Builder`
	buf.WriteString("ok")
	res := buf.String()
	return res
}

func Fixed8(cond bool) string {
	buf := &strings.Builder{}
	buf.WriteString("ok")
	res := buf.String()
	return res
}

func Ignore1() string {
	// buf.Reset() is used.
	var buf bytes.Buffer
	buf.WriteString("foo")
	res := buf.String()
	buf.Reset()
	return res
}

func Ignore2() []byte {
	// buf.Bytes() is used.
	var buf bytes.Buffer
	buf.WriteString("foo")
	return buf.Bytes()
}

func Ignore3() {
	// buf.String() is never used.
	var buf bytes.Buffer
	buf.WriteString("foo")
	writeData(&buf)
}

func Ignore4() bytes.Buffer {
	// buf is returned.
	var buf bytes.Buffer
	buf.WriteString("foo")
	_ = buf.String()
	return buf
}

func Ignore5() *bytes.Buffer {
	// buf is returned.
	var buf bytes.Buffer
	buf.WriteString("foo")
	_ = buf.String()
	return &buf
}

func Ignore6() string {
	// buf is passed as non-interface argument.
	var buf bytes.Buffer
	buf.WriteString("foo")
	writeDataBuf(&buf)
	return buf.String()
}

func writeData(w io.Writer) {
	w.Write([]byte("data"))
}

func writeDataBuf(w *bytes.Buffer) {
	w.WriteString("data")
}
