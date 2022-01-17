package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

func main() {
	{
		var buf strings.Builder
		binary.Write(&buf, binary.LittleEndian, []byte("hello, "))
		binary.Write(&buf, binary.BigEndian, []byte("world"))
		println(buf.String())
	}

	{
		var buf strings.Builder
		binary.Write(&buf, binary.LittleEndian, "hello, ")
		binary.Write(&buf, binary.BigEndian, "world")
		println(buf.String())
	}

	{
		buf := bytes.NewBuffer(make([]byte, 10))
		binary.Write(buf, binary.LittleEndian, []byte("hello, "))
		binary.Write(buf, binary.BigEndian, []byte("world"))
		println(buf.String())
	}

	{
		var buffers [4]bytes.Buffer
		binary.Write(&buffers[0], binary.LittleEndian, []byte("hello, "))
		binary.Write(&buffers[0], binary.BigEndian, []byte("world"))
		println(buffers[0].String())
	}

	{
		var buf strings.Builder
		if err := binary.Write(&buf, binary.LittleEndian, []byte("hello, ")); err != nil {
			panic(err)
		}
		if err := binary.Write(&buf, binary.BigEndian, []byte("world")); err != nil {
			panic(err)
		}
		println(buf.String())
	}

	{
		var buf strings.Builder
		if err := binary.Write(&buf, binary.LittleEndian, "hello, "); err != nil {
			panic(err)
		}
		if err := binary.Write(&buf, binary.BigEndian, "world"); err != nil {
			panic(err)
		}
		println(buf.String())
	}
}
