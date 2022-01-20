package rulestest

func Warn(srcBytes, dstBytes []byte) {
	{
		for i := range srcBytes { // want `for … { … } => copy(dstBytes, srcBytes)`
			dstBytes[i] = srcBytes[i]
		}
	}

	{
		for i, x := range srcBytes { // want `for … { … } => copy(dstBytes, srcBytes)`
			dstBytes[i] = x
		}
	}

	{
		for i := 0; i < len(srcBytes); i++ { // want `for … { … } => copy(dstBytes, srcBytes)`
			dstBytes[i] = srcBytes[i]
		}
	}

	{
		var xs []int
		var ys []int

		{
			for i := range xs { // want `for … { … } => copy(ys, xs)`
				ys[i] = xs[i]
			}
		}

		{
			for i, x := range xs { // want `for … { … } => copy(ys, xs)`
				ys[i] = x
			}
		}
	}
}

func Ignore(srcBytes, dstBytes []byte) {
	copy(dstBytes, srcBytes)

	{
		var xs []int
		var ys []int

		copy(xs, ys)

		for i := range xs {
			println(i)
			ys[i] = xs[i]
		}
	}

	{
		var src []int
		dst := make(map[int]int, len(src))
		for i, x := range src {
			dst[i] = x
		}
		for i := range src {
			dst[i] = src[i]
		}
	}

	{
		var src []int
		dst := make([]interface{}, len(src))
		for i, x := range src {
			dst[i] = x
		}
	}

	{
		type Book struct {
			authorID int
		}
		var books []Book
		m := make(map[int]Book, 10)
		for _, book := range books {
			m[book.authorID] = book
		}
	}
}
