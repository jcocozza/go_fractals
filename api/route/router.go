package route

import (
	"github.com/gin-gonic/gin"
)

// Create a gin router
func CreateRouter() *gin.Engine {
	router := gin.Default()
	return router
}
