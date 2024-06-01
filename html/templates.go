package html

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"reflect"

	"github.com/rs/zerolog/log"
)

type Templates struct {
	layouts map[string]*template.Template

	fs       fs.FS
	siteData *Site
}

func NewTemplates() *Templates {
	site := &Site{
		Title: "Wadoh",
	}
	fs := TemplatesFS()
	tmpl := template.Must(
		template.ParseFS(fs, "templates/*.html"),
	).
		Funcs(commonFuncs(fs, site))

	layouts := make(map[string]*template.Template)
	layoutNames := []string{"empty.html", "base.html", "dashboard.html"}

	for _, name := range layoutNames {
		t, err := tmpl.New(name).ParseFS(fs, "layouts/"+name)
		if err != nil {
			log.Fatal().Err(err).Str("layout", name).Msg("failed to parse layout")
		}
		layouts[name] = t
	}
	// fallback alias
	layouts[""] = layouts["empty.html"]

	return &Templates{
		layouts:  layouts,
		fs:       fs,
		siteData: site,
	}
}

// R renders Renderable
func (t *Templates) R(ctx context.Context, w io.Writer, r Renderable) error {
	layout, render := r.Renderer(t.fs, t.siteData)
	var baseTmpl *template.Template
	layoutTmpl, ok := t.layouts[layout]
	if !ok {
		return fmt.Errorf("templates: wanted layout not found %s", layout)
	}
	layoutTmpl, err := layoutTmpl.Clone()
	if err != nil {
		return fmt.Errorf("templates: error copying layout %s", layout)
	}

	layoutTmpl.Funcs(template.FuncMap{
		"render": fnHTMLRenderer(ctx, t),
	})
	baseTmpl = layoutTmpl

	var buf bytes.Buffer
	if err := render(ctx, baseTmpl, &buf); err != nil {
		return fmt.Errorf("templates: %w", err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return fmt.Errorf("templates: error writing buffer: %w", err)
	}
	return nil
}

// RenderPartial renders partial rawly.
func (t *Templates) RenderPartial(w io.Writer, name string, data interface{}) error {
	partialPath := "partials/" + name
	page, err := parseTemplate(t.fs, partialPath)
	// Funcs(templateFuncs).
	// ParseFS(TemplatesFS(), partialPath)
	if err != nil {
		return err
	}
	page.Funcs(commonFuncs(t.fs, t.siteData))

	if err := page.Execute(w, data); err != nil {
		return err
	}
	return nil
}

func parseTemplate(fs fs.FS, name string) (*template.Template, error) {
	tmpl, err := template.ParseFS(fs, name)
	if err != nil {
		return nil, fmt.Errorf("templates: unable to parse template: %w", err)
	}
	return tmpl, nil
}

// prepareRenderer is common helper that can be used by renderer
// to parse template and build template data
func prepareRenderer(
	ctx context.Context, fs fs.FS,
	site *Site, base *template.Template,
	name string, data any,
) (*template.Template, map[string]interface{}, error) {
	tmpl, err := base.ParseFS(fs, name)
	if err != nil {
		return nil, nil, err
	}
	tmpData, err := buildPageData(ctx, site, data)
	if err != nil {
		return nil, nil, err
	}
	return tmpl, tmpData, nil
}

func commonFuncs(fs fs.FS, s *Site) template.FuncMap {
	funcs := template.FuncMap{}
	funcs["partial"] = fnPartialRenderer(fs)
	funcs["title"] = fnTitle(s)
	return funcs
}

func fnPartialRenderer(fs fs.FS) func(string, interface{}) template.HTML {
	fn := func(name string, data any) template.HTML {
		part, err := parseTemplate(fs, "partials/"+name)
		if err != nil {
			panic(err)
		}

		var out bytes.Buffer
		if err := part.Execute(&out, data); err != nil {
			panic(fmt.Errorf("partial render error: %w", err))
		}
		return template.HTML(out.String())
	}

	return fn
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

// fnTitle makes possible to page title by combining Site.Title
// with current Page title if exists.
//
// this is useful to send new title on HTMX Ajax.
// See html/templates/pages/dashboard/devices/block_detail.html for example.
func fnTitle(s *Site) func(data map[string]interface{}) template.HTML {
	fn := func(data map[string]interface{}) template.HTML {
		title := s.Title
		page, ok := data["Page"].(map[string]interface{})
		if !ok {
			return template.HTML(title)
		}
		pageTitle, ok := page["Title"].(string)
		if !ok {
			return template.HTML(title)
		}

		title = fmt.Sprintf("%s - %s", pageTitle, title)
		return template.HTML(title)
	}

	return fn
}
