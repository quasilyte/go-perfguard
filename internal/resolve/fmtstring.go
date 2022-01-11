package resolve

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

type FmtInfo struct {
	Args []FmtArgInfo
}

type FmtArgInfo struct {
	Flag   byte
	Verb   byte
	ArgNum int
}

func (arg FmtArgInfo) String() string {
	buf := make([]byte, 0, 8)
	buf = append(buf, '%')
	if arg.Flag != 0 {
		buf = append(buf, arg.Flag)
	}
	buf = append(buf, arg.Verb, '<')
	buf = strconv.AppendInt(buf, int64(arg.ArgNum), 10)
	buf = append(buf, '>')
	return string(buf)
}

func FmtString(s string) (FmtInfo, bool) {
	var info FmtInfo

	if !strings.Contains(s, "%") {
		return info, true
	}
	if len(s) > 512 {
		return info, false
	}

	i := 0
	for i < len(s) {
		ch, chLen := utf8.DecodeRuneInString(s[i:])
		if ch == utf8.RuneError {
			return info, false
		}
		i += chLen
		if ch != '%' {
			continue
		}

		if i+1 > len(s) {
			return info, false
		}
		var arg FmtArgInfo
		verbOffset := 0
		switch s[i] {
		case '#', '+', '-', '0':
			arg.Flag = s[i]
			verbOffset++
			if i+2 > len(s) {
				return info, false
			}
		case '%':
			i++
			continue
		}

		switch verb := s[i+verbOffset]; verb {
		case 'd', 'f', 's', 'q', 'v':
			arg.Verb = verb
		default:
			return info, false
		}

		i += verbOffset
		arg.ArgNum = len(info.Args)
		info.Args = append(info.Args, arg)
	}

	return info, true
}
