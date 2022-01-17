package rulestest

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

func Warn(b []byte, s string) {
	{
		var w io.Writer
		binary.Write(w, binary.LittleEndian, b) // want `binary.Write(w, binary.LittleEndian, b) => w.Write(b)`
		binary.Write(w, binary.BigEndian, b)    // want `binary.Write(w, binary.BigEndian, b) => w.Write(b)`
	}

	{
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, b) // want `binary.Write(&buf, binary.LittleEndian, b) => buf.Write(b)`
		binary.Write(&buf, binary.BigEndian, b)    // want `binary.Write(&buf, binary.BigEndian, b) => buf.Write(b)`
	}

	{
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, s) // want `binary.Write(&buf, binary.LittleEndian, s) => buf.WriteString(s)`
		binary.Write(&buf, binary.BigEndian, s)    // want `binary.Write(&buf, binary.BigEndian, s) => buf.WriteString(s)`
	}

	{
		buf := &bytes.Buffer{}
		binary.Write(buf, binary.LittleEndian, b) // want `binary.Write(buf, binary.LittleEndian, b) => buf.Write(b)`
		binary.Write(buf, binary.BigEndian, b)    // want `binary.Write(buf, binary.BigEndian, b) => buf.Write(b)`
	}

	{
		var w io.Writer
		if err := binary.Write(w, binary.LittleEndian, b); err != nil { // want `err := binary.Write(w, binary.LittleEndian, b) => _, err := w.Write(b)`
			panic(err)
		}
		if err := binary.Write(w, binary.BigEndian, b); err != nil { // want `err := binary.Write(w, binary.BigEndian, b) => _, err := w.Write(b)`
			panic(err)
		}
	}

	{
		w := &bytes.Buffer{}
		if err := binary.Write(w, binary.LittleEndian, s); err != nil { // want `err := binary.Write(w, binary.LittleEndian, s) => _, err := w.WriteString(s)`
			panic(err)
		}
		if err := binary.Write(w, binary.BigEndian, s); err != nil { // want `err := binary.Write(w, binary.BigEndian, s) => _, err := w.WriteString(s)`
			panic(err)
		}
	}

	{
		var w io.Writer
		err1 := binary.Write(w, binary.LittleEndian, b) // want `err1 := binary.Write(w, binary.LittleEndian, b) => _, err1 := w.Write(b)`
		err2 := binary.Write(w, binary.BigEndian, b)    // want `err2 := binary.Write(w, binary.BigEndian, b) => _, err2 := w.Write(b)`
		if err1 != nil {
			panic(err1)
		}
		if err2 != nil {
			panic(err2)
		}
	}
}

func Ignore(b []byte, s string) {
	{
		var w io.Writer
		w.Write(b)
	}

	{
		var buf bytes.Buffer
		(&buf).Write(b)
	}

	{
		buf := bytes.NewBuffer(make([]byte, 10))
		buf.Write(b)
	}

	{
		var w io.Writer
		if _, err := w.Write(b); err != nil {
			panic(err)
		}
		if _, err := w.Write(b); err != nil {
			panic(err)
		}
	}

	{
		var w io.Writer
		_, err1 := w.Write(b)
		_, err2 := w.Write(b)
		if err1 != nil {
			panic(err1)
		}
		if err2 != nil {
			panic(err2)
		}
	}

	{
		var w io.Writer
		binary.Write(w, binary.LittleEndian, 29)
		binary.Write(w, binary.BigEndian, 29)
	}

	{
		var i int
		buf := bytes.Buffer{}
		binary.Write(&buf, binary.LittleEndian, uint32(i))
		binary.Write(&buf, binary.BigEndian, uint32(i))
	}

	{
		var buf strings.Builder
		buf.WriteString("hello, ")
		buf.WriteString("world")
		println(buf.String())
	}

	{
		var w io.Writer
		if err := binary.Write(w, binary.LittleEndian, s); err != nil {
			panic(err)
		}
		if err := binary.Write(w, binary.BigEndian, s); err != nil {
			panic(err)
		}
	}
}

func checkError(err error) {}
