package web

import (
	"github.com/MicBun/go-100-coverage-docker-crud/service"
	"github.com/MicBun/go-100-coverage-docker-crud/web/handlers"
	"github.com/MicBun/go-100-coverage-docker-crud/web/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterAPIRoutes(c *service.Container) {
	api := handlers.NewApiHandler(c)

	c.Web.GET("/hello", api.Hello)
	c.Web.POST("/login", api.Login)

	userRoutes := c.Web.Group("/user")
	userRoutes.Use(middleware.JwtAuthMiddleware())
	userRoutes.POST("/register", api.RegisterUser)
	userRoutes.PUT("/update/:id", api.UpdateUser)
	userRoutes.DELETE("/delete/:id", api.DeleteUser)
	userRoutes.GET("/get/:id", api.GetUserByID)
	userRoutes.GET("/get", api.GetUserByToken)
	userRoutes.GET("/list", api.ListUsers)
	userRoutes.GET("/refresh", api.RefreshToken)

	c.Web.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
