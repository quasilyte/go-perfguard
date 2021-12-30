package rulestest

func Warn() {
	const Zero = 0

	{
		var xs []int
		for i := 0; i < len(xs); i++ { // want `for ... { ... } => for i := range xs { xs[i] = 0 }`
			xs[i] = 0
		}
	}

	{
		var xs []byte
		for i := 0; i < len(xs); i++ { // want `for ... { ... } => for i := range xs { xs[i] = Zero }`
			xs[i] = Zero
		}
	}
}

func Ignore() {
	{
		var xs []int
		for i := range xs {
			xs[i] = 0
		}
	}

	{
		var xs []int
		var ys []int
		for i := 0; i < len(xs); i++ {
			ys[i] = 0
		}
	}
}
