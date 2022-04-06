package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQuickFix(t *testing.T) {
	dir := filepath.Join("testdata", "quickfix")

	goRun := func(t *testing.T, filename string) []byte {
		out, err := exec.Command("go", "run", filename).CombinedOutput()
		if err != nil {
			t.Fatalf("run %s: %v: %s", filename, err, out)
		}
		return out
	}

	for _, testDir := range readdir(t, dir) {
		key := filepath.Base(testDir)
		t.Run(key, func(t *testing.T) {
			f, err := os.Create(filepath.Join(testDir, "target.go"))
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(f.Name())

			data, err := os.ReadFile(filepath.Join(testDir, "before.go"))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := f.Write(data); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}

			args := []string{
				"--fix",
				"--no-color",
				"--quiet",
				fmt.Sprintf("./testdata/quickfix/%s/target.go", key),
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			if _, err := cmdLint(&stdout, &stderr, args); err != nil {
				t.Fatal(err)
			}
			if stderr.Len() != 0 {
				t.Fatalf("errors:\n%s", stderr.String())
			}

			want, err := os.ReadFile(filepath.Join(testDir, "out", "after.go"))
			if err != nil {
				t.Fatal(err)
			}
			have, err := os.ReadFile(filepath.Join(testDir, "target.go"))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(have, want); diff != "" {
				t.Fatalf("quick fixes result mismatch (+want -have):\n%s", diff)
			}
			oldOutput := goRun(t, filepath.Join(testDir, "target.go"))
			newOutput := goRun(t, filepath.Join(testDir, "out", "after.go"))
			if diff := cmp.Diff(newOutput, oldOutput); diff != "" {
				t.Fatalf("go run results mismatch (+old -new):\n%s", diff)
			}
		})
	}
}
