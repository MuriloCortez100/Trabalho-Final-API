package main

import (
	"net/http"
	"time"

	"api-gin/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})

		v1.POST("/salas", handlers.CriarSala)
		v1.GET("/salas", handlers.ListarSalas)
		v1.GET("/salas/:id", handlers.BuscarSala)
		v1.PUT("/salas/:id/status", handlers.AlterarStatusSala)
		v1.GET("/salas/:id/disponibilidade", handlers.VerificarDisponibilidadeSala)

		v1.POST("/alunos", handlers.CriarAluno)
		v1.GET("/alunos", handlers.ListarAlunos)
		v1.GET("/alunos/:id", handlers.BuscarAluno)

		v1.POST("/turmas", handlers.CriarTurma)
		v1.GET("/turmas", handlers.ListarTurmas)
		v1.GET("/turmas/:id", handlers.BuscarTurma)
		v1.PUT("/turmas/:id/status", handlers.AlterarStatusTurma)
		v1.POST("/turmas/:id/alunos", handlers.MatricularAluno)
		v1.GET("/turmas/:id/alunos", handlers.ListarAlunosDaTurma)
		v1.POST("/turmas/:id/alocar", handlers.AlocarSala)
	}

	r.Run(":8080")
}