package checkerstest

import (
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
}
