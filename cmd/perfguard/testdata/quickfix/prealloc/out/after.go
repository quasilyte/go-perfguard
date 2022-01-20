package main

import "sort"

func main() {
	xs := []int{1, 2, 3, 4, 5, 2, 1, 9, 1, 8, 1}
	for _, x := range uniq(xs) {
		println(x)
	}
	for _, x := range uniq([]int{}) {
		println(x)
	}
	println(uniq(nil) == nil)
	println(uniq([]int{}) == nil)

	ys := []int{1, 9, 10, 309}
	for _, x := range intersection(xs, ys) {
		println(x)
	}
	for _, x := range intersection(xs, []int{}) {
		println(x)
	}
	println(intersection(nil, nil) == nil)
	println(intersection([]int{}, []int{}) == nil)

	m := map[string]string{
		"a": "b",
	}
	copymap(m)
}

func intersection(xs, ys []int) []int {
	all := map[int]struct{}{}
	for _, x := range xs {
		all[x] = struct{}{}
	}
	for _, y := range ys {
		all[y] = struct{}{}
	}
	var result []int
	for x := range all {
		result = append(result, x)
	}
	sort.Ints(result)
	return result
}

func uniq(xs []int) []int {
	var set map[int]struct{} = make(map[int]struct{})
	for _, x := range xs {
		set[x] = struct{}{}
	}
	result := make([]int, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	sort.Ints(result)
	return result
}

func copymap(m map[string]string) map[string]string {
	copied := make(map[string]string, len(m))
	for k, v := range m {
		copied[k] = v
	}
	return copied
}
