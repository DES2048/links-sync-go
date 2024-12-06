package webapi

import (
	"links-sync-go/internal/config"
	"links-sync-go/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type IApiServer interface {
	Run() error
	Storage() storage.IStorage
	Config() *config.Config
}

type ApiServer struct {
	echoInstance *echo.Echo
	config       *config.Config
	storage      storage.IStorage
}

func NewApiServer(config *config.Config) *ApiServer {
	e := echo.New()
	e.Debug = true

	/*
	   e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	       LogURI: true,
	       LogStatus: true,
	       LogError: true,
	   }))*/

	e.Use(middleware.CORS())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Token",
		Validator: func(auth string, c echo.Context) (bool, error) {
			return auth == config.Server.ApiKey, nil
		},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Links sync here...")
	})

	storage, _ := storage.NewDbStorage(config)

	server := &ApiServer{
		echoInstance: e,
		config:       config,
		storage:      storage,
	}

	// setup handlers
	visitedHandler := NewVisitedHandler(server)

	visitedGroup := e.Group("/get-links")
	visitedGroup.GET("", visitedHandler.getVisitedListHandler)
	visitedGroup.POST("", visitedHandler.getVisitedByIds)
	visitedGroup.GET("/:id", visitedHandler.getVisitedHandler)
	visitedGroup.DELETE("/:id", visitedHandler.deleteVisitedHandler)
	visitedGroup.PATCH("/:id", visitedHandler.patchVisitedHandler2)

	e.POST("/save-links", visitedHandler.saveVisitedHandler)

	return server
}

func (s *ApiServer) Run() error {
	return s.echoInstance.Start(s.config.Server.Address)
}

func (s *ApiServer) Storage() storage.IStorage {
	return s.storage
}

func (s *ApiServer) Config() *config.Config {
	return s.config
}
