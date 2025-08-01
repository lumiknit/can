package server

import (
	_ "embed"
	"text/template"
)

//go:embed index.tpl.html
var indexTemplateContent string

var IndexTemplate *template.Template

func init() {
	var err error
	IndexTemplate, err = template.New("index").Parse(indexTemplateContent)
	if err != nil {
		panic("Failed to parse index template: " + err.Error())
	}
}
