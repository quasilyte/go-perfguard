package rulestest

import (
	"bytes"
	"strings"
)

func Warn() {
	{
		var x, y string
		_ = strings.ToLower(x) == y                    // want `strings.ToLower(x) == y => strings.EqualFold(x, y)`
		_ = strings.ToLower(x) == strings.ToLower(y)   // want `strings.ToLower(x) == strings.ToLower(y) => strings.EqualFold(x, y)`
		_ = x == strings.ToLower(y)                    // want `x == strings.ToLower(y) => strings.EqualFold(x, y)`
		_ = strings.ToLower(x) != "y"                  // want `strings.ToLower(x) != "y" => !strings.EqualFold(x, "y")`
		_ = strings.ToLower(x) == strings.ToLower("y") // want `strings.ToLower(x) == strings.ToLower("y") => strings.EqualFold(x, "y")`
		_ = x == strings.ToLower("y")                  // want `x == strings.ToLower("y") => strings.EqualFold(x, "y")`
		_ = strings.ToUpper(x) == y                    // want `strings.ToUpper(x) == y => strings.EqualFold(x, y)`
		_ = strings.ToUpper(x) != strings.ToUpper(y)   // want `strings.ToUpper(x) != strings.ToUpper(y) => !strings.EqualFold(x, y)`
		_ = x != strings.ToUpper(y)                    // want `x != strings.ToUpper(y) => !strings.EqualFold(x, y)`
		_ = strings.ToUpper(x) == "y"                  // want `strings.ToUpper(x) == "y" => strings.EqualFold(x, "y")`
		_ = strings.ToUpper(x) == strings.ToUpper("y") // want `strings.ToUpper(x) == strings.ToUpper("y") => strings.EqualFold(x, "y")`
		_ = x == strings.ToUpper("y")                  // want `x == strings.ToUpper("y") => strings.EqualFold(x, "y")`
	}

	{
		var x, y []byte
		_ = bytes.Equal(bytes.ToLower(x), y)                          // want `bytes.Equal(bytes.ToLower(x), y) => bytes.EqualFold(x, y)`
		_ = bytes.Equal(bytes.ToLower(x), bytes.ToLower(y))           // want `bytes.Equal(bytes.ToLower(x), bytes.ToLower(y)) => bytes.EqualFold(x, y)`
		_ = !bytes.Equal(x, bytes.ToLower(y))                         // want `bytes.Equal(x, bytes.ToLower(y)) => bytes.EqualFold(x, y)`
		_ = !bytes.Equal(bytes.ToLower(x), []byte("y"))               // want `bytes.Equal(bytes.ToLower(x), []byte("y")) => bytes.EqualFold(x, []byte("y"))`
		_ = bytes.Equal(bytes.ToLower(x), bytes.ToLower([]byte("y"))) // want `bytes.Equal(bytes.ToLower(x), bytes.ToLower([]byte("y"))) => bytes.EqualFold(x, []byte("y"))`
		_ = bytes.Equal(x, bytes.ToLower([]byte("y")))                // want `bytes.Equal(x, bytes.ToLower([]byte("y"))) => bytes.EqualFold(x, []byte("y"))`
		_ = bytes.Equal(bytes.ToUpper(x), y)                          // want `bytes.Equal(bytes.ToUpper(x), y) => bytes.EqualFold(x, y)`
		_ = !bytes.Equal(bytes.ToUpper(x), bytes.ToUpper(y))          // want `bytes.Equal(bytes.ToUpper(x), bytes.ToUpper(y)) => bytes.EqualFold(x, y)`
		_ = bytes.Equal(x, bytes.ToUpper(y))                          // want `bytes.Equal(x, bytes.ToUpper(y)) => bytes.EqualFold(x, y)`
		_ = bytes.Equal(bytes.ToUpper(x), []byte("y"))                // want `bytes.Equal(bytes.ToUpper(x), []byte("y")) => bytes.EqualFold(x, []byte("y"))`
		_ = bytes.Equal(bytes.ToUpper(x), bytes.ToUpper([]byte("y"))) // want `bytes.Equal(bytes.ToUpper(x), bytes.ToUpper([]byte("y"))) => bytes.EqualFold(x, []byte("y"))`
		_ = bytes.Equal(x, bytes.ToUpper([]byte("y")))                // want `bytes.Equal(x, bytes.ToUpper([]byte("y"))) => bytes.EqualFold(x, []byte("y"))`
	}

	{
		var s1, s2 string
		_ = strings.HasPrefix(strings.ToLower(s1), s2) // want `strings.HasPrefix(strings.ToLower(s1), s2) => (len(s1) >= len(s2) && strings.EqualFold(s1[:len(s2)], s2))`
		_ = strings.HasSuffix(strings.ToLower(s1), s2) // want `strings.HasSuffix(strings.ToLower(s1), s2) => (len(s1) >= len(s2) && strings.EqualFold(s1[len(s1)-len(s2):], s2))`
		_ = strings.HasPrefix(strings.ToUpper(s1), s2) // want `strings.HasPrefix(strings.ToUpper(s1), s2) => (len(s1) >= len(s2) && strings.EqualFold(s1[:len(s2)], s2))`
		_ = strings.HasSuffix(strings.ToUpper(s1), s2) // want `strings.HasSuffix(strings.ToUpper(s1), s2) => (len(s1) >= len(s2) && strings.EqualFold(s1[len(s1)-len(s2):], s2))`
	}
	{
		var b1, b2 []byte
		_ = bytes.HasPrefix(bytes.ToLower(b1), b2) // want `bytes.HasPrefix(bytes.ToLower(b1), b2) => (len(b1) >= len(b2) && bytes.EqualFold(b1[:len(b2)], b2))`
		_ = bytes.HasSuffix(bytes.ToLower(b1), b2) // want `bytes.HasSuffix(bytes.ToLower(b1), b2) => (len(b1) >= len(b2) && bytes.EqualFold(b1[len(b1)-len(b2):], b2))`
		_ = bytes.HasPrefix(bytes.ToUpper(b1), b2) // want `bytes.HasPrefix(bytes.ToUpper(b1), b2) => (len(b1) >= len(b2) && bytes.EqualFold(b1[:len(b2)], b2))`
		_ = bytes.HasSuffix(bytes.ToUpper(b1), b2) // want `bytes.HasSuffix(bytes.ToUpper(b1), b2) => (len(b1) >= len(b2) && bytes.EqualFold(b1[len(b1)-len(b2):], b2))`
	}
}

func Ignore() {
	{
		var x string
		// Same lhs and rhs.
		_ = strings.ToLower(x) == x
		_ = x == strings.ToLower(x)
		_ = strings.ToLower(x) != x
		_ = x != strings.ToLower(x)
		_ = strings.ToUpper(x) == x
		_ = x == strings.ToUpper(x)
		_ = strings.ToUpper(x) != x
		_ = x != strings.ToUpper(x)
	}
	{
		var b []byte
		_ = bytes.Equal(bytes.ToLower(b), b)
		_ = bytes.Equal(b, bytes.ToLower(b))
		_ = bytes.Equal(bytes.ToUpper(b), b)
		_ = bytes.Equal(b, bytes.ToUpper(b))
	}

	concat := func(x, y string) string {
		return x + y
	}

	{
		var x, y string
		_ = strings.EqualFold(x, y)
		_ = strings.EqualFold(x, concat(y, "123"))
		_ = strings.EqualFold(concat(y, "123"), x)
	}

	{
		var x, y []byte
		_ = bytes.EqualFold(x, y)
		_ = bytes.EqualFold(x, append(y, 'a'))
		_ = bytes.EqualFold(append(y, 'a'), x)
	}

	{
		var s1, s2 string
		var b1, b2 []byte
		// Side effects.
		_ = strings.ToLower(s1) == concat(s2, "123")
		_ = bytes.Equal(bytes.ToLower(b1), append(b2, 'a'))
	}

	{
		var s1, s2 string
		_ = (len(s1) >= len(s2) && strings.EqualFold(s1[:len(s2)], s2))
		_ = (len(s1) >= len(s2) && strings.EqualFold(s1[len(s1)-len(s2):], s2))
		_ = (len(s1) >= len(s2) && strings.EqualFold(s1[:len(s2)], s2))
		_ = (len(s1) >= len(s2) && strings.EqualFold(s1[len(s1)-len(s2):], s2))
	}
	{
		var b1, b2 []byte
		_ = (len(b1) >= len(b2) && bytes.EqualFold(b1[:len(b2)], b2))
		_ = (len(b1) >= len(b2) && bytes.EqualFold(b1[len(b1)-len(b2):], b2))
		_ = (len(b1) >= len(b2) && bytes.EqualFold(b1[:len(b2)], b2))
		_ = (len(b1) >= len(b2) && bytes.EqualFold(b1[len(b1)-len(b2):], b2))
	}
}
