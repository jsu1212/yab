package main

import (
	"fmt"
	"path"
	"strings"
	"sort"
	"os"
)

func getAllRoots(f string) [][]string {
	var result [][]string
	var roots []string

	f = path.Clean(f)
	parts := strings.Split(f, "/")

	sort.Sort(sort.Reverse(sort.StringSlice(parts)))

	for _, part := range (parts) {
		roots = append(roots, part)
		result = append(result, roots)
	}

	return result
}

func main() {
	for _, root := range getAllRoots(os.Args[1]) {
		fmt.Println(root)
	}
}