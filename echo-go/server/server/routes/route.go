package routes

import (
	"echo-react-serve/config"
	"echo-react-serve/server/controllers"
	"echo-react-serve/server/repositories"
	"echo-react-serve/server/services"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func route(e *echo.Echo) {
	repo := repositories.NewRepo(config.MongoDB)
	service := services.NewService(repo)
	handler := controllers.NewController(service)
	g := e.Group("/group")
	g.GET("", handler.Handle)
}

func ApiProxyHandler(c echo.Context) error {
	if config.Envs.App.ProtectAPI {

	}

	// Forward the request
	req := c.Request()
	client := &http.Client{}
	proxyReq, err := http.NewRequest(req.Method, config.Envs.App.RealBackendUrl+req.URL.Path, req.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	proxyReq.Header = req.Header
	resp, err := client.Do(proxyReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Internal error while contacting backend\n" + err.Error(),
		})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
