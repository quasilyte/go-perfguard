package rulestest

func Warn(m map[string]int) {
	{
		m = make(map[string]int, len(m)) // want `m = make(map[string]int, len(m)) => for k := range m { delete(m, k) }`
	}

	{
		type withMap struct {
			m map[int]bool
		}
		o := &withMap{}
		o.m = make(map[int]bool, len(o.m)) // want `o.m = make(map[int]bool, len(o.m)) => for k := range o.m { delete(o.m, k) }`
	}
}

func Ignore(m map[string]int) {
	{
		for k := range m {
			delete(m, k)
		}
	}

	{
		type withMap struct {
			m map[int]bool
		}
		o := &withMap{}
		for k := range o.m {
			delete(o.m, k)
		}
	}
}
