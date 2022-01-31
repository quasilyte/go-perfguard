package rulestest

import (
	"reflect"
)

func Warn(x interface{}, rv reflect.Value) {
	_ = rv.Type().Kind() // want `rv.Type().Kind() => rv.Kind()`
}

func Ignore(x interface{}, rv reflect.Value) {
	_ = rv.Kind()
}
