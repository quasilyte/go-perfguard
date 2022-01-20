package checkerstest

func Warn1map(xs map[int]int) map[int]int {
	var squares = make(map[int]int) // want `squares = make(map[int]int) => squares = make(map[int]int, len(xs))`
	for k, x := range xs {
		squares[k] = x * x
	}
	return squares
}

func Warn2map(xs map[int]int) map[int]int {
	var squares = map[int]int{} // want `squares = map[int]int{} => squares = make(map[int]int, len(xs))`
	for k, x := range xs {
		squares[k] = x * x
	}
	return squares
}

func Warn3map(xs []int) map[int]int {
	// Like other example, but uses := to define a var.
	squares := make(map[int]int) // want `make(map[int]int) => make(map[int]int, len(xs))`
	for k, x := range xs {
		squares[k] = x * x
	}
	return squares
}

func Warn4map(xs []int) map[int]int {
	// Like other example, but uses := to define a var and a literal instead of a make.
	squares := map[int]int{} // want `map[int]int{} => make(map[int]int, len(xs))`
	for k, x := range xs {
		squares[k] = x * x
	}
	return squares
}

func Warn5map(xs map[string]int) map[string]struct{} {
	// Key-only copying.
	var keys = make(map[string]struct{}) // want `keys = make(map[string]struct{}) => keys = make(map[string]struct{}, len(xs))`
	for k := range xs {
		keys[k] = struct{}{}
	}
	return keys
}

func Ignore1map(xs map[int]int) map[int]int {
	// Copying the values, not the keys.
	var squares = make(map[int]int)
	for _, x := range xs {
		squares[x] = x * x
	}
	return squares
}

func Ignore2map(xs []int) map[int]int {
	// Copying the slice values.
	var squares = make(map[int]int)
	for _, x := range xs {
		squares[x] = x * x
	}
	return squares
}
