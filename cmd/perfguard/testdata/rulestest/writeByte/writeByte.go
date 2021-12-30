package rulestest

import (
	"bytes"
	"strings"
)

func TODO() {
	// See https://github.com/quasilyte/go-ruleguard/issues/331
	var buf bytes.Buffer
	buf.WriteRune('\n')
}

func Warn(builder *strings.Builder) {
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
