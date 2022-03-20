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

	{
		var x comparable1
		_ = reflect.DeepEqual(x, comparable1{}) // want `reflect.DeepEqual(x, comparable1{}) => (x == comparable1{})`
		_ = reflect.DeepEqual(comparable1{}, x) // want `reflect.DeepEqual(comparable1{}, x) => (comparable1 == x{})`
	}

	{
		var x comparable2
		_ = reflect.DeepEqual(x, comparable2{}) // want `reflect.DeepEqual(x, comparable2{}) => (x == comparable2{})`
		_ = reflect.DeepEqual(comparable2{}, x) // want `reflect.DeepEqual(comparable2{}, x) => (comparable2 == x{})`
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

	{
		var x, y comparable1
		_ = reflect.DeepEqual(x, y)
		_ = reflect.DeepEqual(y, x)
	}

	{
		var x uncomparable
		_ = reflect.DeepEqual(x, uncomparable{})
		_ = reflect.DeepEqual(uncomparable{}, x)
	}
}

type comparable1 struct {
	a string
}

type comparable2 struct {
	comparable1
	x int
	y [4]byte
}

type uncomparable struct {
	_ [0]func()
}
