package rulestest

import (
	"fmt"
	"strconv"
)

func Warn() {
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var u uint
	var u8 uint8
	var u16 uint16
	var u32 uint32
	var u64 uint64

	_ = fmt.Sprintf("%d", i)   // want `fmt.Sprintf("%d", i) => strconv.Atoi(i)`
	_ = fmt.Sprintf("%d", i8)  // want `fmt.Sprintf("%d", i8) => strconv.FormatInt(int64(i8), 10)`
	_ = fmt.Sprintf("%d", i16) // want `fmt.Sprintf("%d", i16) => strconv.FormatInt(int64(i16), 10)`
	_ = fmt.Sprintf("%d", i32) // want `fmt.Sprintf("%d", i32) => strconv.FormatInt(int64(i32), 10)`
	_ = fmt.Sprintf("%d", i64) // want `fmt.Sprintf("%d", i64) => strconv.FormatInt(i64, 10)`
	_ = fmt.Sprintf("%d", u)   // want `fmt.Sprintf("%d", u) => strconv.FormatUint(uint64(u), 10)`
	_ = fmt.Sprintf("%d", u8)  // want `fmt.Sprintf("%d", u8) => strconv.FormatUint(uint64(u8), 10)`
	_ = fmt.Sprintf("%d", u16) // want `fmt.Sprintf("%d", u16) => strconv.FormatUint(uint64(u16), 10)`
	_ = fmt.Sprintf("%d", u32) // want `fmt.Sprintf("%d", u32) => strconv.FormatUint(uint64(u32), 10)`
	_ = fmt.Sprintf("%d", u64) // want `fmt.Sprintf("%d", u64) => strconv.FormatUint(u64, 10)`

	_ = fmt.Sprintf("%v", i)   // want `fmt.Sprintf("%v", i) => strconv.Atoi(i)`
	_ = fmt.Sprintf("%v", i8)  // want `fmt.Sprintf("%v", i8) => strconv.FormatInt(int64(i8), 10)`
	_ = fmt.Sprintf("%v", i16) // want `fmt.Sprintf("%v", i16) => strconv.FormatInt(int64(i16), 10)`
	_ = fmt.Sprintf("%v", i32) // want `fmt.Sprintf("%v", i32) => strconv.FormatInt(int64(i32), 10)`
	_ = fmt.Sprintf("%v", i64) // want `fmt.Sprintf("%v", i64) => strconv.FormatInt(i64, 10)`
	_ = fmt.Sprintf("%v", u)   // want `fmt.Sprintf("%v", u) => strconv.FormatUint(uint64(u), 10)`
	_ = fmt.Sprintf("%v", u8)  // want `fmt.Sprintf("%v", u8) => strconv.FormatUint(uint64(u8), 10)`
	_ = fmt.Sprintf("%v", u16) // want `fmt.Sprintf("%v", u16) => strconv.FormatUint(uint64(u16), 10)`
	_ = fmt.Sprintf("%v", u32) // want `fmt.Sprintf("%v", u32) => strconv.FormatUint(uint64(u32), 10)`
	_ = fmt.Sprintf("%v", u64) // want `fmt.Sprintf("%v", u64) => strconv.FormatUint(u64, 10)`

	_ = fmt.Sprint(i)   // want `fmt.Sprint(i) => strconv.Atoi(i)`
	_ = fmt.Sprint(i8)  // want `fmt.Sprint(i8) => strconv.FormatInt(int64(i8), 10)`
	_ = fmt.Sprint(i16) // want `fmt.Sprint(i16) => strconv.FormatInt(int64(i16), 10)`
	_ = fmt.Sprint(i32) // want `fmt.Sprint(i32) => strconv.FormatInt(int64(i32), 10)`
	_ = fmt.Sprint(i64) // want `fmt.Sprint(i64) => strconv.FormatInt(i64, 10)`
	_ = fmt.Sprint(u)   // want `fmt.Sprint(u) => strconv.FormatUint(uint64(u), 10)`
	_ = fmt.Sprint(u8)  // want `fmt.Sprint(u8) => strconv.FormatUint(uint64(u8), 10)`
	_ = fmt.Sprint(u16) // want `fmt.Sprint(u16) => strconv.FormatUint(uint64(u16), 10)`
	_ = fmt.Sprint(u32) // want `fmt.Sprint(u32) => strconv.FormatUint(uint64(u32), 10)`
	_ = fmt.Sprint(u64) // want `fmt.Sprint(u64) => strconv.FormatUint(u64, 10)`

	_ = fmt.Sprintf("%x", i)   // want `fmt.Sprintf("%x", i) => strconv.FormatInt(int64(i), 16)`
	_ = fmt.Sprintf("%x", i8)  // want `fmt.Sprintf("%x", i8) => strconv.FormatInt(int64(i8), 16)`
	_ = fmt.Sprintf("%x", i16) // want `fmt.Sprintf("%x", i16) => strconv.FormatInt(int64(i16), 16)`
	_ = fmt.Sprintf("%x", i32) // want `fmt.Sprintf("%x", i32) => strconv.FormatInt(int64(i32), 16)`
	_ = fmt.Sprintf("%x", i64) // want `fmt.Sprintf("%x", i64) => strconv.FormatInt(i64, 16)`
	_ = fmt.Sprintf("%x", u)   // want `fmt.Sprintf("%x", u) => strconv.FormatUint(uint64(u), 16)`
	_ = fmt.Sprintf("%x", u8)  // want `fmt.Sprintf("%x", u8) => strconv.FormatUint(uint64(u8), 16)`
	_ = fmt.Sprintf("%x", u16) // want `fmt.Sprintf("%x", u16) => strconv.FormatUint(uint64(u16), 16)`
	_ = fmt.Sprintf("%x", u32) // want `fmt.Sprintf("%x", u32) => strconv.FormatUint(uint64(u32), 16)`
	_ = fmt.Sprintf("%x", u64) // want `fmt.Sprintf("%x", u64) => strconv.FormatUint(u64, 16)`
}

func Ignore() {
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var u uint
	var u8 uint8
	var u16 uint16
	var u32 uint32
	var u64 uint64

	_ = strconv.Atoi(i)
	_ = strconv.FormatInt(int64(i8), 10)
	_ = strconv.FormatInt(int64(i16), 10)
	_ = strconv.FormatInt(int64(i32), 10)
	_ = strconv.FormatInt(i64, 10)
	_ = strconv.FormatUint(uint64(u), 10)
	_ = strconv.FormatUint(uint64(u8), 10)
	_ = strconv.FormatUint(uint64(u16), 10)
	_ = strconv.FormatUint(uint64(u32), 10)
	_ = strconv.FormatUint(u64, 10)

	_ = strconv.Atoi(i)
	_ = strconv.FormatInt(int64(i8), 10)
	_ = strconv.FormatInt(int64(i16), 10)
	_ = strconv.FormatInt(int64(i32), 10)
	_ = strconv.FormatInt(i64, 10)
	_ = strconv.FormatUint(uint64(u), 10)
	_ = strconv.FormatUint(uint64(u8), 10)
	_ = strconv.FormatUint(uint64(u16), 10)
	_ = strconv.FormatUint(uint64(u32), 10)
	_ = strconv.FormatUint(u64, 10)

	_ = strconv.Atoi(i)
	_ = strconv.FormatInt(int64(i8), 10)
	_ = strconv.FormatInt(int64(i16), 10)
	_ = strconv.FormatInt(int64(i32), 10)
	_ = strconv.FormatInt(i64, 10)
	_ = strconv.FormatUint(uint64(u), 10)
	_ = strconv.FormatUint(uint64(u8), 10)
	_ = strconv.FormatUint(uint64(u16), 10)
	_ = strconv.FormatUint(uint64(u32), 10)
	_ = strconv.FormatUint(u64, 10)

	_ = strconv.FormatInt(int64(i), 16)
	_ = strconv.FormatInt(int64(i8), 16)
	_ = strconv.FormatInt(int64(i16), 16)
	_ = strconv.FormatInt(int64(i32), 16)
	_ = strconv.FormatInt(i64, 16)
	_ = strconv.FormatUint(uint64(u), 16)
	_ = strconv.FormatUint(uint64(u8), 16)
	_ = strconv.FormatUint(uint64(u16), 16)
	_ = strconv.FormatUint(uint64(u32), 16)
	_ = strconv.FormatUint(u64, 16)
}
