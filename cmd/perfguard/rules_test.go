package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/quasilyte/go-perfguard/internal/testfile"
)

func TestRules(t *testing.T) {
	readdir := func(t *testing.T, dir string) []string {
		var filenames []string
		files, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range files {
			filenames = append(filenames, filepath.Join(dir, f.Name()))
		}
		return filenames
	}

	runRulesTest := func(t *testing.T, name string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join("testdata", "rulestest", name)
			args := []string{
				"./testdata/rulestest/" + name + "/...",
			}

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			if err := cmdLint(&stdout, &stderr, args); err != nil {
				t.Fatal(err)
			}

			filenames := readdir(t, dir)

			var annotations []testfile.Annotation
			for _, filename := range filenames {
				data, err := os.ReadFile(filename)
				if err != nil {
					t.Fatal(err)
				}
				fileAnnotations, err := testfile.Parse(filename, data)
				if err != nil {
					t.Fatalf("parse test file annotations: %v", err)
				}
				annotations = append(annotations, fileAnnotations...)
			}

			compareTestResults(t, annotations, stdout.Bytes())
		})
	}

	rules := readdir(t, filepath.Join("testdata", "rulestest"))
	for _, name := range rules {
		runRulesTest(t, filepath.Base(name))
	}
}

var outputLineRegexp = regexp.MustCompile(`(.*?):(\d+): (\w+): (.*)`)

func compareTestResults(t *testing.T, annotations []testfile.Annotation, output []byte) {
	t.Helper()

	type location struct {
		filename string
		line     int
	}

	wantWarnings := make(map[location]string)
	for _, a := range annotations {
		wantWarnings[location{filename: a.Filename, line: a.Line}] = a.Text
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(wd, "/") {
		wd += "/"
	}

	atoi := func(s string) int {
		v, err := strconv.Atoi(s)
		if err != nil {
			return -1
		}
		return v
	}

	haveWarnings := make(map[location][]string)
	for _, l := range strings.Split(string(output), "\n") {
		if l == "" {
			continue
		}
		if !strings.HasPrefix(l, wd) {
			continue
		}
		s := strings.TrimPrefix(l, wd)
		parts := outputLineRegexp.FindStringSubmatch(s)
		if parts == nil {
			continue
		}
		filename := parts[1]
		tag := parts[2]
		lineNum := atoi(parts[2])
		messageText := parts[4]

		_ = tag // not needed right now

		k := location{filename: filename, line: lineNum}
		haveWarnings[k] = append(haveWarnings[k], messageText)
	}

	unexected := make(map[location][]string)
	unmatched := make(map[location]string)
	for loc, warnings := range haveWarnings {
		want, ok := wantWarnings[loc]
		if !ok {
			for _, w := range warnings {
				unexected[loc] = append(unexected[loc], w)
			}
			continue
		}
		matched := false
		for _, w := range warnings {
			if strings.Contains(w, want) {
				if matched {
					unexected[loc] = append(unexected[loc], w)
					continue
				}
				matched = true
				continue
			} else {
				unexected[loc] = append(unexected[loc], w)
			}
		}
		if !matched {
			unmatched[loc] = want
		}
	}
	for loc, w := range wantWarnings {
		if _, ok := haveWarnings[loc]; ok {
			continue
		}
		unmatched[loc] = w
	}
	for loc, warnings := range unexected {
		for _, w := range warnings {
			t.Errorf("%s:%d: unexpected warn: %s", loc.filename, loc.line, w)
		}
	}
	for loc, w := range unmatched {
		t.Errorf("%s:%d: unmatched warn: %s", loc.filename, loc.line, w)
	}
}
