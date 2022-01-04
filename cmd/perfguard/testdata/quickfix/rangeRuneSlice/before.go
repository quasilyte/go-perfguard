package main

func main() {
	tests := []string{
		"",
		"hello",
		"Hello, 世界",
	}
	for _, test := range tests {
		println("f1")
		f1(test)
	}
}

func f1(s string) {
	i := 0
	for _, ch := range []rune(s) {
		i++
		println(i, ch)
	}
}
