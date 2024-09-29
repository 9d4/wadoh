package html

import (
	"context"
	"html/template"
	"io"
	"io/fs"

	"github.com/9d4/wadoh/devices"
)

type DevicesTmpl struct {
	Title  string
	List   *DevicesListBlock
	Detail *DevicesDetailBlock
}

func (t *DevicesTmpl) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	t.Title = "Devices"
	if t.Detail != nil {
		t.Title = t.Detail.Device.Name
	}

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
	Title   string
	Devices []devices.Device
}

func (b *DevicesListBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	b.Title = "Devices"

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
	Title      string
	Device     *devices.Device
	DetailPane *DevicesDetailPaneBlock
}

func (b *DevicesDetailBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	b.Title = b.Device.Name

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

type DevicesNewTmpl struct{}

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

type DevicesReconnectTmpl struct {
	Device *devices.Device
}

func (t *DevicesReconnectTmpl) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		tmp := template.Must(base.
			ParseFS(fs, "pages/dashboard/devices_reconnect.html"),
		)

		data, err := buildPageData(ctx, site, t)
		if err != nil {
			return err
		}
		return tmp.Execute(w, data)
	}

	return "dashboard.html", fn
}

type DevicesDetailPaneBlock struct {
	Device *devices.Device

	SubAPIKey     bool
	SubTryMessage bool
	SubWebhook    bool
	SubMore       bool
}

func (t *DevicesDetailPaneBlock) Renderer(fs fs.FS, site *Site) (string, RenderFunc) {
	fn := func(ctx context.Context, base *template.Template, w io.Writer) error {
		// set default pane if none true
		if !t.SubAPIKey && !t.SubTryMessage && !t.SubWebhook && !t.SubMore {
			t.SubAPIKey = true
		}

		tmp := template.Must(base.
			ParseFS(fs, "pages/dashboard/devices/block_detail_pane.html"),
		)

		data, err := buildPageData(ctx, site, t)
		if err != nil {
			return err
		}
		return tmp.Execute(w, data)
	}

	return "", fn
}
