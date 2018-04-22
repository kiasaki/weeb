package weeb

import (
	"bytes"
	"html/template"
	"time"
)

var TemplatesFunctionMap = template.FuncMap{
	"date":        func(t time.Time) string { return t.Format("2006-01-02") },
	"datetime":    func(t time.Time) string { return t.Format("2006-01-02 15:04") },
	"datetimesec": func(t time.Time) string { return t.Format("2006-01-02 15:04:05") },
}

type Templates interface {
	Render(name string, value J) (string, error)
	Add(name, contents string) error
}

type TemplatesGo struct {
	t       *template.Template
	funcMap template.FuncMap
}

var _ Templates = Templates(&TemplatesGo{})

func NewTemplatesGo(templates *template.Template) *TemplatesGo {
	templates.Funcs(TemplatesFunctionMap)
	return &TemplatesGo{t: templates, funcMap: funcMap}
}

func (t *TemplatesGo) Render(name string, value J) (string, error) {
	var b bytes.Buffer
	err := t.t.ExecuteTemplate(&b, name, value)
	return b.String(), err
}

func (t *TemplatesGo) Add(name, contents string) error {
	templates, err := t.t.New(name).Parse(contents)
	if err != nil {
		return err
	}
	t.t = templates
	return nil
}

func (t *TemplatesGo) AddFunction(name string, fn interface{}) {
	t.funcMap[name] = fn
	t.t.Funcs(t.funcMap)
}
