package browser

import (
	"embed"
	"io/fs"
)

//go:embed production/*
var assets embed.FS

func GetStaticAssets() fs.FS {
	return assets
}
