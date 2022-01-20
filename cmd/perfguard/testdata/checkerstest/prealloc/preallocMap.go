package checkerstest

func Warn1map(xs []int) map[int]int {
	var squares = make(map[int]int) // want `squares = make(map[int]int) => squares = make(map[int]int, len(xs))`
	for _, x := range xs {
		squares[x] = x * x
	}
	return squares
}

func Warn2map(xs []int) map[int]int {
	// Like other example, but uses an empty map literal instead of a make.
	var squares = map[int]int{} // want `squares = map[int]int{} => squares = make(map[int]int, len(xs))`
	for _, x := range xs {
		squares[x] = x * x
	}
	return squares
}

func Warn3map(xs []int) map[int]int {
	// Like other example, but uses := to define a var.
	squares := make(map[int]int) // want `make(map[int]int) => make(map[int]int, len(xs))`
	for _, x := range xs {
		squares[x] = x * x
	}
	return squares
}

func Warn4map(xs []int) map[int]int {
	// Like other example, but uses := to define a var and a literal instead of a make.
	squares := map[int]int{} // want `map[int]int{} => make(map[int]int, len(xs))`
	for _, x := range xs {
		squares[x] = x * x
	}
	return squares
}
