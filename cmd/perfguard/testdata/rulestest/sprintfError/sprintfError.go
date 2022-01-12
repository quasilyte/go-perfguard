package rulestest

import (
	"errors"
	"fmt"
)

func Warn() {
	{
		var filename string
		var line int
		_ = errors.New(fmt.Sprintf("%s:%d", filename, line)) // want `errors.New(fmt.Sprintf("%s:%d", filename, line)) => fmt.Errorf("%s:%d", filename, line)`
	}
}

func Ignore() {
	{
		var filename string
		var line int
		_ = fmt.Errorf("%s:%d", filename, line)
	}
}
