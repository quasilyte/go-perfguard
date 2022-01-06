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
			strings.HasSuffix(strings.ToLower(test.s), test.sub),
		)
		println("hassuffix+toupper", test.s, test.sub,
			strings.HasSuffix(strings.ToUpper(test.s), test.sub),
		)
		println("hasprefix+tolower", test.s, test.sub,
			strings.HasPrefix(strings.ToLower(test.s), test.sub),
		)
		println("hasprefix+toupper", test.s, test.sub,
			strings.HasPrefix(strings.ToUpper(test.s), test.sub),
		)

		b := []byte(test.s)
		sub := []byte(test.sub)
		println("bytes hassuffix+tolower", test.s, test.sub,
			bytes.HasSuffix(bytes.ToLower(b), sub),
		)
		println("bytes hassuffix+toupper", test.s, test.sub,
			bytes.HasSuffix(bytes.ToUpper(b), sub),
		)
		println("bytes hasprefix+tolower", test.s, test.sub,
			bytes.HasPrefix(bytes.ToLower(b), sub),
		)
		println("bytes hasprefix+toupper", test.s, test.sub,
			bytes.HasPrefix(bytes.ToUpper(b), sub),
		)
	}
}
