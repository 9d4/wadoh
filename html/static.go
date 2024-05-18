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
