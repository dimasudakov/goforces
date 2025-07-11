package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	// исходники (не файлы с решениями задач)
	sourceFiles = []string{
		"./inout.go",
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

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(f)
		var codeLines []string
		inImports := false

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
	result := sourceFiles
	result = append(result, fmt.Sprintf("./task_%s.go", taskLetter))
	return result
}
