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

func Slice(typ types.Type) bool {
	_, ok := typ.(*types.Slice)
	return ok
}

func Map(typ types.Type) bool {
	_, ok := typ.(*types.Map)
	return ok
}

func ByteSlice(typ types.Type) bool {
	if typ, ok := typ.(*types.Slice); ok {
		if typ, ok := typ.Elem().(*types.Basic); ok {
			return typ.Kind() == types.Uint8
		}
	}
	return false
}

func Named(typ types.Type, pkgPath, typeName string) bool {
	if namedType, ok := typ.(*types.Named); ok {
		obj := namedType.Obj()
		return obj.Name() == typeName && obj.Pkg().Path() == pkgPath
	}
	return false
}
