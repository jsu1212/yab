package main

import (
	"fmt"
	"path"
	"strings"
	"sort"
	"os"
	"path/filepath"
)

const IDL_PATH string = "./uber-idl"

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

func find_thrift_files() []string {
	var result []string

	err := filepath.Walk(IDL_PATH, func(p string, info os.FileInfo, err error) error {
		if (err != nil) {
			return err
		}

		if path.Ext(p) == ".thrift" {
			result = append(result, p)
		}

		return nil
	})

	if (err != nil) {
		panic(err)
	}

	return result
}

func main() {
	files := find_thrift_files()

	for _, file := range files {
		fmt.Println(file)
	}
}