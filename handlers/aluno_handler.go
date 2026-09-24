package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"api-gin/models"
	"github.com/gin-gonic/gin"
)

var alunos []models.Aluno
var proximoIDAluno = 1

func CriarAluno(c *gin.Context) {
	var aluno models.Aluno

	if err := c.ShouldBindJSON(&aluno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	aluno.Matricula = strings.TrimSpace(aluno.Matricula)
	aluno.Nome = strings.TrimSpace(aluno.Nome)
	aluno.Email = strings.TrimSpace(aluno.Email)

	if aluno.Matricula == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "A matrícula é obrigatória",
		})
		return
	}

	if aluno.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O nome é obrigatório",
		})
		return
	}

	if aluno.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O e-mail é obrigatório",
		})
		return
	}

	for _, existente := range alunos {
		if strings.EqualFold(existente.Matricula, aluno.Matricula) {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "Aluno já cadastrado com essa matrícula",
			})
			return
		}

		if strings.EqualFold(existente.Email, aluno.Email) {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "Aluno já cadastrado com esse e-mail",
			})
			return
		}
	}

	aluno.ID = proximoIDAluno

	alunos = append(alunos, aluno)
	proximoIDAluno++

	c.JSON(http.StatusCreated, aluno)
}

func ListarAlunos(c *gin.Context) {
	c.JSON(http.StatusOK, alunos)
}

func BuscarAluno(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	for _, aluno := range alunos {
		if aluno.ID == id {
			c.JSON(http.StatusOK, aluno)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Aluno não encontrado",
	})
}