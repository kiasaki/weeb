package weeb

import (
	"bytes"
	"html/template"
)

type Templates interface {
	Render(name string, value J) (string, error)
	Add(name, contents string) error
}

type TemplatesGo struct {
	templates *template.Template
}

var _ Templates = Templates(&TemplatesGo{})

func NewTemplatesGo(templates *template.Template) *TemplatesGo {
	return &TemplatesGo{templates: templates}
}

func (t *TemplatesGo) Render(name string, value J) (string, error) {
	var b bytes.Buffer
	err := t.templates.ExecuteTemplate(&b, name, value)
	return b.String(), err
}

func (t *TemplatesGo) Add(name, contents string) error {
	templates, err := t.templates.New(name).Parse(contents)
	if err != nil {
		return err
	}
	t.templates = templates
	return nil
}
