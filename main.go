package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")        // Backslash
	s = strings.ReplaceAll(s, "\"", "\\\"")        // Double quote
	s = strings.ReplaceAll(s, "\n", "\\n")         // Newline
	s = strings.ReplaceAll(s, "\r", "\\r")         // Carriage return
	s = strings.ReplaceAll(s, "\t", "\\t")         // Tab
	s = strings.ReplaceAll(s, "\u0022", "\\u0022") // Unicode escape for double quote
	s = strings.ReplaceAll(s, "\u0000", "\\u0000") // Unicode escape for null character
	return s
}

func main() {
	// PARAMS
	subject := "competitions"

	csvFile := subject + ".csv"
	goFile := subject + ".go"
	packageName := "inmemory"
	sliceName := subject
	templateName := subject
	// END PARAMS

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip the header row
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	var values []string

	// Read and process each row
	for {
		row, err := reader.Read()
		if err != nil {
			break // Stop reading when no more rows
		}
		if len(row) > 0 {
			escapedValue := escapeString(row[0])  // Escape troublesome characters
			values = append(values, escapedValue) // Assuming single-column CSV
		}
	}

	// Create or overwrite the Go source file (players.go)
	goSourceFile, err := os.Create(goFile)
	if err != nil {
		log.Fatal(err)
	}
	defer goSourceFile.Close()

	// Define a template
	tmpl := `package {{ .PackageName }}

var {{ .SliceName }} = []string{
{{- range .Values }}
	"{{ . }}",
{{- end }}
}
`

	// Parse the template
	t, err := template.New(templateName).Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	// Define template data
	data := struct {
		PackageName string
		SliceName   string
		Values      []string
	}{
		PackageName: packageName,
		SliceName:   sliceName,
		Values:      values,
	}

	// Execute the template and write to the Go source file
	err = t.Execute(goSourceFile, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Go source file generated:", goFile)
}
