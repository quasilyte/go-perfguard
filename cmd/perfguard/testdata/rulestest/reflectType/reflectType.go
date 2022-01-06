package rulestest

import (
	"fmt"
	"reflect"
)

func Warn(x interface{}) {
	_ = reflect.ValueOf(x).Type() // want `reflect.ValueOf(x).Type() => reflect.TypeOf(x)`
	_ = fmt.Sprintf("%T", x)      // want `fmt.Sprintf("%T", x) => reflect.TypeOf(x).String()`

	rv := reflect.ValueOf(x)
	_ = reflect.TypeOf(rv.Interface())    // want `reflect.TypeOf(rv.Interface()) => rv.Type()`
	_ = fmt.Sprintf("%T", rv.Interface()) // want `fmt.Sprintf("%T", rv.Interface()) => rv.Type().String()`

	_ = reflect.ValueOf(x).Type().Size() // want `reflect.ValueOf(x).Type() => reflect.TypeOf(x)`
}

func Ignore(x interface{}) {
	_ = reflect.TypeOf(x)
	_ = reflect.TypeOf(x).String()

	rv := reflect.ValueOf(x)
	_ = rv.Type()
	_ = rv.Type().String()

	_ = reflect.TypeOf(x).Size()
}
