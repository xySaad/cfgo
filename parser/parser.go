package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/xySaad/cfgo/models"
)

type Parser struct {
	imports       map[string]string
	configImports map[string]string
}

func NewParser() Parser {
	return Parser{
		imports:       make(map[string]string),
		configImports: make(map[string]string),
	}
}
func (p *Parser) ConfigImports() map[string]string {
	return p.configImports
}
func (p *Parser) Imports() map[string]string {
	return p.imports
}

func (p *Parser) ParseTopLevelImports(json map[string]any) {
	for key, value := range json {
		if !strings.HasPrefix(key, "@") {
			continue
		}

		file := reflect.ValueOf(value).String()
		p.configImports[key[1:]] = file

		switch filepath.Ext(file) {
		case ".env":
			p.imports["godotenv"] = "github.com/joho/godotenv"
		case ".json":
			// unsupported
		}
		delete(json, key)
	}
}

func (p *Parser) ParseJSON(name string, structTempl *template.Template, json map[string]any, wr io.Writer, recursive bool) {
	params := models.Params{Name: models.EnglishTitle.String(name), NameLower: name, Fields: nil}
	for key, value := range json {
		titleKey := models.EnglishTitle.String(key)
		field := models.Field{
			Key:   titleKey,
			Type:  reflect.TypeOf(value).String(),
			Value: value,
		}

		switch v := value.(type) {
		case map[string]any:
			field.Type = field.Key
			field.Value = key
			if recursive {
				p.ParseJSON(key, structTempl, v, wr, true)
			}
		case string:
			field.Value = p.parseImports(v)
		}

		params.Fields = append(params.Fields, field)
	}

	err := structTempl.Execute(wr, params)
	if err != nil {
		panic(err)
	}
}

func (p Parser) parseImports(path string) string {
	if after, ok := strings.CutPrefix(path, "@"); ok {
		parts := strings.Split(after, ".")
		if len(parts) != 2 {
			panic("invalid format: " + path)
		}

		_, ok := p.configImports[parts[0]]
		if !ok {
			fmt.Fprintf(os.Stderr, "WARNING: %s contains import symbol @ but the target file wasn't imported\n", path)
			return fmt.Sprintf(`"%s"`, path)
		}
		//assume it's env since it's the only type supported

		p.imports["os"] = "os"
		return fmt.Sprintf(`os.Getenv("%s")`, parts[1])
	}
	return fmt.Sprintf(`"%s"`, path)
}
