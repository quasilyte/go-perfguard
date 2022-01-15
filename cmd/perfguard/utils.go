package main

import (
	"go/ast"
	"strings"
)

func isAutogenFile(f *ast.File) bool {
	for _, comment := range f.Comments {
		if isAutogenComment(comment) {
			return true
		}
	}
	return false
}

func isAutogenComment(comment *ast.CommentGroup) bool {
	generated := false
	doNotEdit := false
	for _, c := range comment.List {
		s := strings.ToLower(c.Text)
		if !generated {
			generated = strings.Contains(s, " code generated ") ||
				strings.Contains(s, " generated by ")
		}
		if !doNotEdit {
			doNotEdit = strings.Contains(s, "do not edit") ||
				strings.Contains(s, "don't edit")
		}
		if generated && doNotEdit {
			return true
		}
	}
	return false
}
