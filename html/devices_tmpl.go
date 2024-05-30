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
		tmpl, data, err := prepareRenderer(
			ctx, fs, site, base,
			"pages/dashboard/devices.html", t,
		)
		if err != nil {
			return err
		}
		return tmpl.Execute(w, data)
	}

	return "dashboard.html", fn
}

type DevicesListBlock struct {
	Devices []devices.Device
}

func (b *DevicesListBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmpl, data, err := prepareRenderer(
			ctx, fs, site, base,
			"pages/dashboard/devices/block_list.html", b,
		)
		if err != nil {
			return err
		}
		return tmpl.Execute(w, data)
	}

	return "", fn
}

type DevicesDetailBlock struct {
	Device *devices.Device
}

func (b *DevicesDetailBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmpl, data, err := prepareRenderer(
			ctx, fs, site, base,
			"pages/dashboard/devices/block_detail.html", b,
		)
		if err != nil {
			return err
		}
		return tmpl.Execute(w, data)
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

		data, err := buildPageData(ctx, site, t)
		if err != nil {
			return err
		}
		return tmp.Execute(w, data)
	}

	return "dashboard.html", fn
}
