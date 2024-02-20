package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jcocozza/go_fractals/api/controllers"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {})

	r.POST("/julia", controllers.JuliaHandler)

}