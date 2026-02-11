package routes

import (
	"first-project-go/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Serve Frontend
	r.StaticFile("/", "./public/index.html")

	// API Groups
	api := r.Group("/api")
	{
		api.GET("/history/:myID/:targetID", handlers.GetHistory)
		api.GET("/contacts/:myID", handlers.GetContacts)
	}

	// WebSocket
	r.GET("/ws/:myID", handlers.HandleChat)

	return r
}