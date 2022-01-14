package main

import "fmt"

func main() {
	intTests := []int{0, 1, -1, 295}
	for _, v := range intTests {
		println(fmt.Sprintf("%d", v))
		println(fmt.Sprintf("%v", v))
		println(fmt.Sprint(v))
	}

	int64Tests := []int64{0, 1, -1, 295, 45032}
	for _, v := range int64Tests {
		println(fmt.Sprintf("%d", v))
		println(fmt.Sprintf("%v", v))
		println(fmt.Sprintf("%x", v))
		println(fmt.Sprint(v))
	}

	uint64Tests := []uint64{0, 1, 295, 45032}
	for _, v := range uint64Tests {
		println(fmt.Sprintf("%d", v))
		println(fmt.Sprintf("%v", v))
		println(fmt.Sprintf("%x", v))
		println(fmt.Sprint(v))
	}
}
