package html

import (
	"html/template"
	"net/http"
	"strings"
)

type Templates struct {
	base *template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		base: template.Must(template.ParseFS(TemplatesFS(), "layouts/base.html")),
	}
}

// Render searches page template in pages/ directory. It also applies base template
// based on sub directory. Example of name dashboard/index.html will use base template
// from layouts/dashboard.html if exists or fallback to default layouts/base.html.
func (t *Templates) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) error {
	name = strings.TrimPrefix(name, "/")
	base, _ := t.base.Clone()
	tmplData := &TemplateData{
		Layout: strings.SplitN(base.Name(), ".", 2)[0],
		Data:   data,
	}

	page, err := base.ParseFS(TemplatesFS(), "pages/"+name)
	if err != nil {
		return err
	}

	// Get layout based on directory before the page template html.
	// If layout is defined in the layouts/ directory then parse it.
	// Parsing layout in last order ensures that the page template
	// does not replace the definitions inside the extended layout template.
	layoutString := strings.SplitN(name, "/", 2)
	if len(layoutString) == 2 {
		_, err := page.ParseFS(TemplatesFS(), "layouts/"+layoutString[0]+".html")
		if err == nil {
			tmplData.Layout = layoutString[0]
		}
	}

	return page.Execute(w, tmplData)
}

type TemplateData struct {
	Layout string
	Data   interface{}
}
