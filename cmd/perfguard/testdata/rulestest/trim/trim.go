package trim

import (
	"bytes"
	"strings"
	"unicode"
)

func Warn(s string, b []byte) {
	_ = strings.TrimLeft(strings.TrimRight(s, "+-"), "+-") // want `strings.TrimLeft(strings.TrimRight(s, "+-"), "+-") => strings.Trim(s, "+-")`
	_ = strings.TrimRight(strings.TrimLeft(s, "+-"), "+-") // want `strings.TrimRight(strings.TrimLeft(s, "+-"), "+-") => strings.Trim(s, "+-")`

	_ = bytes.TrimLeft(bytes.TrimRight(b, "+-"), "+-") // want `bytes.TrimLeft(bytes.TrimRight(b, "+-"), "+-") => bytes.Trim(b, "+-")`
	_ = bytes.TrimRight(bytes.TrimLeft(b, "+-"), "+-") // want `bytes.TrimRight(bytes.TrimLeft(b, "+-"), "+-") => bytes.Trim(b, "+-")`

	_ = strings.TrimFunc(s, unicode.IsSpace) // want `strings.TrimFunc(s, unicode.IsSpace) => strings.TrimSpace(s)`
	_ = bytes.TrimFunc(b, unicode.IsSpace)   // want `bytes.TrimFunc(b, unicode.IsSpace) => bytes.TrimSpace(b)`

	_ = strings.Trim(s, " \n\r")       // want `strings.Trim(s, " \n\r") => strings.TrimSpace(s)`
	_ = strings.Trim(s, " \t\r\n")     // want `strings.Trim(s, " \t\r\n") => strings.TrimSpace(s)`
	_ = strings.Trim(s, "\n\r\t ")     // want `strings.Trim(s, "\n\r\t ") => strings.TrimSpace(s)`
	_ = strings.Trim(s, "\n\r\t \v\f") // want `strings.Trim(s, "\n\r\t \v\f") => strings.TrimSpace(s)`
	_ = bytes.Trim(b, " \t\r\n")       // want `bytes.Trim(b, " \t\r\n") => bytes.TrimSpace(b)`
	_ = bytes.Trim(b, "\n\r\t ")       // want `bytes.Trim(b, "\n\r\t ") => bytes.TrimSpace(b)`
	_ = bytes.Trim(b, "\n\r\t \v\f")   // want `bytes.Trim(b, "\n\r\t \v\f") => bytes.TrimSpace(b)`
}

func Ignore() {

}
