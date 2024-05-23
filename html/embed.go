//go:build !dev

package html

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var staticFs embed.FS

func StaticFs() fs.FS {
	sub, err := fs.Sub(staticFs, "static")
	if err != nil {
		panic(err)
	}
	return sub
}

//go:embed templates/*
var templatesFs embed.FS

func TemplatesFS() fs.FS {
	sub, err := fs.Sub(templatesFs, "templates")
	if err != nil {
		panic(err)
	}
	return sub
}
