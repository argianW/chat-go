package routes

import (
	"first-project-go/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine { // Fungsi untuk mengatur dan mengembalikan router Gin
	r := gin.Default() // Membuat instance router Gin dengan default middleware (logger dan recovery)

	// Serve Frontend
	r.StaticFile("/", "./public/index.html") // Menyajikan file statis index.html di root URL

	// API Groups
	api := r.Group("/api") // Membuat grup route untuk API dengan prefix /api
	{
		api.GET("/history/:myID/:targetID", handlers.GetHistory) // Route GET untuk mengambil history chat
		api.GET("/contacts/:myID", handlers.GetContacts)         // Route GET untuk mengambil daftar kontak
	}

	// WebSocket
	r.GET("/ws/:myID", handlers.HandleChat) // Route GET untuk koneksi WebSocket

	return r // Mengembalikan instance router yang sudah dikonfigurasi
}
