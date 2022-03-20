package rulestest

import "bytes"

func Warn(b1, b2 []byte) {
	_ = bytes.Compare(b1, b2) == 0 // want `bytes.Compare(b1, b2) == 0 => bytes.Equal(b1, b2)`
	_ = bytes.Compare(b1, b2) != 0 // want `bytes.Compare(b1, b2) != 0 => !bytes.Equal(b1, b2)`
}

func Ignore(b1, b2 []byte) {
	_ = bytes.Equal(b1, b2)
	_ = !bytes.Equal(b1, b2)
}
