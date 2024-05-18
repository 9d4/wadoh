package html

import (
	"net/http"
	"text/template"
)

type Templates struct {
	template *template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		template: template.Must(template.ParseFS(TemplatesFS(), "*.html")),
	}
}

func (t *Templates) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) error {
	return t.template.ExecuteTemplate(w, name, data)
}
