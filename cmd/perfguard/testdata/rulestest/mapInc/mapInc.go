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
	}

	{
		m[getString()] = m[getString()] + 1
	}
}

func getString() string { return "" }
