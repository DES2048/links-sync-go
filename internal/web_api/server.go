package webapi

import (
	"context"
	"log/slog"
	"net/http"

	"links-sync-go/internal/config"
	"links-sync-go/internal/storage"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type IApiServer interface {
	Run() error
	Storage() storage.IStorage
	Config() *config.Config
	Echo() *echo.Echo
}

type ApiServer struct {
	echoInstance *echo.Echo
	config       *config.Config
	storage      storage.IStorage
}

func NewApiServer(config *config.Config) *ApiServer {
	e := echo.New()
	e.Debug = true

	// error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		err2 := err
		if storage.IsErrNotFound(err) {
			err2 = echo.NewHTTPError(http.StatusNotFound, err)
		}

		if validErr, ok := err.(validation.Errors); ok {
			err2 = echo.NewHTTPError(http.StatusBadRequest, validErr)
		}

		e.DefaultHTTPErrorHandler(err2, c)
	}

	storage, _ := storage.NewDbStorage(config)

	server := &ApiServer{
		echoInstance: e,
		config:       config,
		storage:      storage,
	}
	server.SetupMiddleware()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Links sync here...")
	})

	// setup handlers

	visitedGroup := e.Group("/get-links")
	visitedGroup.GET("", getVisitedListHandler(storage.VisitedRepo()))
	visitedGroup.POST("", getVisitedByIds(storage.VisitedRepo()))
	visitedGroup.GET("/:id", getVisitedHandler(storage.VisitedRepo()))
	visitedGroup.DELETE("/:id", deleteVisitedHandler(storage.VisitedRepo()))
	visitedGroup.PATCH("/:id", patchVisitedHandler2(storage.VisitedRepo()))

	e.POST("/save-links", saveVisitedHandler(storage.VisitedRepo()))

	return server
}

func (s *ApiServer) SetupMiddleware() {
	e := s.echoInstance
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogMethod:        true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logAttrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("bytes_in", v.ContentLength),
				slog.Int64("bytes_out", v.ResponseSize),
			}

			if v.Error != nil {
				logAttrs = append(
					logAttrs,
					slog.String("err", v.Error.Error()),
				)
			}
			slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST", logAttrs...)

			return nil
		},
	}))

	e.Use(middleware.CORS())

	// auth
	if s.config.Server.ApiKey != "" {
		e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			AuthScheme: "Token",
			Validator: func(auth string, c echo.Context) (bool, error) {
				return auth == s.config.Server.ApiKey, nil
			},
		}))
	}
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

func (s *ApiServer) Echo() *echo.Echo {
	return s.echoInstance
}
