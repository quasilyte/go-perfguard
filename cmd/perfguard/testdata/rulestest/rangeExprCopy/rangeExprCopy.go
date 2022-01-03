package rulestest

var globalArray [400]object

type object struct {
	x, y, z int64
	data    []byte
}

func getArray() [50000]byte {
	return [50000]byte{}
}

func Warn() {
	var localArray [1000]int
	for _, v := range localArray { // want `localArray => &localArray`
		println(v)
	}

	for _, v := range globalArray { // want `globalArray => &globalArray`
		println(v.x)
		println(v.y)
	}

	for _, b := range getArray() { // want `range over big array value expression is ineffective`
		println(b)
	}
}

func Ignore() {
	var localArray [1000]int
	for _, v := range &localArray {
		println(v)
	}

	localSlice := localArray[:]
	for _, v := range localSlice {
		println(v)
	}

	for i := range localArray {
		println(localArray[i])
	}

	for i := range localArray {
		v := &localArray[i]
		println(v)
	}

	for range localArray {
	}

	arrayPtr := &localArray
	for _, v := range arrayPtr {
		println(v)
	}
}
