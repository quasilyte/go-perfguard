package main

import (
	"bytes"
	"strings"
)

func main() {
	{
		var buf strings.Builder
		buf.Write([]byte("hello, "))
		buf.Write([]byte("world"))
		println(buf.String())
	}

	{
		var buf strings.Builder
		buf.WriteString("hello, ")
		buf.WriteString("world")
		println(buf.String())
	}

	{
		buf := &bytes.Buffer{}
		buf.Write([]byte("hello, "))
		buf.Write([]byte("world"))
		println(buf.String())
	}

	{
		var buffers [4]bytes.Buffer
		buffers[0].Write([]byte("hello, "))
		buffers[0].Write([]byte("world"))
		println(buffers[0].String())
	}

	{
		var buf strings.Builder
		if _, err := buf.Write([]byte("hello, ")); err != nil {
			panic(err)
		}
		if _, err := buf.Write([]byte("world")); err != nil {
			panic(err)
		}
		println(buf.String())
	}

	{
		var buf strings.Builder
		if _, err := buf.WriteString("hello, "); err != nil {
			panic(err)
		}
		if _, err := buf.WriteString("world"); err != nil {
			panic(err)
		}
		println(buf.String())
	}
}
