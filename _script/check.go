package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
)

func main() {
	var c checker

	steps := []struct {
		name string
		fn   func() error
	}{
		{"check autogen", c.checkAutogen},
	}
	for _, step := range steps {
		if err := step.fn(); err != nil {
			log.Fatalf("%s: error: %v", step.name, err)
		}
	}
}

type checker struct{}

func (c *checker) checkAutogen() error {
	filenames := []string{
		"opt_rules.go",
		"universal_rules.go",
	}
	readContents := func() ([][]byte, error) {
		var result [][]byte
		for _, filename := range filenames {
			data, err := os.ReadFile(filepath.Join("perfguard", "rulesdata", filename))
			if err != nil {
				return nil, err
			}
			result = append(result, data)
		}
		return result, nil
	}

	contents, err := readContents()
	if err != nil {
		return err
	}

	out, err := exec.Command("go", "generate", "./perfguard").CombinedOutput()
	if err != nil {
		return fmt.Errorf("run go generate: %v: %s", err, out)
	}

	newContents, err := readContents()
	if err != nil {
		return err
	}

	for i, filename := range filenames {
		oldData := contents[i]
		newData := newContents[i]
		if !bytes.Equal(oldData, newData) {
			diff := cmp.Diff(
				strings.Split(string(oldData), "\n"),
				strings.Split(string(newData), "\n"),
			)
			log.Printf("diff (-old +new):\n%s", diff)
			return fmt.Errorf("rulesdata for %s is outdated, run 'go generate ./perfguard'", filename)
		}
	}

	return nil
}
