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

// Comment1
import "io"
// Comment2
import "bytes"
// Comment3

func f() string { return "" }`, `
package example

// Comment1
// Comment2
// Comment3

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
package example//OK

import (
	"io"
)
import (
	"bytes"
)

func f() string { return "" }`, `
package example//OK


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

		// Add package from non-std list.
		{`
package example // OK

func f() error {
	return errors.New("ok")
}`, `
package example 
import (
	"github.com/pkg/errors"
)
// OK

func f() error {
	return errors.New("ok")
}`,
		},

		// Remove package from non-std list.
		{`
package example

import (
	"github.com/pkg/errors"
)

func f() string {
	return ""
}`, `
package example


func f() string {
	return ""
}`,
		},

		// Remove renamed package import.
		{`
package example

// This is a comment.
import (
	xerrors "github.com/pkg/errors"
)

func f() string {
	return ""
}`, `
package example

// This is a comment.

func f() string {
	return ""
}`,
		},

		// Do not remove package that is used via local name.
		{`
package example

import (
	pkgerrors "github.com/pkg/errors"
)

func f() error {
	return pkgerrors.New("ok")
}`, `
package example

import (
	pkgerrors "github.com/pkg/errors"
)

func f() error {
	return pkgerrors.New("ok")
}`,
		},

		// Do not remove "_" packages.
		{`
package example

import (
	_ "image/png"
)
import (
	"bytes"
)

func f() int { return rand.Intn(10) }`, `
package example

import (
	_ "image/png"
	"math/rand"
)

func f() int { return rand.Intn(10) }`,
		},

		// No changes in imports: 0 imports.
		{`
package example

func f() int { return 10 }`, `
package example

func f() int { return 10 }`,
		},

		// No changes in imports: 1 import.
		{`
package example

import "io"

func f(r io.Reader) int { return 10 }`, `
package example

import "io"

func f(r io.Reader) int { return 10 }`,
		},

		// No changes in imports: 1 import.
		{`
package example

import (
	"io"
)

func f(r io.Reader) int { return 10 }`, `
package example

import (
	"io"
)

func f(r io.Reader) int { return 10 }`,
		},

		// No changes in imports: 2 imports.
		{`
package example

import (
	"io"
	"strings"
)

func f(r io.Reader) string { return strings.Repeat("x", 10) }`, `
package example

import (
	"io"
	"strings"
)

func f(r io.Reader) string { return strings.Repeat("x", 10) }`,
		},

		// Add second entry due to the first one being renamed.
		{`
package example

import (
	"io"
	str "strings"
)

func f(r io.Reader) string { return str.Repeat("x", 10) }
func f2(r io.Reader) string { return strings.Repeat("x", 10) }`, `
package example

import (
	"io"
	str "strings"
	"strings"
)

func f(r io.Reader) string { return str.Repeat("x", 10) }
func f2(r io.Reader) string { return strings.Repeat("x", 10) }`,
		},

		// Try not to discard the imports grouping.
		{`
package example

import (
	"io"

	"github.com/quasilyte/go-perfguard/otherlib"

	"github.com/quasilyte/lib"
)

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`, `
package example

import (
	"io"
	"strings"

	"github.com/quasilyte/go-perfguard/otherlib"

	"github.com/quasilyte/lib"
)

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`,
		},

		// Try not to discard the imports grouping.
		{`
package example

import (
	"io"
)

import (
	"github.com/quasilyte/go-perfguard/otherlib"

	"github.com/quasilyte/lib"
)

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`, `
package example

import (
	"io"
	"strings"
)

import (
	"github.com/quasilyte/go-perfguard/otherlib"

	"github.com/quasilyte/lib"
)

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`,
		},

		// Try not to add extra () for unchanged import decls.
		{`
package example

import (
	"io"
)

import "github.com/quasilyte/go-perfguard/otherlib"

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`, `
package example

import (
	"io"
	"strings"
)

import "github.com/quasilyte/go-perfguard/otherlib"

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`,
		},

		// Do not inject anything to standalone import C decls.
		{`
package example

// Some C-related comment.
// Should not be removed.
import "C"

import (
	"io"
)

func f(r io.Reader, b *strings.Builder) {}`, `
package example

// Some C-related comment.
// Should not be removed.
import "C"

import (
	"io"
	"strings"
)

func f(r io.Reader, b *strings.Builder) {}`,
		},

		// Try not to lose associated comments.
		{`
package example

// Comment 1.
import "github.com/quasilyte/go-perfguard/otherlib"

// Comment 2.
import (
	// Comment 3.
	"io" // Comment 4.
)

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`, `
package example

// Comment 1.
import (
	"github.com/quasilyte/go-perfguard/otherlib"
	"strings"
)

// Comment 2.
import (
	// Comment 3.
	"io" // Comment 4.
)

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`,
		},

		// Try not to lose associated comments.
		{`
package example

// Comment 1.
import "github.com/quasilyte/go-perfguard/otherlib"

// Comment 2.
// Comment 3.
import (
	/* Comment 4. */
	"io" /* Comment 5. */
)

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`, `
package example

// Comment 1.
import (
	"github.com/quasilyte/go-perfguard/otherlib"
	"strings"
)

// Comment 2.
// Comment 3.
import (
	/* Comment 4. */
	"io" /* Comment 5. */
)

func f(x otherlib.T, y lib.T, r io.Reader, b *strings.Builder) {}`,
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
					"errors":  "errors",
					"io":      "io",
				},
				Packages: map[string]string{
					"errors": "github.com/pkg/errors",
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
