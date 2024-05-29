package html

import (
	"context"
	"html/template"
	"io"
	"io/fs"
)

type LoginTmpl struct {
}

func (lt *LoginTmpl) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		t := template.Must(base.
			ParseFS(fs, "pages/login.html"),
		)

		return t.Execute(w, &pageData{
			Site: site,
			Page: lt,
		})
	}

	return "base.html", fn
}
