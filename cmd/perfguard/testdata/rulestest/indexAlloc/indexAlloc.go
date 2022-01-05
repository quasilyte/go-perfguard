package rulestest

import (
	"bytes"
	"net"
	"strings"
)

func Warn() {
	var b []byte
	var b2 []byte
	var s string
	var s2 string

	_ = strings.Index(string(b), s) // want `strings.Index(string(b), s) => bytes.Index(b, []byte(s))`

	_ = strings.Index(string([]byte("12")), s) // want `strings.Index(string([]byte("12")), s) => bytes.Index([]byte("12"), []byte(s))`

	_ = strings.Index(string([]byte{'1', '2'}), s) // want `strings.Index(string([]byte{'1', '2'}), s) => bytes.Index([]byte{'1', '2'}, []byte(s))`

	_ = strings.Index(string(b), "a"+s) // want `strings.Index(string(b), "a"+s) => bytes.Index(b, []byte("a"+s))`

	_ = strings.Contains(string(b), s)  // want `strings.Contains(string(b), s) => bytes.Contains(b, []byte(s))`
	_ = bytes.Index([]byte(s), b)       // want `bytes.Index([]byte(s), b) => strings.Index(s, string(b))`
	_ = bytes.Contains([]byte(s), b)    // want `bytes.Contains([]byte(s), b) => strings.Contains(s, string(b))`
	_ = strings.HasPrefix(string(b), s) // want `strings.HasPrefix(string(b), s) => bytes.HasPrefix(b, []byte(s))`
	_ = strings.HasSuffix(string(b), s) // want `strings.HasSuffix(string(b), s) => bytes.HasSuffix(b, []byte(s))`
	_ = bytes.HasPrefix([]byte(s), b)   // want `bytes.HasPrefix([]byte(s), b) => strings.HasPrefix(s, string(b))`
	_ = bytes.HasSuffix([]byte(s), b)   // want `bytes.HasSuffix([]byte(s), b) => strings.HasSuffix(s, string(b))`

	{
		var tag []byte
		var tagName string
		_ = strings.Contains(string(tag), tagName+":") // want `strings.Contains(string(tag), tagName+":") => bytes.Contains(tag, []byte(tagName+":"))`
	}

	{
		_ = strings.Contains(string(b), "too many layers of packets") // want `strings.Contains(string(b), "too many layers of packets") => bytes.Contains(b, []byte("too many layers of packets")`
	}

	{
		_ = strings.Contains(string(b), string(b2))  // want `strings.Contains(string(b), string(b2)) => bytes.Contains(b, b2)`
		_ = strings.HasPrefix(string(b), string(b2)) // want `strings.HasPrefix(string(b), string(b2)) => bytes.HasPrefix(b, b2)`
		_ = strings.HasSuffix(string(b), string(b2)) // want `strings.HasSuffix(string(b), string(b2)) => bytes.HasSuffix(b, b2)`
		_ = strings.EqualFold(string(b), string(b2)) // want `strings.EqualFold(string(b), string(b2)) => bytes.EqualFold(b, b2)`
		_ = strings.Compare(string(b), string(b2))   // want `strings.Compare(string(b), string(b2)) => bytes.Compare(b, b2)`

		_ = bytes.Contains([]byte(s), []byte(s2))  // want `bytes.Contains([]byte(s), []byte(s2)) => strings.Contains(s, s2)`
		_ = bytes.HasPrefix([]byte(s), []byte(s2)) // want `bytes.HasPrefix([]byte(s), []byte(s2)) => strings.HasPrefix(s, s2)`
		_ = bytes.HasSuffix([]byte(s), []byte(s2)) // want `bytes.HasSuffix([]byte(s), []byte(s2)) => strings.HasSuffix(s, s2)`
		_ = bytes.EqualFold([]byte(s), []byte(s2)) // want `bytes.EqualFold([]byte(s), []byte(s2)) => strings.EqualFold(s, s2)`
		_ = bytes.Compare([]byte(s), []byte(s2))   // want `bytes.Compare([]byte(s), []byte(s2)) => strings.Compare(s, s2)`
	}
}

func Ignore() {
	var b []byte
	var b1 []byte
	var b2 []byte
	var s string
	var s1 string
	var s2 string

	_ = bytes.Index(b1, b2)
	_ = bytes.Index([]byte(b1), b2)
	_ = bytes.Index(b1, []byte(b2))
	_ = bytes.Index([]byte(b1), []byte(b2))

	_ = bytes.Index(b1, []byte(s1))
	_ = bytes.Index([]byte("12"), []byte(s1))
	_ = bytes.Index([]byte{'1', '2'}, []byte(s1))
	_ = bytes.Index(b1, []byte("a"+s2))

	_ = bytes.Index(getBytes(), []byte(s1))
	_ = bytes.Index(getBytes(), []byte("a"+s1))
	_ = bytes.Index(b1, []byte(getString()))
	_ = bytes.Index([]byte("12"), []byte(getString()))
	_ = bytes.Index([]byte{'1', '2'}, []byte(getString()))
	_ = bytes.Index(b1, []byte("a"+getString()))
	_ = strings.Index(string(getBytes()), s1)
	_ = strings.Index(string(getBytes()), "a"+s1)
	_ = strings.Index(string(b1), getString())
	_ = strings.Index(string(getString()), getString())
	_ = strings.Index(string(b1), "a"+getString())

	{
		var b []byte
		var s string
		_ = bytes.Contains(b, []byte(s))
		_ = strings.Index(s, string(b))
		_ = strings.Contains(s, string(b))
		_ = bytes.HasPrefix(b, []byte(s))
		_ = bytes.HasSuffix(b, []byte(s))
		_ = strings.HasPrefix(s, string(b))
		_ = strings.HasSuffix(s, string(b))
	}

	{
		var tag []byte
		var tagName string
		_ = bytes.Contains(tag, []byte(tagName+":"))
	}
	{
		var tag string
		var tagName string
		_ = strings.Contains(string(tag), tagName+":")
	}

	{
		_ = bytes.Contains(b, b2)
		_ = bytes.HasPrefix(b, b2)
		_ = bytes.HasSuffix(b, b2)
		_ = bytes.EqualFold(b, b2)
		_ = bytes.Compare(b, b2)

		_ = strings.Contains(s, s2)
		_ = strings.HasPrefix(s, s2)
		_ = strings.HasSuffix(s, s2)
		_ = strings.EqualFold(s, s2)
		_ = strings.Compare(s, s2)
	}

	{
		var ip1, ip2 net.IP
		_ = bytes.Compare([]byte(ip1), []byte(ip2))
	}
}

func getBytes() []byte  { return nil }
func getString() string { return "" }
