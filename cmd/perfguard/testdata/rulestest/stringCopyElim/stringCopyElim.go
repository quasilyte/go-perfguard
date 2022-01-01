package rulestest

func Warn() {
	var b []byte
	var b2 []byte
	var s string

	copy(b, []byte(s)) // want `copy(b, []byte(s)) => copy(b, s)`

	b = append(b, []byte(s)...) // want `append(b, []byte(s)...) => append(b, s...)`

	_ = len(string(b))       // want `len(string(b)) => len(b)`
	_ = len(string(b2)) == 0 // want `len(string(b2)) => len(b2)`
}

func Ignore() {
	var b []byte
	var s string

	copy(b, s)

	{
		copy := func(int) {}
		copy(1)
	}
}
