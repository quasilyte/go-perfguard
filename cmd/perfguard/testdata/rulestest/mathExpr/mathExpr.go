package rulestest

import (
	"math"
)

func Warn(x, y float64) {
	_ = math.Abs(x) * math.Abs(y) // want `math.Abs(x) * math.Abs(y) => math.Abs((x) * (y))`
	_ = math.Abs(x) / math.Abs(y) // want `math.Abs(x) / math.Abs(y) => math.Abs((x) / (y))`
}

func Ignore(x, y float64) {
	_ = math.Abs(x * y)
	_ = math.Abs(x / y)

	_ = math.Abs(x) + math.Abs(y)
	_ = math.Abs(x) - math.Abs(y)
}
