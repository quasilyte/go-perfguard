package main

import (
	"bytes"
	"encoding/binary"
)

func main() {
	{
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, []byte("hello, "))
		binary.Write(&buf, binary.BigEndian, []byte("world"))
		println(buf.String())
	}

	{
		buf := &bytes.Buffer{}
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
		var buf bytes.Buffer
		if err := binary.Write(&buf, binary.LittleEndian, []byte("hello, ")); err != nil {
			panic(err)
		}
		if err := binary.Write(&buf, binary.BigEndian, []byte("world")); err != nil {
			panic(err)
		}
		println(buf.String())
	}
}
