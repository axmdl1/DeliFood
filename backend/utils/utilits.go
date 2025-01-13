package utils

import "html/template"

// Utility functions for templates
func add(x, y int) int { return x + y }
func sub(x, y int) int { return x - y }
func iter(n int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i + 1
	}
	return result
}

var TmplFuncs = template.FuncMap{
	"add":  add,
	"sub":  sub,
	"iter": iter,
}
