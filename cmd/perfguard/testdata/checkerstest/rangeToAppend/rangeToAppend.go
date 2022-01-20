package rulestest

func Warn() {
	dstBytes := make([]byte, 0, 10)
	var srcBytes []byte

	mapOfBytes := make(map[string][]byte, 100)

	for _, b := range srcBytes { // want `for … { … } => dstBytes = append(dstBytes, srcBytes...)`
		dstBytes = append(dstBytes, b)
	}

	for _, b := range srcBytes { // want `for … { … } => mapOfBytes["k"] = append(mapOfBytes["k"], srcBytes...)`
		mapOfBytes["k"] = append(mapOfBytes["k"], b)
	}

	for _, b := range srcBytes { // want `for … { … } => srcBytes = append(srcBytes, srcBytes...)`
		srcBytes = append(srcBytes, b)
	}

	{
		type object struct {
			bytes []byte
		}
		o := new(object)
		for _, b := range srcBytes { // want `for … { … } => o.bytes = append(o.bytes, srcBytes...)`
			o.bytes = append(o.bytes, b)
		}
	}
}

func Ignore() {
	var dstBytes []byte
	var srcBytes []byte

	mapOfBytes := make(map[string][]byte)

	dstBytes = append(dstBytes, srcBytes...)
	mapOfBytes["k"] = append(mapOfBytes["k"], srcBytes...)

	for _, b := range srcBytes {
		dstBytes = append(dstBytes, b)
		println(b)
	}

	for _, b := range srcBytes {
		println(b)
	}

	srcBytes = append(srcBytes, srcBytes...)

	{
		var arr [10]byte
		for _, b := range arr {
			dstBytes = append(dstBytes, b)
		}
	}
	{
		m := make(map[string]byte)
		for _, b := range m {
			dstBytes = append(dstBytes, b)
		}
	}

	{
		type object struct {
			bytes []byte
		}
		var objects []*object
		for i, b := range srcBytes {
			objects[i].bytes = append(objects[i].bytes, b)
		}
	}

	{
		type Book struct {
			authorID int
		}
		var books []Book
		m := make(map[int][]Book, 10)
		for _, book := range books {
			m[book.authorID] = append(m[book.authorID], book)
		}
	}

	{
		var ifaces []interface{}
		ints := make([]int, 0, len(ifaces))
		for _, x := range ints {
			ifaces = append(ifaces, x)
		}
		_ = ifaces
	}

	_ = srcBytes
	_ = dstBytes
}
