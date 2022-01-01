package rulestest

func Warn(b1, b2 []byte) {
	type vector2 struct {
		x, y float64
	}

	{
		b1 = append([]byte(nil), b2...) // want `b1 = append([]byte(nil), b2...) => b1 = make([]byte, len(b2)); copy(b1, b2)`

		b1 = append([]byte{}, b2...) // want `b1 = append([]byte{}, b2...) => b1 = make([]byte, len(b2)); copy(b1, b2)`

		dst := append([]byte(nil), b2...) // want `dst := append([]byte(nil), b2...) => dst := make([]byte, len(b2)); copy(dst, b2)`

		dst2 := append([]byte(nil), b2...) // want `dst2 := append([]byte(nil), b2...) => dst2 := make([]byte, len(b2)); copy(dst2, b2)`

		var v1 []vector2
		v2 := append([]vector2(nil), v1...) // want `append([]vector2(nil), v1...) => v2 := make([]vector2, len(v1)); copy(v2, v1)`
		v3 := append([]vector2{}, v2...)    // want `append([]vector2{}, v2...) => v3 := make([]vector2, len(v2)); copy(v3, v2)`

		_, _, _ = v1, v2, v3
		_ = dst
		_ = dst2
	}

	_ = b1
	_ = b2
}

func Ignore(b1, b2 []byte) {
	type vector2 struct {
		x, y float64
	}

	type withPtr struct {
		x *int
	}

	{
		b1 = make([]byte, len(b2))
		copy(b1, b2)

		b1 = make([]byte, len(b2))
		copy(b1, b2)

		dst := make([]byte, len(b2))
		copy(dst, b2)

		dst2 := make([]byte, len(b2))
		copy(dst2, b2)

		var v1 []vector2
		v2 := make([]vector2, len(v1))
		copy(v2, v1)
		v3 := make([]vector2, len(v2))
		copy(v3, v2)

		var p1 []withPtr
		p2 := append([]withPtr(nil), p1...)
		p3 := append([]withPtr{}, p2...)

		_, _, _ = v1, v2, v3
		_, _, _ = p1, p2, p3
		_ = dst
		_ = dst2
	}

	_ = b1
	_ = b2
}
