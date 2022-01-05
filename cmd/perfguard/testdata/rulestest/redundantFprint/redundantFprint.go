package rulestest

import (
	"fmt"
	"io"
)

func Warn() {
	{
		var w buffer
		var foo withStringer
		fmt.Fprint(w, foo)        // want `fmt.Fprint(w, foo) => w.WriteString(foo.String())`
		fmt.Fprintf(w, "%s", foo) // want `fmt.Fprintf(w, "%s", foo) => w.WriteString(foo.String())`
		fmt.Fprintf(w, "%v", foo) // want `fmt.Fprintf(w, "%v", foo) => w.WriteString(foo.String())`
	}

	{
		var w buffer
		var err error
		fmt.Fprint(w, err)        // want `fmt.Fprint(w, err) => w.WriteString(err.Error())`
		fmt.Fprintf(w, "%s", err) // want `fmt.Fprintf(w, "%s", err) => w.WriteString(err.Error())`
		fmt.Fprintf(w, "%v", err) // want `fmt.Fprintf(w, "%v", err) => w.WriteString(err.Error())`
	}

	{
		var w buffer
		var s string
		fmt.Fprint(w, s)        // want `fmt.Fprint(w, s) => w.WriteString(s)`
		fmt.Fprintf(w, "%s", s) // want `fmt.Fprintf(w, "%s", s) => w.WriteString(s)`
		fmt.Fprintf(w, "%v", s) // want `fmt.Fprintf(w, "%v", s) => w.WriteString(s)`
	}

	{
		var w buffer
		var b []byte
		fmt.Fprint(w, b)        // want `fmt.Fprint(w, b) => w.Write(b)`
		fmt.Fprintf(w, "%s", b) // want `fmt.Fprintf(w, "%s", b) => w.Write(b)`
		fmt.Fprintf(w, "%v", b) // want `fmt.Fprintf(w, "%v", b) => w.Write(b)`
	}
}

func Ignore() {
	{
		var w io.Writer
		var foo withStringer
		fmt.Fprint(w, foo)
		fmt.Fprintf(w, "%s", foo)
		fmt.Fprintf(w, "%v", foo)
	}

	{
		var w buffer
		var foo withStringer
		w.WriteString(foo.String())
	}

	{
		var w buffer
		var err error
		w.WriteString(err.Error())
	}

	{
		var w buffer
		var s string
		w.WriteString(s)
	}

	{
		var w buffer
		var b []byte
		w.Write(b)
	}
}

type withStringer struct{}

func (withStringer) String() string { return "" }

type buffer struct{}

func (buffer) WriteString(s string) (int, error) { return 0, nil }
func (buffer) Write(b []byte) (int, error)       { return 0, nil }
