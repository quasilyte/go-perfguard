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

func Apply(src []byte, edits []TextEdit) []byte {
	if len(edits) == 0 {
		return src
	}

	sort.SliceStable(edits, func(i, j int) bool {
		if edits[i].StartOffset != edits[j].StartOffset {
			return edits[i].StartOffset < edits[j].StartOffset
		}
		return edits[i].EndOffset < edits[j].EndOffset
	})

	var buf bytes.Buffer
	buf.Grow(len(src))
	offset := 0
	for _, fix := range edits {
		// If we have a nested replacement, apply only outer replacement.
		if offset > fix.StartOffset {
			continue
		}

		buf.Write(src[offset:fix.StartOffset])
		buf.Write(fix.Replacement)

		offset = fix.EndOffset
	}
	buf.Write(src[offset:])

	return buf.Bytes()
}
