package rulestest

func Warn() {
	var dstBytes []byte
	var srcBytes []byte

	mapOfBytes := make(map[string][]byte)

	for _, b := range srcBytes { // want `for ... { ... } => dstBytes = append(dstBytes, srcBytes...)`
		dstBytes = append(dstBytes, b)
	}

	for _, b := range srcBytes { // want `for ... { ... } => mapOfBytes["k"] = append(mapOfBytes["k"], srcBytes...)`
		mapOfBytes["k"] = append(mapOfBytes["k"], b)
	}

	for _, b := range srcBytes { // want `for ... { ... } => srcBytes = append(srcBytes, srcBytes...)`
		srcBytes = append(srcBytes, b)
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

	_ = srcBytes
	_ = dstBytes
}
