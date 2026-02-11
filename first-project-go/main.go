package main

import (
	"first-project-go/config"
	"first-project-go/routes"
)

func main() {
	// 1. Inisialisasi Database & NATS
	config.Init()

	// 2. Setup Router & Jalankan Server
	r := routes.SetupRouter()
	r.Run(":8080")
}
