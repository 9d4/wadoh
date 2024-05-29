package html

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"io/fs"
	"reflect"
	"strings"
)

type Templates struct {
	layouts map[string]*template.Template

	fs       fs.FS
	siteData *Site
}

func NewTemplates() *Templates {
	tmpl := template.Must(
		template.ParseFS(TemplatesFS(), "templates/*.html"),
	).Funcs(templateFuncs)

	empty := template.Must(
		tmpl.New("empty.html").
			ParseFS(TemplatesFS(), "layouts/empty.html"),
	)
	base := template.Must(
		tmpl.New("base.html").
			ParseFS(TemplatesFS(), "layouts/base.html"),
	)
	dash := template.Must(
		tmpl.New("dashboard.html").
			ParseFS(TemplatesFS(), "layouts/dashboard.html"),
	)

	layouts := make(map[string]*template.Template)
	layouts[""] = empty.Funcs(templateFuncs)
	layouts["empty.html"] = layouts[""]
	layouts["base.html"] = base.Funcs(templateFuncs)
	layouts["dashboard.html"] = dash.Funcs(templateFuncs)

	return &Templates{
		layouts: layouts,
		fs:      TemplatesFS(),
		siteData: &Site{
			Title: "Wadoh",
		},
	}
}

// R renders Renderable
func (t *Templates) R(ctx context.Context, w io.Writer, r Renderable) error {
	layout, render := r.Renderer(t.fs, t.siteData)
	var baseTmpl *template.Template
	if layoutTmpl, ok := t.layouts[layout]; ok {
		layoutTmpl, _ := layoutTmpl.Clone()
		layoutTmpl.Funcs(template.FuncMap{
			"render": fnHTMLRenderer(ctx, t),
		})
		baseTmpl = layoutTmpl
	}

	var buf bytes.Buffer
	if err := render(ctx, baseTmpl, &buf); err != nil {
		return err
	}
	_, err := buf.WriteTo(w)
	return err
}

// RenderPartial renders partial rawly.
func (t *Templates) RenderPartial(w io.Writer, name string, data interface{}) error {
	names := strings.Split(name, "/")
	partialPath := "partials/" + name
	page, err := template.New(names[len(names)-1]).
		Funcs(templateFuncs).
		ParseFS(TemplatesFS(), partialPath)
	if err != nil {
		return err
	}
	if err := page.Execute(w, data); err != nil {
		return err
	}
	return nil
}

var templateFuncs = template.FuncMap{}

func init() {
	templateFuncs["partial"] = fnPartial
	templateFuncs["devicePhone"] = devicePhone
}

func fnPartial(name string, data ...interface{}) template.HTML {
	part, err := template.ParseFS(TemplatesFS(), "partials/"+name)
	var out bytes.Buffer
	var dataToExec interface{}
	if len(data) > 0 {
		dataToExec = data[0]
	}
	if err != nil {
		return template.HTML(err.Error())
	}

	part = part.Funcs(templateFuncs)
	part.Execute(&out, dataToExec)
	return template.HTML(out.String())
}

// fnHTMLRenderer returns function that can be used to render
// Renderable from html template.
func fnHTMLRenderer(ctx context.Context, t *Templates) func(r Renderable) template.HTML {
	fn := func(r Renderable) template.HTML {
		if r == nil || reflect.ValueOf(r).IsNil() {
			return template.HTML("")
		}
		var buf bytes.Buffer
		if err := (*Templates).R(t, ctx, &buf, r); err != nil {
			return template.HTML(err.Error())
		}
		return template.HTML(buf.String())
	}

	return fn
}
