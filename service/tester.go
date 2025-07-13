package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var line = "----------------------------------------------------------------------"

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	bold   = "\033[1m"
)

type Test struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Problem struct {
	Name  string `json:"name"`
	Group string `json:"group"`
	Tests []Test `json:"tests"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: tester <program_binary> <test_prefix>")
		os.Exit(1)
	}

	binPath := "./" + os.Args[1]
	prefix := os.Args[2]

	// Ищем все файлы тестов, начинающиеся с prefix
	files, err := filepath.Glob("./tests/" + prefix + ".*.json")
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		fmt.Printf("No test files found with prefix '%s'\n", prefix)
		os.Exit(1)
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("Failed to read %s: %v\n", file, err)
			continue
		}

		var prob Problem
		if err := json.Unmarshal(data, &prob); err != nil {
			fmt.Printf("Failed to parse JSON in %s: %v\n", file, err)
			continue
		}

		for i, test := range prob.Tests {
			fmt.Printf("Test #%d: ", i+1)
			ok, got, want := runTest(binPath, test.Input, test.Output)
			if ok {
				fmt.Println(green + "✅" + reset)
			} else {
				fmt.Println("❌")
				fmt.Println(bold + "Input:" + reset)
				fmt.Println(line)
				fmt.Print(test.Input)
				fmt.Println(line)
				fmt.Println(bold + "Expected output:" + reset)
				fmt.Println(line)
				fmt.Println(want)
				fmt.Println(line)
				fmt.Println(bold + "Got output:" + reset)
				fmt.Println(line)
				fmt.Println(got)
				fmt.Println(line)
			}
		}
	}
}

func runTest(binPath, input, expectedOutput string) (bool, string, string) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return false, "", expectedOutput
	}

	got := normalize(stdout.String())
	want := normalize(expectedOutput)

	return got == want, got, want
}

var spaceCollapse = regexp.MustCompile(`\s+`)

// normalize - Заменяет все \n, \t, множественные пробелы и т.д. на один пробел
func normalize(s string) string {
	return strings.TrimSpace(spaceCollapse.ReplaceAllString(s, " "))
}
