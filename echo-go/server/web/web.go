package web

import (
	"echo-react-serve/config"
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TODO: change the go embed 'dist' if the folder name is different (like build)
// eg.: go:embed all:build

var (
	//go:embed all:dist
	dist embed.FS
	//go:embed dist/index.html
	indexHTML     embed.FS
	distDirFS     = echo.MustSubFS(dist, config.Envs.App.ClientBuildsFolder)
	distIndexHtml = echo.MustSubFS(indexHTML, config.Envs.App.ClientBuildsFolder)
)

func RegisterHandlers(e *echo.Echo) {
	e.FileFS("/", "index.html", distIndexHtml)
	e.StaticFS("/", distDirFS)

	// taken from: https://blog.stackademic.com/golang-with-vite-react-with-live-reload-9e929099752d
	// This is needed to serve the index.html file for all routes that are not /api/*
	// neede for SPA to work when loading a specific url directly
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			// Skip the proxy if the prefix is /api
			return len(c.Path()) >= 4 && c.Path()[:4] == "/api"
		},
		// Root directory from where the static content is served.
		Root: "/",
		// Enable HTML5 mode by forwarding all not-found requests to root so that
		// SPA (single-page application) can handle the routing.
		HTML5:      true,
		Browse:     false,
		IgnoreBase: true,
		Filesystem: http.FS(distDirFS),
	}))
}

// source: https://dev.to/pacholoamit/one-of-the-coolest-features-of-go-embed-reactjs-into-a-go-binary-41e9
