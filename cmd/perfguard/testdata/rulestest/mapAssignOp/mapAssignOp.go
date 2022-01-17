package rulestest

func Warn(m map[string]int) {
	{
		m["foo"] = m["foo"] + 1 // want `m["foo"] = m["foo"] + 1 => m["foo"]++`
		m["foo"] += 1           // want `m["foo"] += 1 => m["foo"]++`
	}

	{
		var k string
		m[k] = m[k] + 1 // want `m[k] = m[k] + 1 => m[k]++`
		m[k] += 1       // want `m[k] += 1 => m[k]++`

		m[k] = m[k] + 5 // want `m[k] = m[k] + 5 => m[k] += 5`
		m[k] = m[k] - 5 // want `m[k] = m[k] - 5 => m[k] -= 5`
		m[k] = m[k] / 2 // want `m[k] = m[k] / 2 => m[k] /= 2`
		m[k] = m[k] * 4 // want `m[k] = m[k] * 4 => m[k] *= 4`
	}
}

func Ignore(m map[string]int) {
	{
		m["foo"]++
	}

	{
		var k string
		m[k]++
		m[k] = m["foo"] + 1
		m[k] += 5
		m[k] -= 5
		m[k] /= 5
		m[k] *= 15
	}

	{
		m[getString()] = m[getString()] + 1
	}

	{
		type counter struct {
			value int
		}
		m := make(map[int]*counter)
		m[10].value++
		m[10].value += 1
		m[10].value = m[10].value + 1
	}
}

func getString() string { return "" }
