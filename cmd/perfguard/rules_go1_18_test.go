//go:build go1.18
// +build go1.18

package main

import (
	"path/filepath"
	"testing"
)

func TestRulesGo1_18(t *testing.T) {
	rules := readdir(t, filepath.Join("testdata", "rulestest"))
	for _, name := range rules {
		key := filepath.Base(name)
		if ver := testVersionConstraints[key]; ver != "1.18" {
			continue
		}
	}
}
