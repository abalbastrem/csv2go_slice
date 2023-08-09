package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	// PARAMS
	csvFile := "players.csv"
	goFileName := "players"
	goFile := "players.go"
	packageName := "inmemory"
	sliceName := "players"

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var values []string

	// Read and process each row
	for {
		row, err := reader.Read()
		if err != nil {
			break // Stop reading when no more rows
		}
		if len(row) > 0 {
			values = append(values, row[0]) // Assuming single-column CSV
		}
	}

	// Create or overwrite the Go source file
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
	t, err := template.New(goFileName).Parse(tmpl)
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
