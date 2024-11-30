package webapi

import (
	"links-sync-go/internal/config"

	"github.com/labstack/echo/v4"
)

type ApiServer struct {
	echoInstance *echo.Echo
	config       *config.Config
}

func NewApiServer(config *config.Config) *ApiServer {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Links sync here...")
	})

	server := &ApiServer{
		echoInstance: e,
		config:       config,
	}

	// setup handlers
	visitedGroup := e.Group("/get-links")
	visitedGroup.POST("", server.getVisitedListHandler)
	visitedGroup.GET("/:id", server.getVisitedHandler)
	visitedGroup.DELETE("/:id", server.deleteVisitedHandler)
	visitedGroup.PATCH("/:id", server.patchVisitedHandler)

	e.POST("/save-links", server.saveVisitedHandler)

	return server
}

func (s *ApiServer) Run() error {
	return s.echoInstance.Start(s.config.Server.Address)
}

// handlers

func (s *ApiServer) saveVisitedHandler(c echo.Context) error {
	return c.String(200, "save links stub")
}

func (s *ApiServer) getVisitedListHandler(c echo.Context) error {
	data := []visitedResponse{
		{
			Id:        111,
			Title:     "Some topic",
			PosterUrl: "some url",
			Status:    1,
		},
		{
			Id:        2,
			Title:     "topic 2",
			PosterUrl: "url 2",
			Status:    2,
		},
	}
	return c.JSON(200, data)
}

func (s *ApiServer) getVisitedHandler(c echo.Context) error {
	return c.String(200, "get one links")
}

func (s *ApiServer) deleteVisitedHandler(c echo.Context) error {
	return c.NoContent(204)
}

func (s *ApiServer) patchVisitedHandler(c echo.Context) error {
	return c.String(200, "patch visited")
}
