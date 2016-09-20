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
	var parts sort.StringSlice = strings.Split(f, "/")

	for i:=len(parts)-1; i>=0; i-- {
		part := parts[i]
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

func fmt_root(root []string) string {
	result := root[0]
	for _, part := range root[1:] {
		result += fmt.Sprintf("[%s]", part)
	}
	return result
}

func main() {
	files := find_thrift_files()

	for _, file := range files {
		fmt.Println(file)
		for _, root := range getAllRoots(file) {
			fmt.Println(fmt_root(root))
		}
		fmt.Println("----------")
	}
}