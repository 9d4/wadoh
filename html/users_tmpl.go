package html

import (
	"context"
	"html/template"
	"io"
	"io/fs"

	"github.com/9d4/wadoh/users"
)

type UsersTmpl struct {
	List   *UsersListBlock
	Detail *UsersDetailBlock
}

func (t *UsersTmpl) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmp, data, err := prepareRenderer(
			ctx, fs, site, base,
			"pages/dashboard/users.html", t,
		)
		if err != nil {
			return err
		}

		return tmp.Execute(w, data)
	}

	return "dashboard.html", fn
}

type UsersListBlock struct {
	Rows *UsersRowsBlock
}

func (t *UsersListBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmp, data, err := prepareRenderer(
			ctx, fs, site, base,
			"pages/dashboard/users/block_list.html", t,
		)
		if err != nil {
			return err
		}

		return tmp.Execute(w, data)
	}

	return "", fn
}

type UsersRowsBlock struct {
	Users []users.User
}

func (p *UsersRowsBlock) Renderer(fs fs.FS, site *Site) (layout string, renderFn RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmp, data, err := prepareRenderer(
			ctx, fs, site, base,
			"pages/dashboard/users/block_rows.html", p,
		)
		if err != nil {
			return err
		}

		return tmp.Execute(w, data)
	}

	return "", fn
}

type UsersDetailBlock struct {
	User *users.User
}

func (t *UsersDetailBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmp, data, err := prepareRenderer(
			ctx, fs, site, base,
			"pages/dashboard/users/block_detail.html", t,
		)
		if err != nil {
			return err
		}

		return tmp.Execute(w, data)
	}

	return "", fn
}
