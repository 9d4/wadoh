package html

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/9d4/wadoh/users"
	"github.com/rs/zerolog/log"
)

type Templates struct {
	base *template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		base: template.Must(
			template.ParseFS(TemplatesFS(), "layouts/base.html"),
		).Funcs(templateFuncs),
	}
}

// Render searches page template in pages/ directory. It also applies base template
// based on sub directory. Example of name dashboard/index.html will use base template
// from layouts/dashboard.html if exists or fallback to default layouts/base.html.
func (t *Templates) Render(w http.ResponseWriter, r *http.Request, user *users.User, name string, data interface{}) error {
	name = strings.TrimPrefix(name, "/")
	base, _ := t.base.Clone()
	tmplData := &TemplateData{
		Layout: strings.SplitN(base.Name(), ".", 2)[0],
		Data:   data,
		User:   user,
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

func (t *Templates) RenderPartial(w http.ResponseWriter, name string, data *PartialData) {
	names := strings.Split(name, "/")
	partialPath := "partials/" + name
	page, err := template.New(names[len(names)-1]).
		Funcs(templateFuncs).
		ParseFS(TemplatesFS(), partialPath)
	if err != nil {
		http.Error(w, "unable to render view", http.StatusInternalServerError)
		log.Debug().Caller().Err(err).Send()
		return
	}
	if err := page.Execute(w, map[string]interface{}(*data)); err != nil {
		http.Error(w, "unable to render view", http.StatusInternalServerError)
		log.Debug().Caller().Err(err).Send()
		return
	}
}

var templateFuncs = template.FuncMap{}

func init() {
	templateFuncs["partial"] = partial
	templateFuncs["devicePhone"] = devicePhone
}

func partial(name string, data ...interface{}) template.HTML {
	part, err := template.ParseFS(TemplatesFS(), "partials/"+name)
	var out bytes.Buffer
	var dataToExec interface{}
	if len(data) > 0 {
		dataToExec = data[0]
	}

	if err != nil {
		template.New("-").Funcs(templateFuncs).Execute(&out, dataToExec)
		return template.HTML(out.String())
	}
	part = part.Funcs(templateFuncs)
	part.Execute(&out, dataToExec)
	return template.HTML(out.String())
}

type TemplateData struct {
	Layout string
	Data   interface{}
	User   *users.User
}

type PartialData map[string]interface{}

func NewPartialData() *PartialData {
	return &PartialData{}
}

func (d *PartialData) Set(key string, value interface{}) *PartialData {
	dMap := map[string]interface{}(*d)
	dMap[key] = value
	return (*PartialData)(&dMap)
}
