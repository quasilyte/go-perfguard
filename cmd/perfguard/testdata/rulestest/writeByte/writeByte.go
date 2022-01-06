package rulestest

import (
	"bytes"
	"strings"
)

func Warn(builder *strings.Builder) {
	var buf bytes.Buffer
	buf.WriteRune('\n') // want `buf.WriteRune('\n') => buf.WriteByte('\n')`

	builder.WriteRune('\n') // want `builder.WriteRune('\n') => builder.WriteByte('\n')`
}

func Ignore(builder *strings.Builder, w RuneWriter) {
	builder.WriteRune('ь')
	builder.WriteByte('\n')

	w.WriteRune('\n')
	w.WriteRune('0')
}

type RuneWriter interface {
	WriteRune(r rune) (int, error)
}
