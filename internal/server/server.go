package server

import (
	"context"
	"encoding/gob"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/config"
	"github.com/sullivtr/k8s_platform/internal/handlers"
	"github.com/sullivtr/k8s_platform/internal/providers"
	"github.com/sullivtr/k8s_platform/internal/server/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	gob.Register(time.Time{})
}

// Server represents the primary entry point for this service.
// It provides a web server containing a collection of REST controllers used by client side
// interactions.
//
// The web server also serves up the static client front-end.
type Server struct {
	*echo.Echo
}

// NewApp will create a new Server using the provider configuration to run the khub application.
func NewApp(c *config.Config) Server {
	e := Server{Echo: echo.New()}
	docs.SwaggerInfo.Version = c.Version
	prvds := &providers.ModuleProviders{
		Config: c,
	}

	prvds.InitCacheProvider()
	prvds.InitStorageProvider()
	prvds.InitK8sProvider()
	prvds.InitAWSProvider()

	e.Use(getMiddleware(c, prvds)...)

	if err := handlers.RegisterRoutes(e.Echo, prvds); err != nil {
		log.Fatal().Msgf("unable to register route handlers: %v", err)
	}

	// Add the swagger docs
	if c.IsProduction() {
		e.GET("/swagger/doc.json", echoSwagger.WrapHandler)
	} else {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	return e
}

// NewDataSink will create a new data sink server using the provider configuration.
func NewDataSink(c *config.Config) {
	prvds := &providers.ModuleProviders{
		Config: c,
	}
	prvds.InitCacheProvider()
	prvds.InitStorageProvider()
	prvds.InitK8sProvider()

	log.Info().Msg("Starting data sink server")
	prvds.StartDataSink(context.Background(), c.K8sDataSinkIntervalSeconds)
}
