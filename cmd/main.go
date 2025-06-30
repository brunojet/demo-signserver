package main

import (
	"demo-signserver/internal/http_handler"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Rotas de perfil
	http_handler.RegisterProfileRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor iniciado na porta %s", port)
	r.Run(":" + port)
}
