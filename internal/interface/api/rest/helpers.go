package rest

import (
	"bytes"
	"html/template"
)

func renderTemplate(tmpl *template.Template, name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
