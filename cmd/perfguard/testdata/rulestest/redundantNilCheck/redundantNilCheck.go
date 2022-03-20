package rulestest

func Warn() {
	var b []byte
	var ints []int
	var m map[int]int

	_ = b == nil || len(b) == 0       // want `b == nil || len(b) == 0 => len(b) == 0`
	_ = ints == nil || len(ints) == 0 // want `ints == nil || len(ints) == 0 => len(ints) == 0`
	_ = m == nil || len(m) == 0       // want `m == nil || len(m) == 0 => len(m) == 0`
	_ = b == nil || cap(b) == 0       // want `b == nil || cap(b) == 0 => cap(b) == 0`
	_ = ints == nil || cap(ints) == 0 // want `ints == nil || cap(ints) == 0 => cap(ints) == 0`
	_ = len(b) == 0 || b == nil       // want `len(b) == 0 || b == nil => len(b) == 0`
	_ = len(ints) == 0 || ints == nil // want `len(ints) == 0 || ints == nil => len(ints) == 0`
	_ = cap(b) == 0 || b == nil       // want `cap(b) == 0 || b == nil => cap(b) == 0`
	_ = cap(ints) == 0 || ints == nil // want `cap(ints) == 0 || ints == nil => cap(ints) == 0`
	_ = len(m) == 0 || m == nil       // want `len(m) == 0 || m == nil => len(m) == 0`
}

func Ignore() {
	var b []byte
	var ints []int
	var m map[int]int

	_ = len(b) == 0
	_ = cap(b) == 0
	_ = len(ints) == 0
	_ = cap(ints) == 0
	_ = len(m) == 0
}
