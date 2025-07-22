package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/voxtmault/psc/routers"
	"github.com/voxtmault/psc/validator"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// read env file
	validator.InitValidator()

	e := echo.New()
	e.Validator = validator.GetEchoAdapter()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())

	e.HideBanner = false
	e.HidePort = false

	root := e.Group("/api/v1")
	if err := routers.KasusRoute(root); err != nil {
		panic(err)
	}

	// e.Logger.Fatal(e.Start(":8080"))

	server := &http.Server{
		Addr:    ":8080",
		Handler: e,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start http server", "error", err)
			return
		}
	}()

	<-ctx.Done()
	slog.Warn("shutting down http server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown http server", "error", err)
	} else {
		slog.Debug("successfully shutdown http server")
	}
}
