package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


func main() {
	logger := createLogger()
	zap.ReplaceGlobals(logger)
	defer logger.Sync() 

	zap.L().Info("app starting...")

	app := fiber.New()

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler())) // monitoring

	app.Get("/", func(c *fiber.Ctx) error {
		zap.L().Info("Hello World!")
		return c.SendString("Hello World")
	})




	// Start server as go routine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	zap.L().Info("Server started on port 3000")

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

func createLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config {
		Level: zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		DisableCaller: false,
		DisableStacktrace: false,
		Sampling: nil,
		Encoding: "json",
		EncoderConfig: encoderCfg,
		OutputPaths: []string {
			"stderr",
		},
		ErrorOutputPaths: []string {
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getegid(),
		},
	}

	return zap.Must(config.Build())
}
