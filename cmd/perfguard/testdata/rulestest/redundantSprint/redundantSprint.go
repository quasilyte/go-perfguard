package rulestest

import "fmt"

func Warn() {
	{
		var foo withStringer
		_ = fmt.Sprint(foo)        // want `fmt.Sprint(foo) => foo.String()`
		_ = fmt.Sprintf("%s", foo) // want `fmt.Sprintf("%s", foo) => foo.String()`
		_ = fmt.Sprintf("%v", foo) // want `fmt.Sprintf("%v", foo) => foo.String()`
	}

	{
		var err error
		_ = fmt.Sprint(err)        // want `fmt.Sprint(err) => err.Error()`
		_ = fmt.Sprintf("%s", err) // want `fmt.Sprintf("%s", err) => err.Error()`
		_ = fmt.Sprintf("%v", err) // want `fmt.Sprintf("%v", err) => err.Error()`
	}

	{
		var s string
		_ = fmt.Sprint(s)        // want `fmt.Sprint(s) => s`
		_ = fmt.Sprintf("%s", s) // want `fmt.Sprintf("%s", s) => s`
		_ = fmt.Sprintf("%v", s) // want `fmt.Sprintf("%v", s) => s`

		_ = fmt.Sprint("x")        // want `fmt.Sprint("x") => "x"`
		_ = fmt.Sprintf("%s", "x") // want `fmt.Sprintf("%s", "x") => "x"`
		_ = fmt.Sprintf("%v", "x") // want `fmt.Sprintf("%v", "x") => "x"`
	}

	{
		var b []byte
		_ = fmt.Sprint(b)        // want `fmt.Sprint(b) => string(b)`
		_ = fmt.Sprintf("%s", b) // want `fmt.Sprintf("%s", b) => string(b)`
		_ = fmt.Sprintf("%v", b) // want `fmt.Sprintf("%v", b) => string(b)`
	}
	{
		type myBytes []byte
		var b myBytes
		_ = fmt.Sprint(b)        // want `fmt.Sprint(b) => string(b)`
		_ = fmt.Sprintf("%s", b) // want `fmt.Sprintf("%s", b) => string(b)`
		_ = fmt.Sprintf("%v", b) // want `fmt.Sprintf("%v", b) => string(b)`
	}
	{
		type myString string
		var s myString
		_ = fmt.Sprint(s)        // want `fmt.Sprint(s) => string(s)`
		_ = fmt.Sprintf("%s", s) // want `fmt.Sprintf("%s", s) => string(s)`
		_ = fmt.Sprintf("%v", s) // want `fmt.Sprintf("%v", s) => string(s)`
	}
}

func Ignore() {
	{
		var foo withStringer
		_ = foo.String()
	}

	{
		var err error
		_ = err.Error()
	}

	{
		var s string
		_ = s
		_ = "x"
	}

	{
		var b []byte
		_ = string(b)
	}
	{
		var runes []rune
		_ = string(runes)
	}
	{
		type myBytes []byte
		var b myBytes
		_ = string(b)
	}
	{
		type myString string
		var s myString
		_ = string(s)
	}
	{
		var runes []rune
		_ = fmt.Sprint(runes)
	}
	{
		var runes []int32
		_ = fmt.Sprint(runes)
	}
	{
		var runes []rune
		_ = fmt.Sprint(runes)
		_ = fmt.Sprintf("%s", runes)
		_ = fmt.Sprintf("%v", runes)
	}
}

type withStringer struct{}

func (withStringer) String() string { return "" }
