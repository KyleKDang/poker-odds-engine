// Package api provides HTTP server setup and routing.
package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns a Gin router.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	router.GET("/health", HandleHealth)
	router.POST("/evaluate", HandleEvaluate)
	router.POST("/odds", HandleOdds)

	return router
}
