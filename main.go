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

var englishTitle = cases.Title(language.English)

const STRUCT_TEMPLATE string = `
type {{.Name}} struct {
{{range .Fields}}	{{.Key}} {{.Type}}
{{end}}}

var {{.NameLower}} = {{.Name}}{
{{range .Fields}}	{{.Key}}: {{.Value}},
{{end}}}

func Get{{.Name}}() {{.Name}} { return {{.NameLower}} }
`

type Field struct {
	Type  string
	Key   string
	Value any
}

type Params struct {
	NameLower string
	Name      string
	Fields    []Field
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

	packageName := filepath.Base(filepath.Dir(outputPath))
	if packageName == "." || packageName == string(filepath.Separator) {
		packageName = "main"
	}

	fmt.Fprintf(file, "package %s\n", packageName)
	transformObject(fileName, structTempl, mappedJSON, file, true)
}

func strValue(key string, anyValue any) any {
	switch anyValue.(type) {
	case map[string]any:
		return key
	case string:
		return fmt.Sprintf(`"%s"`, anyValue)
	default:
		return anyValue
	}
}

func transformObject(name string, structTempl *template.Template, json map[string]any, wr io.Writer, recursive bool) {
	params := Params{Name: englishTitle.String(name), NameLower: name, Fields: nil}
	for key, value := range json {
		titleKey := englishTitle.String(key)
		field := Field{
			Key:   titleKey,
			Type:  reflect.TypeOf(value).String(),
			Value: strValue(key, value),
		}

		if v, isObject := value.(map[string]any); isObject {
			field.Type = field.Key
			if recursive {
				transformObject(key, structTempl, v, wr, true)
			}
		}

		params.Fields = append(params.Fields, field)
	}

	err := structTempl.Execute(wr, params)
	if err != nil {
		panic(err)
	}
}
