package rulestest

import (
	"bytes"
	"strings"
)

func Warn(buf *bytes.Buffer, sb *strings.Builder, s string) {
	buf.Write([]byte(s)) // want `buf.Write([]byte(s)) => buf.WriteString(s)`
	sb.Write([]byte(s))  // want `sb.Write([]byte(s)) => sb.WriteString(s)`

	{
		var b bytes.Buffer
		b.Write([]byte(s)) // want `b.Write([]byte(s)) => b.WriteString(s)`
	}
	{
		var sb strings.Builder
		sb.Write([]byte(s)) // want `sb.Write([]byte(s)) => sb.WriteString(s)`
	}
}

func Ignore(buf *bytes.Buffer, sb *strings.Builder, s string) {
	buf.WriteString(s)
	sb.WriteString(s)

	{
		var b bytes.Buffer
		b.WriteString(s)
	}
	{
		var sb strings.Builder
		sb.WriteString(s)
	}
}
