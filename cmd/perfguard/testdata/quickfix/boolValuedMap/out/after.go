package main

import "sort"

func main() {
	stringsX := []string{
		"",
		"a",
		"str",
		"hello",
		"world",
		"str",
		"a",
	}
	stringsY := []string{
		"hello",
		"world",
	}

	println(intersectionLen(stringsX, stringsY))
	println(intersectionLen(stringsX, stringsX))
	for _, s := range uniq(stringsX) {
		println(s)
	}
	for _, s := range uniq(stringsY) {
		println(s)
	}
	for _, s := range keep(stringsX, stringsY) {
		println(s)
	}
	for _, s := range keep(stringsY, stringsX) {
		println(s)
	}
	for _, s := range keep(stringsX, stringsX) {
		println(s)
	}
}

func uniq(xs []string) []string {
	set := make(map[string]struct{}, len(xs)) // want `change map[T]bool to map[T]struct{}`
	for _, x := range xs {
		set[x] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func intersectionLen(xs, ys []string) int {
	intersection := make(map[string]struct{}, len(xs)) // want `change map[T]bool to map[T]struct{}`
	for _, x := range xs {
		intersection[x] = struct{}{}
	}
	for _, y := range ys {
		intersection[y] = struct{}{}
	}
	return len(intersection)
}

func keep(orig, allowlist []string) []string {
	set := make(map[string]struct{}, len(allowlist))
	for _, s := range allowlist {
		set[s] = struct{}{}
	}
	var result []string
	for _, s := range orig {
		if _, ok := set[s]; !ok {
			continue
		}
		result = append(result, s)
	}
	return result
}
