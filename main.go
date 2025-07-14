package main

import (
	"arslanoktay/denemeler/pkg/config"
	_ "arslanoktay/denemeler/pkg/log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)





func main() {

	appConfig := config.Read()


	defer zap.L().Sync() 

	zap.L().Info("app starting...")

	app := fiber.New()

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		// TODO: check some dependencies
		return c.SendString("OK")
	})

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler())) // monitoring

	app.Get("/", func(c *fiber.Ctx) error {
		zap.L().Info("Hello World!")
		return c.SendString("Hello World")
	})

	// Start server as go routine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", appConfig.Port)); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	zap.L().Info("Server started on port", zap.String("port", appConfig.Port))

	gracefulShutdown(app)

}

func gracefulShutdown(app *fiber.App) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait shutdown signal
	<-sigChan
	zap.L().Info("Shutting down server...")

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Error("Error during server shutdown", zap.Error(err))
	}


	zap.L().Info("Server gracefully stopped.")
}
