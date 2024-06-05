package html

import (
	"context"
	"html/template"
	"io"
	"io/fs"
)

type ErrorTmpl struct {
	Message string
	Code    string
	Status  int
}

func (t *ErrorTmpl) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmpl, data, err := prepareRenderer(
			ctx, fs, site, base,
			"pages/error.html", t,
		)
		if err != nil {
			return err
		}
		return tmpl.Execute(w, data)
	}

	return "base.html", fn
}
