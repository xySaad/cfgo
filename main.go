package main

import (
	"bytes"
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

	packageName := filepath.Base(filepath.Dir(outputPath))
	if packageName == "." || packageName == string(filepath.Separator) {
		packageName = "main"
	}

	buffer := bytes.NewBuffer(nil)
	envUsed := transformObject(fileName, structTempl, mappedJSON, buffer, true)
	topLevel := fmt.Sprintln("package", packageName)
	if envUsed {
		topLevel += `import "os"` + "\n"
	}

	finalBuf := bytes.NewBuffer([]byte(topLevel))
	buffer.WriteTo(finalBuf)
	err = os.WriteFile(outputPath, finalBuf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func transformObject(name string, structTempl *template.Template, json map[string]any, wr io.Writer, recursive bool) (envUsed bool) {
	params := Params{Name: englishTitle.String(name), NameLower: name, Fields: nil}
	for key, value := range json {
		if strings.HasPrefix(key, "@") {
			// str := value.(string)
			// parseImports(key, str, wr)
			continue
		}

		titleKey := englishTitle.String(key)
		field := Field{
			Key:   titleKey,
			Type:  reflect.TypeOf(value).String(),
			Value: value,
		}

		switch v := value.(type) {
		case map[string]any:
			field.Type = field.Key
			field.Value = key
			if recursive {
				if transformObject(key, structTempl, v, wr, true) {
					envUsed = true
				}
			}
		case string:
			fmt.Println(key, value)
			field.Value, envUsed = parseImports(v)
		}

		params.Fields = append(params.Fields, field)
	}

	err := structTempl.Execute(wr, params)
	if err != nil {
		panic(err)
	}
	return
}
