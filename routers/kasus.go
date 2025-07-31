package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/voxtmault/psc/controllers"
	"github.com/voxtmault/psc/db"
	projectMiddleware "github.com/voxtmault/psc/middleware"
	"github.com/voxtmault/psc/services"
)

// /api/v1
func KasusRoute(root *echo.Group) error {
	route := root.Group("/kasus")

	service := services.NewKasusService(db.GetDBCon(), nil)
	controllers := controllers.NewKasusController(service)

	route.Use(middleware.RemoveTrailingSlash())
	route.Use(projectMiddleware.VerifyToken)
	route.GET("", controllers.Get3)
	route.POST("", controllers.Create)

	return nil
}

// flow normal
// router -> controller -> service

// flow normal
// router -> middleware -> controller -> service
