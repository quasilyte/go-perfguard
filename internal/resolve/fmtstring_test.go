package resolve

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFmtString(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{"", []string{}},
		{"%%", []string{}},
		{"%%d", []string{}},

		{"%", nil},
		{"%[0]d", nil},
		{"%z", nil},

		{"%%%d", []string{"%d<0>"}},
		{"%%%+d", []string{"%+d<0>"}},
		{"%%%-d", []string{"%-d<0>"}},
		{"%d%d", []string{"%d<0>", "%d<1>"}},
		{"%s%d", []string{"%s<0>", "%d<1>"}},
	}

	resolveArgs := func(s string) []string {
		info, ok := FmtString(s)
		if !ok {
			return nil
		}
		args := make([]string, len(info.Args))
		for i, a := range info.Args {
			args[i] = a.String()
		}
		return args
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			have := resolveArgs(test.s)
			if have == nil && test.want != nil {
				t.Fatalf("%q: unexpected nil, want %q", test.s, test.want)
			}
			if !reflect.DeepEqual(have, test.want) {
				t.Fatalf("%q: results mismatched:\nhave: %q\nwant: %q",
					test.s, have, test.want)
			}
		})
	}
}
