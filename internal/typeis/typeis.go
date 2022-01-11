package typeis

import (
	"go/types"
)

func String(typ types.Type) bool {
	if typ, ok := typ.(*types.Basic); ok {
		if typ.Info()&types.IsString != 0 {
			return true
		}
	}
	return false
}

func ByteSlice(typ types.Type) bool {
	if typ, ok := typ.(*types.Slice); ok {
		if typ, ok := typ.Elem().(*types.Basic); ok {
			return typ.Kind() == types.Uint8
		}
	}
	return false
}
