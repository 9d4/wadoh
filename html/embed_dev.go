//go:build dev
// +build dev

package html

import (
	"io/fs"
	"os"
)

func StaticFs() fs.FS {
	return os.DirFS("html/static")
}

func TemplatesFS() fs.FS {
	return os.DirFS("html/templates")
}
