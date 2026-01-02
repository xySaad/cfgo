package models

import (
	"regexp"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var PlaceHolderPattern = regexp.MustCompile(`\$\{([\w.]+)\}`)
var EnglishTitle = cases.Title(language.English, cases.NoLower)

const IMPORTS_TEMPLATE = `
import (
	{{range $k, $v := .}} "{{ $v }}"
	{{end}}
)`

const ENV_LOADER_TEMPLATE = `
func mustLoadEnv(path string) map[string]string {
	env, err := godotenv.Read(path)
	if err != nil {
		panic(err)
	}
	return env
}
`

const STRUCT_TEMPLATE string = `
type {{.Name}} struct {
{{range .Fields}}	{{.Key}} {{.Type}}
{{end}}}

var {{.NameLower}} = {{.Name}}{
{{range .Fields}}	{{.Key}}: {{.Value}},
{{end}}}

func Get{{.Name}}() {{.Name}} { return {{.NameLower}} }
`
