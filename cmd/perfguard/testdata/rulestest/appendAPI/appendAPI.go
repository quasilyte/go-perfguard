package rulestest

import (
	"math/big"
	"strconv"
	"time"
)

func Warn(b []byte, i int, i64 int64, u64 uint64, t *time.Time) {
	b = append(b, strconv.Itoa(i)...) // want `b = append(b, strconv.Itoa(i)...) => b = strconv.AppendInt(b, int64(i), 10)`

	b = append(b, strconv.FormatInt(i64, 10)...) // want `b = append(b, strconv.FormatInt(i64, 10)...) => b = strconv.AppendInt(b, i64, 10)`
	b = append(b, strconv.FormatInt(i64, 16)...) // want `b = append(b, strconv.FormatInt(i64, 16)...) => b = strconv.AppendInt(b, i64, 16)`

	b = append(b, strconv.FormatUint(u64, 10)...) // want `b = append(b, strconv.FormatUint(u64, 10)...) => b = strconv.AppendUint(b, u64, 10)`
	b = append(b, strconv.FormatUint(u64, 8)...)  // want `b = append(b, strconv.FormatUint(u64, 8)...) => b = strconv.AppendUint(b, u64, 8)`

	b = append(b, t.Format(time.UnixDate)...) // want `b = append(b, t.Format(time.UnixDate)...) => b = t.AppendFormat(b, time.UnixDate)`

	{
		var bf big.Float
		b = append(b, bf.String()...)      // want `b = append(b, bf.String()...) => b = bf.Append(b, 'g', 10)`
		b = append(b, bf.Text('g', 15)...) // want `b = append(b, bf.Text('g', 15)...) => b = bf.Append(b, 'g', 15)`
	}
	{
		var bi big.Int
		b = append(b, bi.String()...) // want `b = append(b, bi.String()...) => b = bi.Append(b, 10)`
		b = append(b, bi.Text(16)...) // want `b = append(b, bi.Text(16)...) => b = bi.Append(b, 16)`
	}
}

func Ignore(b []byte, i int, i64 int64, u64 uint64, t *time.Time) {
	b = strconv.AppendInt(b, i64, 10)

	b = strconv.AppendInt(b, i64, 10)
	b = strconv.AppendInt(b, i64, 16)

	b = strconv.AppendUint(b, u64, 10)
	b = strconv.AppendUint(b, u64, 8)

	b = t.AppendFormat(b, time.UnixDate)

	{
		var bf big.Float
		b = bf.Append(b, 'g', 10)
		b = bf.Append(b, 'g', 15)
	}
	{
		var bi big.Int
		b = bi.Append(b, 10)
		b = bi.Append(b, 16)
	}
}
