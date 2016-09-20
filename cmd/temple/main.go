package main

import (
	"fmt"
	"github.com/thriftrw/thriftrw-go/ast"
	"github.com/thriftrw/thriftrw-go/idl"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"reflect"
)

// TOOD: kill camel case
// TODO: try to reduce copy/paste in the different stages

const IDL_PATH string = "./uber-idl"

func getAllRoots(f string) []string {
	var result []string
	var roots []string

	f = path.Clean(f)
	var parts sort.StringSlice = strings.Split(f, "/")
	//fmt.Println(parts)

	for i := len(parts) - 1; i >= 0; i-- {
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
		if err != nil {
			return err
		}

		if path.Ext(p) == ".thrift" {
			result = append(result, p)
		}

		return nil
	})

	if err != nil {
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

func keys(m map[string]string) []string {
	var result []string
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

func complete(options []string, prefix string) []string {
	var result []string

	for _, opt := range options {
		if strings.HasPrefix(opt, prefix) {
			result = append(result, opt)
		}
	}

	return result
}

func shellAutocomplete() {
	line := os.Args[1]
	args := strings.Fields(line)
	args = args[1:] // remove actual name of command

	file_map := make(map[string]string)
	populate_thrift_file_map(file_map)

	if strings.HasSuffix(line, " ") {
		args = append(args, "") // we're working on the next new argument
	}

	var opts []string

	if len(args) == 1 {
		// looking for idl file
		idl := args[0]

		idls := keys(file_map)

		opts = complete(idls, idl)

		if len(opts) == 1 {
			opts = []string{file_map[opts[0]]}
		}

	} else if len(args) == 2 {
		// looking for thrift service name
		fileName := args[0]
		serviceName := args[1]

		bytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}

		thrift, err := idl.Parse(bytes)
		if err != nil {
			panic(fileName)
		}

		for _, def := range thrift.Definitions {
			switch t := def.(type) {
			case *ast.Service:
				opts = append(opts, t.Name)
			}
		}

		opts = complete(opts, serviceName)

	} else if len(args) == 3 {
		// looking for function name
		fileName := args[0]
		serviceName := args[1]
		functionName := args[2]

		bytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}

		thrift, err := idl.Parse(bytes)
		if err != nil {
			panic(fileName)
		}

		var service *ast.Service
		for _, def := range thrift.Definitions {
			switch t := def.(type) {
			case *ast.Service:
				if t.Name == serviceName {
					service = t
					break
				}
			}
		}

		if service == nil {
			log.Fatalf("No service named %s", serviceName)
		}

		for _, proc := range service.Functions {
			opts = append(opts, proc.Name)
		}

		opts = complete(opts, functionName)
	}

	fmt.Println(strings.Join(opts, " "))
}

func defaultValue(astType ast.Type, structs map[string]*ast.Struct) string {
	switch t := astType.(type) {
	case ast.BaseType:
		switch t.ID {
		case ast.BoolTypeID:
			return `"false"`
		case ast.I8TypeID, ast.I16TypeID, ast.I32TypeID, ast.I64TypeID:
			return "0"
		case ast.DoubleTypeID:
			return "0.0"
		case ast.StringTypeID:
			return `''`
		case ast.BinaryTypeID:
			return `''`	//TODO: needs somt thought
		}
	case ast.TypeReference:
		s := structs[t.Name]
		return fmt.Sprintf(`"struct %s (%d fields)"`, s.Name, len(s.Fields))

	default:
		log.Fatalf("Unknown type: %s", reflect.TypeOf(astType).String())
	}
	panic("Should not get here")
}


func generateTemplate(function *ast.Function, structs map[string]*ast.Struct) {
	// TODO: structs feels like it should be a member variable of something

	for _, param := range function.Parameters {
		fmt.Printf("%s: %s\n", param.Name, defaultValue(param.Type, structs))
	}
}

func main() {
	if len(os.Getenv("SHELL_AUTOCOMPLETE")) > 0 {
		shellAutocomplete()
	} else {
		fileName := os.Args[1]
		serviceName := os.Args[2]
		functionName := os.Args[3]

		bytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}

		thrift, err := idl.Parse(bytes)
		if err != nil {
			panic(fileName)
		}

		structs := make(map[string]*ast.Struct)

		for _, def := range thrift.Definitions {
			switch t := def.(type) {
			case *ast.Struct:
				structs[t.Name] = t
			}
		}

		var service *ast.Service
		for _, def := range thrift.Definitions {
			switch t := def.(type) {
			case *ast.Service:
				if t.Name == serviceName {
					service = t
					break
				}
			}
		}

		if service == nil {
			log.Fatalf("No service named %s", serviceName)
		}

		var function *ast.Function




		for _, proc := range service.Functions {
			if proc.Name == functionName {
				function = proc
				break
			}
		}

		if function == nil {
			log.Fatalf("No function named %s", serviceName)
		}

		generateTemplate(function, structs)
	}
}
