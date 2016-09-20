package main

import (
	"fmt"
	"path"
	"strings"
	"sort"
	"os"
	"path/filepath"
)

// TOOD: kill camel case

const IDL_PATH string = "./uber-idl"

func getAllRoots(f string) []string {
	var result []string
	var roots []string

	f = path.Clean(f)
	var parts sort.StringSlice = strings.Split(f, "/")
	//fmt.Println(parts)

	for i:=len(parts)-1; i>=0; i-- {
		part := parts[i]
		roots = append(roots, part)
		full := strings.Join(roots, "/")
		result = append(result, full)
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

func fmt_root(root string) string {
	parts := strings.Split(root, "/")
	result := parts[0]
	for _, part := range parts[1:] {
		result += fmt.Sprintf("[%s]", part)
	}
	return result
}

func populate_thrift_file_map(file_map map[string]string) {
	files := find_thrift_files()
	rootCounts := make(map[string]int)

	// find duplicate thrift files, working backwards along the path
	for _, full := range files {
		for _, root := range getAllRoots(full) {
			key := fmt_root(root)
			rootCounts[key]++
		}
	}

	// for each file, insert the shortest non-duplicate root along the backwards path
	for _, full := range files {
		for _, root := range getAllRoots(full) {
			key := fmt_root(root)
			count := rootCounts[key]
			if count == 1 {
				file_map[key] = full
				break
			}
		}
	}
}

func main() {
	//files := find_thrift_files()

	//for _, file := range files {
	//	fmt.Println(file)
	//	for _, root := range getAllRoots(file) {
	//		fmt.Println(fmt_root(root))
	//	}
	//	fmt.Println("----------")
	//}

	file_map := make(map[string]string)
	populate_thrift_file_map(file_map)

	for k, v := range file_map {
		fmt.Println(k, v)
	}
}