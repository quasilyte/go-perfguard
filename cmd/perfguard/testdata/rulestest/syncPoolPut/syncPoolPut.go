package syncPoolNonPtr

import (
	"io"
	"sync"
	"unsafe"
)

var pool sync.Pool

func Warn() {
	{
		var i int
		pool.Put(i)  // want `non-pointer values in sync.Pool involve extra allocation`
		pool.Put(10) // want `non-pointer values in sync.Pool involve extra allocation`
	}

	{
		var u uintptr
		pool.Put(u) // want `non-pointer values in sync.Pool involve extra allocation`
	}

	{
		type point struct {
			x, y int
		}
		var pt point
		pool.Put(pt) // want `non-pointer values in sync.Pool involve extra allocation`
	}

	{
		var a [4]int
		pool.Put(a) // want `non-pointer values in sync.Pool involve extra allocation`
	}

	{
		var b []byte
		pool.Put(b) // want `non-pointer values in sync.Pool involve extra allocation`
	}
}

func Ignore() {
	{
		var i int
		pool.Put(&i)
	}

	{
		type myChan chan int
		var ch chan int
		var myCh myChan
		pool.Put(ch)
		pool.Put(&ch)
		pool.Put(myCh)
		pool.Put(&myCh)
	}

	{
		type point struct {
			x, y int
		}
		var pt point
		pool.Put(&pt)
	}

	{
		type myMap map[int]string
		var m map[int]string
		var myM myMap
		pool.Put(m)
		pool.Put(&m)
		pool.Put(myM)
		pool.Put(&myM)
	}

	{
		type myEface interface{}
		var eface interface{}
		var e2 myEface
		var r io.Reader
		pool.Put(eface)
		pool.Put(&eface)
		pool.Put(e2)
		pool.Put(&e2)
		pool.Put(r)
		pool.Put(&r)
	}

	{
		type myFunc func(int) bool
		var f1 func()
		var f2 myFunc
		pool.Put(f1)
		pool.Put(&f1)
		pool.Put(f2)
		pool.Put(&f2)
		pool.Put(Ignore)
	}

	{
		var ptr unsafe.Pointer
		pool.Put(ptr)
		pool.Put(&ptr)
	}
}
