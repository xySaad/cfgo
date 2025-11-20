package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const STRUCT_TEMPLATE string = `type {{.Name}} struct {
{{range .Fields}}	{{.Key}} {{.Type}}
{{end}}}

`

type Field struct {
	Type string
	Key  string
}

type Params struct {
	Name   string
	Fields []Field
}

func main() {
	input := os.Args[1]
	output := os.Args[2]

	bfile, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}

	filePathNoExt := strings.TrimSuffix(input, filepath.Ext(input))
	mappedJSON := map[string]any{}
	err = json.Unmarshal(bfile, &mappedJSON)
	if err != nil {
		panic(err)
	}

	structTempl, err := template.New("struct").Parse(STRUCT_TEMPLATE)
	if err != nil {
		panic(err)
	}

	fileName := filepath.Base(filePathNoExt)
	outputPath := path.Join(output)
	file, err := os.OpenFile(outputPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "package %s\n\n", fileName)
	structName := cases.Title(language.English).String(fileName)
	transformObject(structName, structTempl, mappedJSON, file, true)
}

func transformObject(name string, structTempl *template.Template, json map[string]any, wr io.Writer, recursive bool) {
	params := Params{Name: name, Fields: nil}
	for key, _type := range json {
		titleKey := cases.Title(language.English).String(key)
		field := Field{Key: titleKey, Type: reflect.TypeOf(_type).String()}

		if v, isObject := _type.(map[string]any); isObject {
			field.Type = titleKey
			if recursive {
				transformObject(field.Type, structTempl, v, wr, true)
			}
		}

		params.Fields = append(params.Fields, field)
	}

	err := structTempl.Execute(wr, params)
	if err != nil {
		panic(err)
	}
}
