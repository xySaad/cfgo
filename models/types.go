package models

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
