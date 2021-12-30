package testfile

import (
	"fmt"
	"strings"
)

type Annotation struct {
	Filename string
	Line     int
	Text     string
}

func Parse(filename string, data []byte) ([]Annotation, error) {
	var result []Annotation

	lines := strings.Split(string(data), "\n")
	for i, l := range lines {
		if l == "" {
			continue
		}
		lineNum := i + 1
		start := strings.Index(l, "// want `")
		if start == -1 {
			continue
		}
		s := l[start+len("// want `"):]
		end := strings.IndexByte(s, '`')
		if start == -1 {
			return nil, fmt.Errorf("line %d: can't find closing `", lineNum)
		}
		s = s[:end]
		result = append(result, Annotation{
			Filename: filename,
			Line:     lineNum,
			Text:     s,
		})
	}

	return result, nil
}
