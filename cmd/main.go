package main

import (
	"demo-signserver/internal/request/http_handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Rotas de perfil
	http_handlers.RegisterProfileRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor iniciado na porta %s", port)
	r.Run(":" + port)
}
