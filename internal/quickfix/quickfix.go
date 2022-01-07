package quickfix

import (
	"bytes"
	"sort"
)

// TextEdit is a suggested issue fix.
//
// More or less, it represents our version of the https://godoc.org/golang.org/x/tools/go/analysis#TextEdit
// which is a part of https://godoc.org/golang.org/x/tools/go/analysis#SuggestedFix.
type TextEdit struct {
	StartOffset int
	EndOffset   int

	// Replacement is a text to be inserted as a Replacement.
	Replacement []byte
}

func Sort(slice interface{}, get func(i int) TextEdit) {
	sort.SliceStable(slice, func(i, j int) bool {
		x := get(i)
		y := get(j)
		if x.StartOffset != y.StartOffset {
			return x.StartOffset < y.StartOffset
		}
		return x.EndOffset < y.EndOffset
	})
}

// Apply returns updated src with edits applied to it.
//
// Second return value contains indexes from edits slice
// that were not applied due to overlapping.
func Apply(src []byte, edits []TextEdit) (out []byte, indexes []int) {
	if len(edits) == 0 {
		return src, nil
	}

	Sort(edits, func(i int) TextEdit {
		return edits[i]
	})

	var overlapping []int
	var buf bytes.Buffer
	buf.Grow(len(src))
	offset := 0
	for i, fix := range edits {
		// If we have a nested replacement, apply only outer replacement.
		if offset > fix.StartOffset {
			overlapping = append(overlapping, i)
			continue
		}

		buf.Write(src[offset:fix.StartOffset])
		buf.Write(fix.Replacement)

		offset = fix.EndOffset
	}
	buf.Write(src[offset:])

	return buf.Bytes(), overlapping
}
