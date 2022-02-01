package rulestest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func Warn(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("%+x", 10))) // want `w.Write([]byte(fmt.Sprintf("%+x", 10))) => fmt.Fprintf(w, "%+x", 10)`

	w.Write([]byte(fmt.Sprint(1, 2, 3, 4))) // want `w.Write([]byte(fmt.Sprint(1, 2, 3, 4))) => fmt.Fprint(w, 1, 2, 3, 4)`

	w.Write([]byte(fmt.Sprintln(1, 2, 3, 4))) // want `w.Write([]byte(fmt.Sprintln(1, 2, 3, 4))) => fmt.Fprintln(w, 1, 2, 3, 4)`

	buf := &bytes.Buffer{}

	buf.Write([]byte(fmt.Sprintf("%+x", 10))) // want `buf.Write([]byte(fmt.Sprintf("%+x", 10))) => fmt.Fprintf(buf, "%+x", 10)`

	buf.Write([]byte(fmt.Sprint(1, 2, 3, 4))) // want `buf.Write([]byte(fmt.Sprint(1, 2, 3, 4))) => fmt.Fprint(buf, 1, 2, 3, 4)`

	buf.Write([]byte(fmt.Sprintln(1, 2, 3, 4))) // want `buf.Write([]byte(fmt.Sprintln(1, 2, 3, 4))) => fmt.Fprintln(buf, 1, 2, 3, 4)`

	var i uintptr

	io.WriteString(buf, fmt.Sprint(i)) // want `io.WriteString(buf, fmt.Sprint(i)) => fmt.Fprint(buf, i)`

	io.WriteString(buf, fmt.Sprintf("<%4d>", i)) // want `io.WriteString(buf, fmt.Sprintf("<%4d>", i)) => fmt.Fprintf(buf, "<%4d>", i)`

	io.WriteString(buf, fmt.Sprintln(i, i)) // want `io.WriteString(buf, fmt.Sprintln(i, i)) => fmt.Fprintln(buf, i, i)`

	io.WriteString(os.Stdout, fmt.Sprint(i, i)) // want `io.WriteString(os.Stdout, fmt.Sprint(i, i)) => fmt.Fprint(os.Stdout, i, i)`

	{
		var bw *bufio.Writer
		var key, value string
		if _, err := io.WriteString(bw, fmt.Sprintf("SET %s %s", key, value)); err != nil { // want `io.WriteString(bw, fmt.Sprintf("SET %s %s", key, value)) => fmt.Fprintf(bw, "SET %s %s", key, value)`
			panic(err)
		}
	}
}

func Ignore() {
	{
		var fmt formatter
		var w io.Writer
		w.Write([]byte(fmt.Sprintf()))
		w.Write([]byte(fmt.Sprint()))
		w.Write([]byte(fmt.Sprintln()))

		w.Write([]byte(fmt.Sprintf("%d", 1)))
		w.Write([]byte(fmt.Sprint(1)))
		w.Write([]byte(fmt.Sprintln(1, 2)))
	}
}

type formatter struct{}

func (formatter) Sprintf(args ...interface{}) string  { return "abc" }
func (formatter) Sprint(args ...interface{}) string   { return "abc" }
func (formatter) Sprintln(args ...interface{}) string { return "abc" }
