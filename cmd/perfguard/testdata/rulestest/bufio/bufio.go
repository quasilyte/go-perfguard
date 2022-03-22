package rulestest

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func read(r io.Reader)                 {}
func readFromBuffered(r *bufio.Reader) {}

func Warn(buf *bytes.Buffer, s string, b []byte) io.Reader {
	{
		read(bufio.NewReader(bytes.NewReader(b)))   // want `bufio.NewReader(bytes.NewReader(b)) => bytes.NewReader(b)`
		read(bufio.NewReader(strings.NewReader(s))) // want `bufio.NewReader(strings.NewReader(s)) => strings.NewReader(s)`
	}
	{
		var br *bytes.Reader
		var sr *strings.Reader
		read(bufio.NewReader(br)) // want `bufio.NewReader(br) => br`
		read(bufio.NewReader(sr)) // want `bufio.NewReader(sr) => sr`
	}
	if buf != nil {
		return bufio.NewReader(bytes.NewReader(b)) // want `bufio.NewReader(bytes.NewReader(b)) => bytes.NewReader(b)`
	}
	var br *bytes.Reader
	_ = []io.Reader{
		bufio.NewReader(br), // want `bufio.NewReader(br) => br`
	}
	{
		read(bufio.NewReader(buf)) // want `bufio.NewReader(buf) => buf`
	}
	return nil
}

func Ignore(buf *bytes.Buffer, s string, b []byte) io.Reader {
	{
		read(bytes.NewReader(b))
		read(strings.NewReader(s))
	}
	{
		var br *bytes.Reader
		var sr *strings.Reader
		read(br)
		read(sr)
	}
	if buf != nil {
		return bytes.NewReader(b)
	}
	var br *bytes.Reader
	var r io.Reader
	_ = []io.Reader{
		br,
		r,
		bufio.NewReader(r),
	}

	{
		var buffered *bufio.Reader = bufio.NewReader(bytes.NewReader(b))
		_ = buffered
	}
	readFromBuffered(bufio.NewReader(bytes.NewReader(b)))
	readFromBuffered(bufio.NewReader(strings.NewReader(s)))

	read(buf)

	return nil
}
