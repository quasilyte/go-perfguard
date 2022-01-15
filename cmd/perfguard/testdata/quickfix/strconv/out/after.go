package main

import (
	"strconv"
)

func main() {
	intTests := []int{0, 1, -1, 295}
	for _, v := range intTests {
		println(strconv.Itoa(v))
		println(strconv.Itoa(v))
		println(strconv.Itoa(v))
	}

	int64Tests := []int64{0, 1, -1, 295, 45032}
	for _, v := range int64Tests {
		println(strconv.FormatInt(v, 10))
		println(strconv.FormatInt(v, 10))
		println(strconv.FormatInt(v, 16))
		println(strconv.FormatInt(v, 10))
	}

	uint64Tests := []uint64{0, 1, 295, 45032}
	for _, v := range uint64Tests {
		println(strconv.FormatUint(v, 10))
		println(strconv.FormatUint(v, 10))
		println(strconv.FormatUint(v, 16))
		println(strconv.FormatUint(v, 10))
	}
}
