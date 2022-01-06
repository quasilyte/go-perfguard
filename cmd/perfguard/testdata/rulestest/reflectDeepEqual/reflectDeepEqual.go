package rulestest

import (
	"bytes"
	"reflect"
)

func Warn() {
	{
		var s1, s2 string
		_ = reflect.DeepEqual(s1, s2)  // want `reflect.DeepEqual(s1, s2) => (s1 == s2)`
		_ = !reflect.DeepEqual(s1, s2) // want `reflect.DeepEqual(s1, s2) => (s1 == s2)`
	}

	{
		var b1, b2 []byte
		_ = reflect.DeepEqual(b1, b2)  // want `reflect.DeepEqual(b1, b2) => bytes.Equal(b1, b2)`
		_ = !reflect.DeepEqual(b1, b2) // want `reflect.DeepEqual(b1, b2) => bytes.Equal(b1, b2)`
	}

	{
		var i1, i2 int
		_ = reflect.DeepEqual(i1, i2)  // want `reflect.DeepEqual(i1, i2) => (i1 == i2)`
		_ = !reflect.DeepEqual(i1, i2) // want `reflect.DeepEqual(i1, i2) => (i1 == i2)`
	}
}

func Ignore() {
	{
		var s1, s2 string
		_ = s1 == s2
		_ = !(s1 == s2)
		_ = s1 != s2
	}

	{
		var b1, b2 []byte
		_ = bytes.Equal(b1, b2)
		_ = !bytes.Equal(b1, b2)
	}

	{
		var i1, i2 int
		_ = i1 == i2
		_ = !(i1 == i2)
		_ = i1 != i2
	}
}
