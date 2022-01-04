package main

func main() {
	tests := []string{
		"",
		"hello",
		"Hello, 世界",
	}
	for _, test := range tests {
		f1(test)
		f2(test)
		f3([]rune(test))
		f4(test, "ab")
	}
}

func f1(s string) {
	i := 0
	for _, ch := range s {
		i++
		println("f1", i, ch)
	}
}

func f2(s string) {
	var ch rune
	for _, ch = range s {
		println("f2", ch)
		break
	}
	println(ch)
}

func f3(runes []rune) {
	var ch rune
	for _, ch = range runes {
		println("f3", ch)
	}
	println(ch)
}

func f4(s1, s2 string) {
	for _, ch1 := range s1 {
		for _, ch2 := range s2 {
			println("f4", ch1, ch2)
		}
	}
}
