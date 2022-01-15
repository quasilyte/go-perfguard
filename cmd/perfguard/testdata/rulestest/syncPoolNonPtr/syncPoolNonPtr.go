package syncPoolNonPtr

import (
	"sync"
)

type r string
type u struct {
	l []byte
}
type e struct {
	g string
}
type uu struct {
	a [102]byte
	r string
	d int
}

func foo() {
	var s = sync.Pool{}

	gu := ""
	s.Put(gu) // want `don't use sync.Pool on non pointer objects`

	bar := r("")

	s.Put(bar) // want `don't use sync.Pool on non pointer objects`
	s.Put(&bar)

	uv := u{}
	s.Put(uv) // want `don't use sync.Pool on non pointer objects`
	s.Put(&uv)
	s.Put(u{}) // want `don't use sync.Pool on non pointer objects`

	ee := e{}
	s.Put(ee)  // want `don't use sync.Pool on non pointer objects`
	s.Put(e{}) // want `don't use sync.Pool on non pointer objects`

	uuu := uu{}
	s.Put(uuu) // want `don't use sync.Pool on non pointer objects`
	s.Put(&uuu)
	s.Put(uu{}) // want `don't use sync.Pool on non pointer objects`
	s.Put(0)    // want `don't use sync.Pool on non pointer objects`
	s.Put("")   // want `don't use sync.Pool on non pointer objects`
}
