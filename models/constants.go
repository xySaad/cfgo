package models

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var EnglishTitle = cases.Title(language.English)

const IMPORTS_TEMPLATE = `
import (
	{{range $k, $v := .}} "{{ $v }}"
	{{end}}
)`

const STRUCT_TEMPLATE string = `
type {{.Name}} struct {
{{range .Fields}}	{{.Key}} {{.Type}}
{{end}}}

var {{.NameLower}} = {{.Name}}{
{{range .Fields}}	{{.Key}}: {{.Value}},
{{end}}}

func Get{{.Name}}() {{.Name}} { return {{.NameLower}} }
`
