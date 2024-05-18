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

func (t *Templates) Render(w http.ResponseWriter, r *http.Request) error {
	return t.template.ExecuteTemplate(w, r.URL.Path, map[string]interface{}{"Title": "Hello World"})
}
