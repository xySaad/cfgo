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
	template      *template.Template
	root          map[string]any
}

func NewParser(root map[string]any, template *template.Template) Parser {
	return Parser{
		root:          root,
		imports:       make(map[string]string),
		configImports: make(map[string]string),
		template:      template,
	}
}

func (p *Parser) Parse(fileName string, wr io.Writer) {
	p.ParseJSON(fileName, p.root, wr)
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

func (p *Parser) resolveValue(val string) string {
	result := models.PlaceHolderPattern.ReplaceAllStringFunc(val, func(match string) string {
		key := match[2 : len(match)-1]

		parent := p.root
		parts := strings.Split(key, ".")
		for i := 0; i < len(parts)-1; i++ {
			part := parts[i]
			object, ok := parent[part].(map[string]any)
			if !ok {
				panic("key: " + key + " not found")
			}
			parent = object
		}
		last := parts[len(parts)-1]
		return fmt.Sprintf("%v", parent[last])
	})

	return fmt.Sprintf(`%s`, result)
}

func (p *Parser) ParseJSON(name string, json map[string]any, wr io.Writer) {
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
			p.ParseJSON(key, v, wr)
		case string:
			field.Value = p.parseImports(v)
		}

		params.Fields = append(params.Fields, field)
	}

	for i, field := range params.Fields {
		if str, ok := field.Value.(string); ok {
			field.Value = p.resolveValue(str)
			params.Fields[i] = field
		}
	}

	err := p.template.Execute(wr, params)
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

		return fmt.Sprintf(`env_%s["%s"]`, parts[0], parts[1])
	}
	return fmt.Sprintf(`"%s"`, path)
}
