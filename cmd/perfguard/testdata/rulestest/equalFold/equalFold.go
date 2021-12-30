package rulestest

import (
	"bytes"
	"strings"
)

func Warn(s1, s2 string) {
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

func Ignore() {
	// Same lhs and rhs.
	_ = strings.ToLower(x) == x
	_ = x == strings.ToLower(x)
	_ = strings.ToLower(x) != x
	_ = x != strings.ToLower(x)
	_ = strings.ToUpper(x) == x
	_ = x == strings.ToUpper(x)
	_ = strings.ToUpper(x) != x
	_ = x != strings.ToUpper(x)
	_ = bytes.Equal(bytes.ToLower(b), b)
	_ = bytes.Equal(b, bytes.ToLower(b))
	_ = bytes.Equal(bytes.ToUpper(b), b)
	_ = bytes.Equal(b, bytes.ToUpper(b))

	_ = strings.EqualFold(x, y)
	_ = strings.EqualFold(x, concat(y, "123"))
	_ = strings.EqualFold(concat(y, "123"), x)

	_ = bytes.EqualFold(x, y)
	_ = bytes.EqualFold(x, append(y, 'a'))
	_ = bytes.EqualFold(append(y, 'a'), x)

	concat := func(x, y string) string {
		return x + y
	}

	// Side effects.
	_ = strings.ToLower(x) == concat(y, "123")
	_ = bytes.Equal(bytes.ToLower(x), append(y, 'a'))
}
