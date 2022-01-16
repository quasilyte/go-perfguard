package syncPoolNonPtr

import (
	"sync"
)

type (
	r string
	u struct {
		l []byte
	}
	e struct {
		g string
	}
	uu struct {
		a [102]byte
	}
	rr     []byte
	d      []string
	dd     *d
	aliasD d
)

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

	s.Put([]int{123, 213}) // want `don't use sync.Pool on non pointer objects`
	s.Put(make(chan string))
	s.Put(make(map[int]int, 10))

	s.Put(rr{}) // want `don't use sync.Pool on non pointer objects`
	s.Put(&rr{})

	var ddObj dd = &d{"123", "1333"}
	s.Put(ddObj)

	s.Put(aliasD{}) // want `don't use sync.Pool on non pointer objects`
}

func (rec *r) FooBar() {
	k := sync.Pool{}

	k.Put(rec)
}
