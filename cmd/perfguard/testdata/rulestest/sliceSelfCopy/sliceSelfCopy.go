package rulestest

func Warn(b []byte) {
	{
		keep := 10
		from := 20
		for i := 0; i < keep; i++ { // want `for ... { ... } => copy(b[:keep], b[from:])`
			b[i] = b[from+i]
		}
	}
}

func Ignore(b []byte, a [100]int) {
	{
		keep := 10
		from := 20
		copy(b[:keep], b[from:])
	}

	{
		keep := 10
		from := 20
		for i := 0; i < keep; i++ {
			a[i] = a[from+i]
		}
	}
}
