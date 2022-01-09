package main

import (
	"bytes"
)

func main() {
	{
		var buf bytes.Buffer
		buf.Write([]byte("hello, "))
		buf.Write([]byte("world"))
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
		var buf bytes.Buffer
		if _, err := buf.Write([]byte("hello, ")); err != nil {
			panic(err)
		}
		if _, err := buf.Write([]byte("world")); err != nil {
			panic(err)
		}
		println(buf.String())
	}
}
