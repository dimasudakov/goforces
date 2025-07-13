package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samber/lo"
)

type SrcFile struct {
	Name         string
	Dependencies []string
}

const (
	inout      = "./inout.go"
	set        = "./set.go"
	multiset   = "./multiset.go"
	orderedSet = "./set_ordered.go"
	math       = "./math.go"

	// only deps
	rbtree = "./rbtree.go"
)

var (
	sourceFiles = []SrcFile{
		{
			Name: inout,
		},
		{
			Name: set,
		},
		{
			Name: multiset,
		},
		{
			Name: orderedSet,
			Dependencies: []string{
				rbtree,
			},
		},
		{
			Name: math,
		},
	}
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run merger.go <output-file>")
		return
	}
	taskLetter := os.Args[1]
	outFile := "./solutions/task_" + taskLetter + ".go"

	files := getSrcFilesList(taskLetter)

	var builder strings.Builder
	builder.WriteString("package main\n\n")

	importsSet := make(map[string]struct{})
	var codeParts []string

	for i, file := range files {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(f)
		var codeLines []string
		inImports := false

		if i == len(files)-1 {
			// последний файл - это файл с решением
			solutionDividerLine := fmt.Sprintf(
				"// %s solution %s",
				strings.Repeat("=", 50),
				strings.Repeat("=", 50),
			)
			codeLines = append(codeLines, solutionDividerLine, "")
		}

		for scanner.Scan() {
			line := scanner.Text()

			// Skip package main
			if strings.HasPrefix(line, "package main") {
				continue
			}

			// Collect imports
			if strings.HasPrefix(line, "import") {
				inImports = true
				if strings.Contains(line, "(") {
					continue // skip multi-import line start 'import "math"'
				} else {
					// single import
					importLine := strings.TrimSpace(line[6:])
					importLine = strings.Trim(importLine, `" `)
					importsSet[importLine] = struct{}{}
					inImports = false
					continue
				}
			}

			// Handle multi-line imports
			if inImports {
				if strings.Contains(line, ")") {
					inImports = false
					continue
				}
				importLine := strings.TrimSpace(line)
				importLine = strings.Trim(importLine, `"`)
				if importLine != "" {
					importsSet[importLine] = struct{}{}
				}
				continue
			}

			codeLines = append(codeLines, line)
		}

		codeParts = append(codeParts, strings.Join(codeLines, "\n"))
		_ = f.Close()
	}

	// Write collected imports
	if len(importsSet) > 0 {
		importsArr := make([]string, 0, len(importsSet))
		for importLine := range importsSet {
			importsArr = append(importsArr, importLine)
		}
		sort.Strings(importsArr)

		builder.WriteString("import (\n")
		for _, imp := range importsArr {
			builder.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
		}
		builder.WriteString(")\n\n")
	}

	// Write code parts
	for _, part := range codeParts {
		builder.WriteString(strings.TrimSpace(part))
		builder.WriteString("\n\n")
	}

	// Получаем директорию из пути
	dir := filepath.Dir(outFile)

	// Создаём директорию и все недостающие родительские директории
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(outFile, []byte(builder.String()), 0644)
	if err != nil {
		panic(err)
	}

	const (
		reset = "\033[0m"
		bold  = "\033[1m"
	)

	fmt.Printf(bold+"task_%s\n"+reset, os.Args[1])
}

func getSrcFilesList(taskLetter string) []string {
	solutionFileName := fmt.Sprintf("./task_%s.go", taskLetter)
	used := collectUsedNames(solutionFileName)

	filesSet := make(map[string]struct{})

	for _, file := range sourceFiles {
		declared := collectDeclaredNames(file.Name)

		for name := range declared {
			if used[name] {
				fmt.Printf("File %s: symbol %s is used in task_A.go\n", file, name)
				filesSet[file.Name] = struct{}{}
				for _, dep := range file.Dependencies {
					filesSet[dep] = struct{}{}
				}
				break
			}
		}
	}
	result := lo.Keys(filesSet)
	result = append(result, solutionFileName)

	return result
}

// collectUsedNames parses task_X.go and returns all used identifiers
func collectUsedNames(path string) map[string]bool {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	used := make(map[string]bool)

	ast.Inspect(node, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			used[ident.Name] = true
		}
		return true
	})

	fmt.Printf("Used in %v: \t\t%+v\n", path, used)

	return used
}

// collectDeclaredNames returns top-level function and type names declared in the given file
func collectDeclaredNames(path string) map[string]bool {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	declared := make(map[string]bool)

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil { // только функции, не методы
				declared[d.Name.Name] = true
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					declared[ts.Name.Name] = true
				}
			}
		}
	}

	fmt.Printf("Declared in %v: \t\t%+v\n", path, declared)

	return declared
}
