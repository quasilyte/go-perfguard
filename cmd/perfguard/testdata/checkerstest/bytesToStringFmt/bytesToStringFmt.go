package checkerstest

import (
	"bytes"
	"fmt"
	"io"
)

func Warn(b []byte, w io.Writer) {
	_ = fmt.Sprintf("(%s)", string(b)) // want `string(b) => b`

	_, _ = fmt.Fprintf(w,
		"%s+%d+%s",
		string(b), // want `string(b) => b`
		10,
		string([]byte{'a', 'b'}), // want `string([]byte{'a', 'b'}) => []byte{'a', 'b'}`
	)

	_ = fmt.Sprintf("%d, %s\n", 10, customStringer(b))
	_ = fmt.Sprintf("%s, %d\n", customStringer(b), 10)
}

func Ignore(b []byte, w io.Writer) {
	// Other checkers report this.
	// TODO: #9
	_ = fmt.Sprintf("%s", string(b)) // want `fmt.Sprintf("%s", string(b)) => string(b)`
	_ = fmt.Sprintf("%s", []byte(b)) // want `fmt.Sprintf("%s", []byte(b)) => string([]byte(b))`
	_ = fmt.Sprintf("%s", b)         // want `fmt.Sprintf("%s", b) => string(b)`

	_ = fmt.Sprintf("(%s)", b)

	_, _ = fmt.Fprintf(w,
		"%s+%d+%s",
		b,
		10,
		[]byte{'a', 'b'},
	)

	{
		var out bytes.Buffer
		var e exampleType
		_, _ = fmt.Fprintf(&out, "%d - %d = %s", 1, 2, e.String())
	}
}

type exampleType struct{}

func (e exampleType) String() string { return " " }

func customStringer(data []byte) string { return "" }
