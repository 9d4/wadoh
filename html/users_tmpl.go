package html

import (
	"context"
	"html/template"
	"io"
	"io/fs"

	"github.com/9d4/wadoh/users"
)

type UsersTmpl struct {
	Rows *UsersRowsPartial
}

func (t *UsersTmpl) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmp := template.Must(base.
			ParseFS(fs, "pages/dashboard/users.html"),
		)

		return tmp.Execute(w, &pageData{
			Site: site,
			Page: t,
		})
	}

	return "dashboard.html", fn
}

type UsersRowsPartial struct {
	Users []users.User
}

func (p *UsersRowsPartial) Renderer(fs fs.FS, site *Site) (layout string, renderFn RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		t := template.Must(template.ParseFS(fs, "partials/users/rows.html"))
		return t.Execute(w, &pageData{
			Site: site,
			Page: p,
		})
	}

	return "", fn
}
