package rulestest

func Warn() {
	{
		var a1 [30]int
		var a2 [30]int
		copy(a1[:], a2[:]) // want `copy(a1[:], a2[:]) => a1 = a2`
		a1 = a2
	}

	{
		type withArray struct {
			arr [10]byte
		}
		var o1 withArray
		var o2 withArray
		copy(o1.arr[:], o2.arr[:]) // want `copy(o1.arr[:], o2.arr[:]) => o1.arr = o2.arr`
		o1.arr = o2.arr
	}

	{
		type withArray struct {
			arr [10]byte
		}
		o1 := new(withArray)
		o2 := new(withArray)
		copy(o1.arr[:], o2.arr[:]) // want `copy(o1.arr[:], o2.arr[:]) => o1.arr = o2.arr`
		o1.arr = o2.arr
	}
}

func Ignore() {
	{
		var a1 [30]int
		var a2 [30]int
		copy(a1[:5], a2[:5])
		copy(a1[:10], a2[:5])
		copy(a1[:5], a2[:10])
	}

	{
		var a1 [10]byte
		var a2 [12]byte
		copy(a1[:], a2[:])
		copy(a2[:], a1[:])
	}

	{
		a1 := new([10]byte)
		a2 := new([12]byte)
		copy(a1[:], a2[:])
		copy(a2[:], a1[:])
	}

	{
		a1 := new([10]byte)
		a2 := new([12]byte)
		copy((*a1)[:], (*a2)[:])
		copy((*a2)[:], (*a1)[:])
	}
}
