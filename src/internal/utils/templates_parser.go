package utils

import (
	"bytes"
	"text/template"
)

func ParseTemplate(templateName string, data any) ([]byte, error) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
