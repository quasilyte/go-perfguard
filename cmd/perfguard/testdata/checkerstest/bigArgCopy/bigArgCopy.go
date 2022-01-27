package checkerstest

func f(x bigObject) {}

func ignored(_ bigObject) {}
func ignored2(args ...bigObject)
func ignored3(args ...interface{})

type example struct{}

func (example) f(y [1000]byte) {}

func Warn() {
	var o bigObject
	f(o)                              // want `expensive x arg copy`
	f(bigObject(anotherBigObject(o))) // want `expensive x arg copy`

	o.byValue() // want `expensive o receiver copy (10000 bytes)`
	o.byPointer()

	var e example
	e.f([1000]byte{}) // want `expensive y arg copy (1000 bytes), consider passing it by pointer`
}

type bigObject struct {
	data [10000]byte
}

func (o bigObject) byValue()    {}
func (o *bigObject) byPointer() {}

type anotherBigObject struct {
	data [10000]byte
}

func Ignore(o bigObject) {
	// Argument is not really evaluated.
	_ = new([90000]byte)

	// Not a function call: it's a type conversion.
	_ = anotherBigObject(o)
	_ = bigObject(o)

	ignored(o)
	ignored2(o)
	ignored3(o)
}
