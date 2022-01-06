package rulestest

func Warn(s string, b []byte) {
	_ = string(b)[:5] // want `string(b)[:5] => string(b[:5])`
	_ = []byte(s)[:5] // want `[]byte(s)[:5] => []byte(s[:5])`
}

func Ignore(s string, b []byte) {
	_ = string(b[:5])
	_ = []byte(s[:5])
}
