package checkerstest

func Warn1(xs []int) []int {
	if len(xs) == 0 {
		return nil
	}
	var ys = []int{} // want `ys = []int{} => ys = make([]int, 0, len(xs))`
	for _, x := range xs {
		ys = append(ys, x+1)
	}
	if len(xs) < 10 {
		var (
			ys2 = []int{}        // want `ys2 = []int{} => ys2 = make([]int, 0, len(xs))`
			zs2 = make([]int, 0) // want `zs2 = make([]int, 0) => zs2 = make([]int, 0, len(xs))`
		)
		for _, x := range xs {
			println(x)
			ys2 = append(ys2, x+1)
			zs2 = append(zs2, 1)
		}
		return append(ys2, zs2...)
	}
	return ys
}

func Warn2(xs []int) ([]int, []int) {
	var ys []int // want `can use len(xs) as make size hint for ys`
	var zs []int // want `can use len(xs) as make size hint for zs`
	for _, x := range xs {
		ys = append(ys, x+1)
		zs = append(zs, 1)
	}
	return ys, zs
}

func Warn3(xs []int) []int {
	ys := []int{} // want `[]int{} => make([]int, 0, len(xs))`
	for _, x := range xs {
		ys = append(ys, x+1)
	}
	return ys
}

func Warn4(xs map[int]struct{}) []int {
	ys := []int{} // want `[]int{} => make([]int, 0, len(xs))`
	for x := range xs {
		ys = append(ys, x+1)
	}
	return ys
}

func Warn5(xs []map[int]struct{}) []int {
	ys := []int{} // want `[]int{} => make([]int, 0, len(xs[0]))`
	for range xs[0] {
		ys = append(ys, 1)
	}
	return ys
}

func Warn6(xs []map[int]struct{}) []int {
	var ys []int = []int{} // want `ys []int = []int{} => ys []int = make([]int, 0, len(xs[0]))`
	for range xs[0] {
		ys = append(ys, 1)
	}
	return ys
}

func Ignore1(xs []int) []int {
	var ys []int
	if len(xs) == 0 {
		return ys // Untrack ys: returned before appended to
	}
	for _, x := range xs {
		ys = append(ys, x+1)
	}
	return ys
}

func Ignore2(xs []int) []int {
	// Conditional append.
	var ys []int
	for _, x := range xs {
		if x%2 == 0 {
			ys = append(ys, x)
		}
	}
	return ys
}

func Ignore3(xs []int) []int {
	// Conditional range with append.
	var ys []int
	if len(xs) > 10 {
		for _, x := range xs {
			ys = append(ys, x*2)
		}
	}
	return ys
}

func Ignore4() []int {
	// Range expression has side effects.
	var ys []int
	for _, x := range getInts() {
		ys = append(ys, x+1)
	}
	return ys
}

func Ignore5(xs []int) []int {
	// Has continue, so it's effectively a conditional append.
	var ys []int
	for _, x := range xs {
		if x == 0 {
			continue
		}
		ys = append(ys, x+1)
	}
	return ys
}

func Ignore6(xs []int) []int {
	var ys []int
	println(len(ys)) // Untrack ys: used before appended to
	for _, x := range xs {
		ys = append(ys, x+1)
	}
	return ys
}

func Ignore7(xs []int) []int {
	// Has break.
	var ys []int
	for _, x := range xs {
		if x == 0 {
			break
		}
		ys = append(ys, x+1)
	}
	return ys
}

func Ignore8(xs []int) []int {
	// Has break.
	var ys []int
	for _, x := range xs {
		ys = append(ys, x+1)
		break
	}
	return ys
}

func getInts() []int { return []int{42} }
