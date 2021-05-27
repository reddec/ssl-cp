package ui

import (
	"embed"
	"net/http"
)

//go:embed app/dist
var assets embed.FS

const Path = "app/dist"

func Handler() http.Handler {
	return http.FileServer(http.FS(assets))
}
