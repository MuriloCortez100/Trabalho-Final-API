package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"api-gin/models"
	"github.com/gin-gonic/gin"
)

var salas []models.Sala
var proximoIDSala = 1

func CriarSala(c *gin.Context) {
	var sala models.Sala

	if err := c.ShouldBindJSON(&sala); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	sala.Nome = strings.TrimSpace(sala.Nome)

	if sala.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O nome da sala é obrigatório",
		})
		return
	}

	if sala.Capacidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "A capacidade deve ser maior que zero",
		})
		return
	}

	for _, existente := range salas {
		if strings.EqualFold(existente.Nome, sala.Nome) {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "Já existe uma sala com esse nome",
			})
			return
		}
	}

	if sala.Recursos == nil {
		sala.Recursos = []string{}
	}

	sala.ID = proximoIDSala
	sala.Ativa = true

	salas = append(salas, sala)
	proximoIDSala++

	c.JSON(http.StatusCreated, sala)
}

func ListarSalas(c *gin.Context) {
	c.JSON(http.StatusOK, salas)
}

func BuscarSala(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	for _, sala := range salas {
		if sala.ID == id {
			c.JSON(http.StatusOK, sala)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Sala não encontrada",
	})
}

func AlterarStatusSala(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	var dados struct {
		Ativa bool `json:"ativa"`
	}

	if err := c.ShouldBindJSON(&dados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	for i := range salas {
		if salas[i].ID == id {
			salas[i].Ativa = dados.Ativa

			c.JSON(http.StatusOK, gin.H{
				"mensagem": "Status da sala alterado com sucesso",
				"sala":     salas[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Sala não encontrada",
	})
}

func VerificarDisponibilidadeSala(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	diaSemana := normalizarDia(c.Query("dia"))
	horaInicio := c.Query("hora_inicio")
	horaFim := c.Query("hora_fim")

	if diaSemana == "" || horaInicio == "" || horaFim == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Informe dia, hora_inicio e hora_fim",
		})
		return
	}

	if !validarHorario(horaInicio) || !validarHorario(horaFim) {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Os horários devem estar no formato HH:MM",
		})
		return
	}

	if horaInicio >= horaFim {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O horário de início deve ser anterior ao horário de fim",
		})
		return
	}

	salaEncontrada := -1

	for i, sala := range salas {
		if sala.ID == id {
			salaEncontrada = i
			break
		}
	}

	if salaEncontrada == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "Sala não encontrada",
		})
		return
	}

	if !salas[salaEncontrada].Ativa {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "A sala está inativa",
		})
		return
	}

	for _, turma := range turmas {
		if turma.SalaID == nil {
			continue
		}

		if *turma.SalaID != id {
			continue
		}

		if !turma.Ativa {
			continue
		}

		if turma.DiaSemana != diaSemana {
			continue
		}

		if horaInicio < turma.HoraFim && horaFim > turma.HoraInicio {
			c.JSON(http.StatusOK, gin.H{
				"disponivel": false,
				"mensagem":   "A sala está ocupada nesse horário",
				"turma":      montarTurmaResposta(turma),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"disponivel":  true,
		"mensagem":    "A sala está disponível nesse horário",
		"sala":        salas[salaEncontrada],
		"dia":         diaSemana,
		"hora_inicio": horaInicio,
		"hora_fim":    horaFim,
	})
}

func validarHorario(horario string) bool {
	_, err := time.Parse("15:04", horario)
	return err == nil
}

func normalizarDia(dia string) string {
	dia = strings.ToLower(strings.TrimSpace(dia))

	switch dia {
	case "segunda", "segunda-feira":
		return "segunda"
	case "terca", "terça", "terca-feira", "terça-feira":
		return "terca"
	case "quarta", "quarta-feira":
		return "quarta"
	case "quinta", "quinta-feira":
		return "quinta"
	case "sexta", "sexta-feira":
		return "sexta"
	case "sabado", "sábado":
		return "sabado"
	case "domingo":
		return "domingo"
	default:
		return ""
	}
}