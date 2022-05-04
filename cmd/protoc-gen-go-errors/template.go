package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}
var {{.LowerCamelValue}} *errors.RespError
{{- end }}

func init() {
{{- range .Errors }}
{{.LowerCamelValue}} = errors.New({{.HTTPCode}}, "{{.Key}}", "{{.Pretty}}", {{.Name}}_{{.Value}}.String())
errors.Register({{.LowerCamelValue}})
{{- end }}
}

{{ range .Errors }}
{{if .HasComment}}{{.Comment}}{{end}}func {{.UpperCamelValue}}() errors.Error {
	 return {{.LowerCamelValue}}
}
{{ end }}
`

type errorInfo struct {
	Name            string
	Value           string
	HTTPCode        int
	CamelValue      string
	UpperCamelValue string
	LowerCamelValue string
	Key             string
	Comment         string
	HasComment      bool
	Pretty          string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
