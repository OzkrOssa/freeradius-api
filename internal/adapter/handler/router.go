package handler

import (
	"github.com/OzkrOssa/freeradius-api/internal/adapter/config"
	"github.com/OzkrOssa/freeradius-api/internal/core/port"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
)

type Router struct {
	*gin.Engine
}

func NewRouter(config *config.Http, token port.TokenService, authHandler *AuthHandler, userHandler *UserHandler) (*Router, error) {
	// Disable debug mode in production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	ginConfig := cors.DefaultConfig()
	ginConfig.AllowOrigins = config.AllowedOrigins

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))

	v1 := router.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.POST("/login", authHandler.Login)
		}

		authUser := user.Group("/").Use(authMiddleware(token))
		{
			authUser.GET("/", userHandler.ListUsers)
			authUser.POST("/:id", userHandler.GetUser)
			authUser.PUT("/:id", userHandler.UpdateUser)
			authUser.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	return &Router{
		router,
	}, nil

}

func (r *Router) RunServer(address string) error {
	return r.Run(address)
}
