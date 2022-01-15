package main

type withMap struct {
	m map[int]string
}

func main() {
	o := &withMap{
		m: map[int]string{},
	}

	o.m[40] = "40"
	o.m[20] = "20"
	println(len(o.m))
	clearMap(o)
	println(len(o.m))
	o.m[1] = "1"
	println(len(o.m))
	clearMap(o)
	println(len(o.m))
}

func clearMap(o *withMap) {
	o.m = make(map[int]string, len(o.m))
}
