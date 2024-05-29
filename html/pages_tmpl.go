package html

import (
	"context"
	"html/template"
	"io"
	"io/fs"

	"github.com/9d4/wadoh/users"
)

type Site struct {
	Title string
}

var SiteData = &Site{
	Title: "Wadoh",
}

// Ctx is request context.
type Ctx struct {
	User *users.User
}

type pageData struct {
	Site *Site
	Ctx  *Ctx
	Page interface{}
}

type RenderFunc func(ctx context.Context, base *template.Template, w io.Writer) error
type Renderable interface {
	// RenderData() *RenderData
	Renderer(fs.FS, *Site) (layout string, fn RenderFunc)
}

