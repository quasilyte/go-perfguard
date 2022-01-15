package main

import (
	"path/filepath"
	"testing"
)

func TestCheckers(t *testing.T) {
	checkers := readdir(t, filepath.Join("testdata", "checkerstest"))
	for _, name := range checkers {
		key := filepath.Base(name)
		runLintTest(t, "checkerstest", key)
	}
}
