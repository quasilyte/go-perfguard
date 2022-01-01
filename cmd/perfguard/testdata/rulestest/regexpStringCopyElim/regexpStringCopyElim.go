package rulestest

import (
	"regexp"
)

func Warn() {
	var s string
	var b []byte
	var p []byte

	re := regexp.MustCompile(`\w+`)

	_ = re.Match([]byte(s))            // want `re.Match([]byte(s)) => re.MatchString(s)`
	_ = re.FindIndex([]byte(s))        // want `re.FindIndex([]byte(s)) => re.FindStringIndex(s)`
	_ = re.FindAllIndex([]byte(s), -1) // want `re.FindAllIndex([]byte(s), -1) => re.FindAllStringIndex(s, -1)`

	_ = string(re.ReplaceAll([]byte(s), []byte("foo"))) // want `string(re.ReplaceAll([]byte(s), []byte("foo"))) => re.ReplaceAllString(s, "foo")`

	_ = re.MatchString(string(b))            // want `re.MatchString(string(b)) => re.Match(b)`
	_ = re.FindStringIndex(string(b))        // want `re.FindStringIndex(string(b)) => re.FindIndex(b)`
	_ = re.FindAllStringIndex(string(b), -1) // want `re.FindAllStringIndex(string(b), -1) => re.FindAllIndex(b, -1)`

	_ = []byte(re.ReplaceAllString(string(b), string(p))) // want `[]byte(re.ReplaceAllString(string(b), string(p))) => re.ReplaceAll(b, p)`

	{
		var content []byte
		var targetPath string
		content = []byte(re.ReplaceAllString(string(content), targetPath)) // want `[]byte(re.ReplaceAllString(string(content), targetPath)) => re.ReplaceAll(content, []byte(targetPath))`
	}
	{
		var content string
		var targetPath []byte
		content = string(re.ReplaceAll([]byte(content), targetPath)) // want `string(re.ReplaceAll([]byte(content), targetPath)) => re.ReplaceAllString(content, string(targetPath))`
	}
}

func Ignore() {
	re := regexp.MustCompile(`\w+`)

	{
		var s string
		var b []byte
		_ = re.MatchString(s)
		_ = re.FindStringIndex(s)
		_ = re.FindAllStringIndex(s, -1)
		_ = re.Match(b)
		_ = re.FindIndex(b)
		_ = re.FindAllIndex(b, -1)
	}

	{
		var content []byte
		var targetPath string
		content = re.ReplaceAll(content, []byte(targetPath))
	}
	{
		var content string
		var targetPath []byte
		content = re.ReplaceAllString(content, string(targetPath))
	}
}
