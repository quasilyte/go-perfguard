package imports

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFix(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		// Add a single import to a file without imports.
		{`
package example

func f() string {
	return strings.Repeat("foo", 2)
}`, `
package example

import (
	"strings"
)

func f() string {
	return strings.Repeat("foo", 2)
}`,
		},

		// Add two imports to a file without imports.
		{`
package example

func f() string {
	return strings.Repeat("foo", 2)
}
func g() []byte {
	return bytes.Repeat([]byte("foo"), 2)
}`, `
package example

import (
	"bytes"
	"strings"
)

func f() string {
	return strings.Repeat("foo", 2)
}
func g() []byte {
	return bytes.Repeat([]byte("foo"), 2)
}`,
		},

		// Add two imports to a file with imports in () group.
		{`
package example

import (
	"io"
)

func f(r io.Reader) (string, []byte) {
	return strings.Repeat("foo", 2),
           bytes.Repeat(nil, 2)
}`, `
package example

import (
	"io"
	"bytes"
	"strings"
)

func f(r io.Reader) (string, []byte) {
	return strings.Repeat("foo", 2),
           bytes.Repeat(nil, 2)
}`,
		},

		// Add two imports to a file with 1 unit import.
		{`
package example

import "io"

func f(r io.Reader) (string, []byte) {
	return strings.Repeat("foo", 2),
           bytes.Repeat(nil, 2)
}`, `
package example

import (
	"io"
	"bytes"
	"strings"
)

func f(r io.Reader) (string, []byte) {
	return strings.Repeat("foo", 2),
           bytes.Repeat(nil, 2)
}`,
		},

		// Add one import to a file with 1 unit import.
		{`
package example

import "io"

func f(r io.Reader) string {
	return strings.Repeat("foo", 2)
}`, `
package example

import (
	"io"
	"strings"
)

func f(r io.Reader) string {
	return strings.Repeat("foo", 2)
}`,
		},

		// Add and remove one.
		{`
package example

import "io"

func f() string { return strings.Repeat("foo", 2)
}`, `
package example

import (
	"strings"
)

func f() string { return strings.Repeat("foo", 2)
}`,
		},

		// Removed all imports.
		{`
package example

import "io"

func f() string { return "" }`, `
package example


func f() string { return "" }`,
		},

		// Removed all imports.
		{`
package example

import "io"

import "bytes"

func f() string { return "" }`, `
package example



func f() string { return "" }`,
		},

		// Removed all imports.
		{`
package example

import "io"
import "bytes"

func f() string { return "" }`, `
package example


func f() string { return "" }`,
		},

		// Removed all imports.
		{`
package example

import ("io")
import ("bytes")

func f() string { return "" }`, `
package example


func f() string { return "" }`,
		},

		// Removed all imports.
		{`
package example

import (
	"io"
)
import (
	"bytes"
)

func f() string { return "" }`, `
package example


func f() string { return "" }`,
		},

		// Removed two, added one.
		{`
package example

import (
	"io"
)
import (
	"bytes"
)

func f() int { return rand.Intn(10) }`, `
package example

import (
	"math/rand"
)

func f() int { return rand.Intn(10) }`,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			config := FixConfig{
				StdlibPackages: map[string]string{
					"strings": "strings",
					"bytes":   "bytes",
					"rand":    "math/rand",
				},
			}
			fixed, err := Fix(config, []byte(test.before))
			if err != nil {
				t.Fatal(err)
			}
			have := string(bytes.TrimSpace(fixed))
			want := strings.TrimSpace(test.after)
			if diff := cmp.Diff(have, want); diff != "" {
				t.Errorf("results mismatches (+want -have):\n%s", diff)
				fmt.Println("Before:")
				fmt.Println(test.before)
				fmt.Println("After:")
				fmt.Println(have)
			}
		})
	}
}
