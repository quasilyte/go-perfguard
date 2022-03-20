package rulestest

func Warn() {
	_ = append([]byte{1}, 2)       // want `append([]byte{1}, 2) => []byte{1, 2}`
	_ = append([]byte{1}, 2, 3)    // want `append([]byte{1}, 2, 3) => []byte{1, 2, 3}`
	_ = append([]byte{1}, 2, 3, 4) // want `append([]byte{1}, 2, 3, 4) => []byte{1, 2, 3, 4}`

	_ = append([]byte{}, 1)       // want `append([]byte{}, 1) => []byte{1}`
	_ = append([]byte{}, 1, 2)    // want `append([]byte{}, 1, 2) => []byte{1, 2}`
	_ = append([]byte{}, 1, 2, 3) // want `append([]byte{}, 1, 2, 3) => []byte{1, 2, 3}`

	_ = append([]byte(nil), 1)       // want `append([]byte(nil), 1) => []byte{1}`
	_ = append([]byte(nil), 1, 2)    // want `append([]byte(nil), 1, 2) => []byte{1, 2}`
	_ = append([]byte(nil), 1, 2, 3) // want `append([]byte(nil), 1, 2, 3) => []byte{1, 2, 3}`
}

func Ignore(b []byte) {
	_ = append([]byte{1}, b...)
	_ = append([]byte{1, 2}, 3)
	_ = []byte{1, 2, 3}
}
