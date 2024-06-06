package html

import (
	"context"
	"html/template"
	"io"
	"io/fs"
	"maps"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/9d4/wadoh/users"
)

type ctxKey string

const (
	contextCtxKey ctxKey = "context"
	flashCtxKey   ctxKey = "flash"
	errorCtxKey   ctxKey = "error"
)

type Site struct {
	Title string
}

var SiteData = &Site{
	Title: "Wadoh",
}

// pageData is data that will be passed to template.
// this should be
type pageData struct {
	// Site represents site data.
	Site *Site
	// Page represents data from current Renderable or page.
	Page interface{}
	// Others are additional data that accessible from template.
	Others map[string]interface{}
}

func (d *pageData) set(k string, v interface{}) {
	if d.Others == nil {
		d.Others = make(map[string]interface{})
	}
	d.Others[k] = v
}

func (d *pageData) toMap() map[string]interface{} {
	// The usage should be
	// .Site.Title
	// .Page.I should be accessible using .Items
	// additional or Others simply .accessMe
	out := make(map[string]interface{})
	out["Site"] = d.Site
	out["Page"] = d.Page

	pageMap := make(map[string]interface{})
	mapstructure.Decode(d.Page, &pageMap)
	// merge
	for k, v := range pageMap {
		out[k] = v
	}
	for k, v := range d.Others {
		out[k] = v
	}

	return out
}

type RenderFunc func(ctx context.Context, base *template.Template, w io.Writer) error

type Renderable interface {
	Renderer(fs.FS, *Site) (layout string, fn RenderFunc)
}

// buildPageData creates map than can be used inside template.
func buildPageData(ctx context.Context, site *Site, data any) (map[string]interface{}, error) {
	out := structToMap(data)
	out["Page"] = maps.Clone(out)
	out["Site"] = site
	out["User"] = users.UserFromContext(ctx)
	out["Flash"] = FlashFromContext(ctx)
	if err := ErrorFromContext(ctx); err != nil {
		out["Error"] = err
	}

	return out, nil
}

// structToMap converts a struct to a map[string]interface{} while preserving nested struct types
func structToMap(obj interface{}) map[string]interface{} {
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)

	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
		objType = objType.Elem()
	}

	result := make(map[string]interface{})
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := objType.Field(i)

		// Handle only exported fields
		if fieldType.PkgPath != "" {
			continue
		}

		result[fieldType.Name] = field.Interface()
	}
	return result
}
