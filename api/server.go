package main

import (
	"github.com/jcocozza/go_fractals/api/middleware"
	"github.com/jcocozza/go_fractals/api/route"
)

// Set up and run the API
func main() {
	router := route.CreateRouter()
	middleware.SetCors(router)
	route.InitRoutes(router)
	router.Run(":8080")
}