package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}
var {{.LowerCamelValue}} *apierrors.Error
{{- end }}

func init() {
{{- range .Errors }}
{{.LowerCamelValue}} = apierrors.New({{.HTTPCode}}, "{{.Key}}", "{{.Pretty}}", {{.Msg}})
apierrors.Register({{.LowerCamelValue}})
{{- end }}
}

{{ range .Errors }}
{{if .HasComment}}{{.Comment}}{{end}}func {{.UpperCamelValue}}() *apierrors.Error {
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
	Msg             string
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
