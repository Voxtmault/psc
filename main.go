package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogmulti "github.com/samber/slog-multi"
	"github.com/voxtmault/psc/config"
	"github.com/voxtmault/psc/db"
	logging "github.com/voxtmault/psc/logging"
	"github.com/voxtmault/psc/routers"
	"github.com/voxtmault/psc/validator"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.New("")

	// read env file
	validator.InitValidator()
	if err := db.InitConnection(&cfg.DBConfig); err != nil {
		panic(err)
	}

	if err := logging.InitLogger(&cfg.LoggingConfig); err != nil {
		panic(err)
	}

	logger := slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: false,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						a.Key = "time"
						a.Value = slog.StringValue(time.Now().Format(time.DateTime))
					}

					return a
				},
			}),
			slog.NewJSONHandler(logging.GetServerLogger(), &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						a.Key = "time"
						a.Value = slog.StringValue(time.Now().Format(time.DateTime))
					}

					return a
				},
			}),
		),
	)

	errLogger := slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						a.Key = "time"
						a.Value = slog.StringValue(time.Now().Format(time.DateTime))
					}

					return a
				},
			}),
			slog.NewJSONHandler(logging.GetErrorLogger(), &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						a.Key = "time"
						a.Value = slog.StringValue(time.Now().Format(time.DateTime))
					}

					return a
				},
			}),
		),
	)

	child := logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
		),
	)

	errChild := errLogger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
		),
	)

	e := echo.New()
	e.Validator = validator.GetEchoAdapter()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogURIPath:   true,
		LogHost:      true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogProtocol:  true,
		LogUserAgent: true,
		LogLatency:   true,
		LogError:     true,
		LogRequestID: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				child.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.Any("latency", v.Latency.String()),
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
					slog.String("protocol", v.Protocol),
					slog.String("URI", v.URI),
					slog.String("Path", v.URIPath),
					slog.String("Host", v.Host),
					slog.String("Client IP", v.RemoteIP),
					slog.String("User Agent", v.UserAgent),
				)
			} else {
				errChild.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.Any("latency", v.Latency.String()),
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.String("status", v.Protocol),
					slog.String("Path", v.URIPath),
					slog.String("Host", v.Host),
					slog.String("Client IP", v.RemoteIP),
					slog.String("User Agent", v.UserAgent),
					slog.String("err", v.Error.Error()),
				)
			}

			return nil
		},
	}))

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

	db.Close()
}
