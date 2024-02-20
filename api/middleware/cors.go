package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Create cors handler
func createCors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	return cors.New(config)
}

// Set cors for the router
func SetCors(router *gin.Engine) {
	c := createCors()
	router.Use(c)
}