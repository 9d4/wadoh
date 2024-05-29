package html

import (
	"context"
	"html/template"
	"io"
	"io/fs"

	"github.com/9d4/wadoh/devices"
)

type DevicesTmpl struct {
	List   *DevicesListBlock
	Detail *DevicesDetailBlock
}

func (t *DevicesTmpl) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmp := template.Must(base.
			ParseFS(fs, "pages/dashboard/devices.html"),
		)

		return tmp.Execute(w, &pageData{
			Site: site,
			Page: t,
		})
	}

	return "dashboard.html", fn
}

type DevicesListBlock struct {
	Devices []devices.Device
}

func (p *DevicesListBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		t := template.Must(base.ParseFS(fs, "pages/dashboard/devices/block_list.html"))
		return t.Execute(w, &pageData{
			Site: site,
			Page: p,
		})
	}

	return "", fn
}

type DevicesDetailBlock struct {
	Device *devices.Device
}

func (p *DevicesDetailBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		t := template.Must(base.ParseFS(fs, "pages/dashboard/devices/block_detail.html"))
		return t.Execute(w, &pageData{
			Site: site,
			Page: p,
		})
	}

	return "", fn
}

type DevicesNewTmpl struct {
}

func (t *DevicesNewTmpl) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmp := template.Must(base.
			ParseFS(fs, "pages/dashboard/devices_new.html"),
		)

		return tmp.Execute(w, &pageData{
			Site: site,
			Page: t,
		})
	}

	return "dashboard.html", fn
}
