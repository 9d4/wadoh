package html

import (
	"embed"
	"io/fs"
)

//go:embed templates/*
var templatesFs embed.FS

func TemplatesFS() fs.FS {
	sub, err := fs.Sub(templatesFs, "templates")
	if err != nil {
		panic(err)
	}
	return sub
}
