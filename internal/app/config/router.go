package config

import (
	"github.com/gin-gonic/gin"
	"seg-red-auth/internal/app/common"
	"seg-red-auth/internal/app/controller"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(common.GlobalErrorHandler())
	r.NoRoute(common.HandleNoRoute())

	v1 := r.Group("/api/v1")
	controller.NewAuthController(v1)
	return r
}
