package main

import (
	"bytes"
	"strings"
)

func main() {
	type testCase struct {
		s   string
		sub string
	}

	tests := []testCase{
		{"", ""},
		{"", "x"},
		{"x", ""},
		{"x", "x"},
		{"hello, world", ""},
		{"hello, world", "hello"},
		{"hello, world", "world"},
		{"world", "hello, world"},
		{"hello", "hello, world"},
	}

	var allTests []testCase
	for _, test := range tests {
		allTests = append(allTests, test)
		allTests = append(allTests, testCase{
			s:   " " + test.s,
			sub: test.sub,
		})
		allTests = append(allTests, testCase{
			s:   test.s,
			sub: " " + test.sub,
		})
		allTests = append(allTests, testCase{
			s:   " " + test.s,
			sub: " " + test.sub,
		})
		allTests = append(allTests, testCase{
			s:   strings.ToLower(test.s),
			sub: test.sub,
		})
		allTests = append(allTests, testCase{
			s:   test.s,
			sub: strings.ToLower(test.sub),
		})
		allTests = append(allTests, testCase{
			s:   strings.ToLower(test.s),
			sub: strings.ToLower(test.sub),
		})
		allTests = append(allTests, testCase{
			s:   strings.ToUpper(test.s),
			sub: test.sub,
		})
		allTests = append(allTests, testCase{
			s:   test.s,
			sub: strings.ToUpper(test.sub),
		})
		allTests = append(allTests, testCase{
			s:   strings.ToUpper(test.s),
			sub: strings.ToUpper(test.sub),
		})
	}

	for _, test := range allTests {
		println("hassuffix+tolower", test.s, test.sub,
			(len(test.s) >= len(test.sub) && strings.EqualFold(test.s[len(test.s)-len(test.sub):], test.sub)),
		)
		println("hassuffix+toupper", test.s, test.sub,
			(len(test.s) >= len(test.sub) && strings.EqualFold(test.s[len(test.s)-len(test.sub):], test.sub)),
		)
		println("hasprefix+tolower", test.s, test.sub,
			(len(test.s) >= len(test.sub) && strings.EqualFold(test.s[:len(test.sub)], test.sub)),
		)
		println("hasprefix+toupper", test.s, test.sub,
			(len(test.s) >= len(test.sub) && strings.EqualFold(test.s[:len(test.sub)], test.sub)),
		)

		b := []byte(test.s)
		sub := []byte(test.sub)
		println("bytes hassuffix+tolower", test.s, test.sub,
			(len(b) >= len(sub) && bytes.EqualFold(b[len(b)-len(sub):], sub)),
		)
		println("bytes hassuffix+toupper", test.s, test.sub,
			(len(b) >= len(sub) && bytes.EqualFold(b[len(b)-len(sub):], sub)),
		)
		println("bytes hasprefix+tolower", test.s, test.sub,
			(len(b) >= len(sub) && bytes.EqualFold(b[:len(sub)], sub)),
		)
		println("bytes hasprefix+toupper", test.s, test.sub,
			(len(b) >= len(sub) && bytes.EqualFold(b[:len(sub)], sub)),
		)
	}
}
