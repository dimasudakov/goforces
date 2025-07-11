package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func template(taskLetter string) string {
	return fmt.Sprintf(`package main

func main() {
	// t := 1
	t := scanT[int]()

	for ; t > 0; t-- {
		solve%s()
	}

	_ = out.Flush()
}

func solve%s() {

}
`, taskLetter, taskLetter)
}

type TestData struct {
	Name string `json:"name"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var testData TestData
	if err := json.Unmarshal(body, &testData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := saveTest(body, testData.Name); err != nil {
		http.Error(w, "Failed to save test", http.StatusInternalServerError)
		return
	}
	if err := createSolutionTemplate(testData.Name); err != nil {
		http.Error(w, "Failed to create solution template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)
	port := 10043
	log.Printf("Listening on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func saveTest(data []byte, name string) error {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", "    ") // 4 пробела для отступа
	if err != nil {
		return err
	}

	dir := "./tests"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.json", name)
	path := filepath.Join(dir, filename)
	return os.WriteFile(path, prettyJSON.Bytes(), 0644)
}

func createSolutionTemplate(taskName string) error {
	taskID := extractTaskID(taskName)
	dir := "./"
	fileName := fmt.Sprintf("task_%s.go", taskID)
	outFile := filepath.Join(dir, fileName)

	// Проверка: если файл уже существует — не перезаписываем
	if _, err := os.Stat(outFile); err == nil {
		fmt.Printf("File %s already exists, skipping.\n", outFile)
		return err
	}

	err := os.WriteFile(outFile, []byte(template(taskID)), 0644)
	if err != nil {
		fmt.Printf("Failed to write file %s: %s\n", outFile, err)
		return err
	}

	return nil
}

func extractTaskID(s string) string {
	parts := strings.SplitN(s, ".", 2)
	return strings.TrimSpace(parts[0])
}
