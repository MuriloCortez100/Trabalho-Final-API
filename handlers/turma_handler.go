package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"api-gin/models"
	"github.com/gin-gonic/gin"
)

var turmas []models.Turma
var proximoIDTurma = 1

type TurmaResposta struct {
	ID               int          `json:"id"`
	Nome             string       `json:"nome"`
	Disciplina       string       `json:"disciplina"`
	Professor        string       `json:"professor"`
	QuantidadeAlunos int          `json:"quantidade_alunos"`
	Alunos           []int        `json:"alunos"`
	Ativa            bool         `json:"ativa"`
	SalaAlocada      bool         `json:"sala_alocada"`
	Sala             *models.Sala `json:"sala,omitempty"`
	DiaSemana        string       `json:"dia_semana"`
	HoraInicio       string       `json:"hora_inicio"`
	HoraFim          string       `json:"hora_fim"`
}

func montarTurmaResposta(turma models.Turma) TurmaResposta {
	resposta := TurmaResposta{
		ID:               turma.ID,
		Nome:             turma.Nome,
		Disciplina:       turma.Disciplina,
		Professor:        turma.Professor,
		QuantidadeAlunos: len(turma.Alunos),
		Alunos:           turma.Alunos,
		Ativa:            turma.Ativa,
		SalaAlocada:      turma.SalaID != nil,
		DiaSemana:        turma.DiaSemana,
		HoraInicio:       turma.HoraInicio,
		HoraFim:          turma.HoraFim,
	}

	if turma.SalaID != nil {
		for _, sala := range salas {
			if sala.ID == *turma.SalaID {
				salaCopia := sala
				resposta.Sala = &salaCopia
				break
			}
		}
	}

	return resposta
}

func CriarTurma(c *gin.Context) {
	var turma models.Turma

	if err := c.ShouldBindJSON(&turma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	turma.Nome = strings.TrimSpace(turma.Nome)
	turma.Disciplina = strings.TrimSpace(turma.Disciplina)
	turma.Professor = strings.TrimSpace(turma.Professor)

	if turma.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O nome da turma é obrigatório",
		})
		return
	}

	if turma.Disciplina == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "A disciplina é obrigatória",
		})
		return
	}

	if turma.Professor == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O professor é obrigatório",
		})
		return
	}

	for _, existente := range turmas {
		if strings.EqualFold(existente.Nome, turma.Nome) {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "Já existe uma turma com esse nome",
			})
			return
		}
	}

	turma.ID = proximoIDTurma
	turma.Alunos = []int{}
	turma.SalaID = nil
	turma.DiaSemana = ""
	turma.HoraInicio = ""
	turma.HoraFim = ""
	turma.Ativa = true

	turmas = append(turmas, turma)
	proximoIDTurma++

	c.JSON(http.StatusCreated, montarTurmaResposta(turma))
}

func ListarTurmas(c *gin.Context) {
	respostas := make([]TurmaResposta, 0, len(turmas))

	for _, turma := range turmas {
		respostas = append(respostas, montarTurmaResposta(turma))
	}

	c.JSON(http.StatusOK, respostas)
}

func BuscarTurma(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	for _, turma := range turmas {
		if turma.ID == id {
			c.JSON(http.StatusOK, montarTurmaResposta(turma))
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Turma não encontrada",
	})
}

func AlterarStatusTurma(c *gin.Context) {
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

	for i := range turmas {
		if turmas[i].ID == id {
			turmas[i].Ativa = dados.Ativa

			c.JSON(http.StatusOK, gin.H{
				"mensagem": "Status da turma alterado com sucesso",
				"turma":    montarTurmaResposta(turmas[i]),
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Turma não encontrada",
	})
}

func MatricularAluno(c *gin.Context) {
	turmaID, err := strconv.Atoi(c.Param("id"))

	if err != nil || turmaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID da turma inválido",
		})
		return
	}

	var dados struct {
		AlunoID int `json:"aluno_id"`
	}

	if err := c.ShouldBindJSON(&dados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	if dados.AlunoID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O aluno_id deve ser maior que zero",
		})
		return
	}

	turmaEncontrada := -1

	for i, turma := range turmas {
		if turma.ID == turmaID {
			turmaEncontrada = i
			break
		}
	}

	if turmaEncontrada == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "Turma não encontrada",
		})
		return
	}

	if !turmas[turmaEncontrada].Ativa {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "A turma está inativa",
		})
		return
	}

	alunoExiste := false

	for _, aluno := range alunos {
		if aluno.ID == dados.AlunoID {
			alunoExiste = true
			break
		}
	}

	if !alunoExiste {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "Aluno não encontrado",
		})
		return
	}

	for _, alunoID := range turmas[turmaEncontrada].Alunos {
		if alunoID == dados.AlunoID {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "Aluno já está matriculado nesta turma",
			})
			return
		}
	}

	if turmas[turmaEncontrada].SalaID != nil {
		salaID := *turmas[turmaEncontrada].SalaID

		salaEncontrada := -1

		for i, sala := range salas {
			if sala.ID == salaID {
				salaEncontrada = i
				break
			}
		}

		if salaEncontrada == -1 {
			c.JSON(http.StatusNotFound, gin.H{
				"erro": "A sala alocada para a turma não foi encontrada",
			})
			return
		}

		if !salas[salaEncontrada].Ativa {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "A sala alocada para a turma está inativa",
			})
			return
		}

		if len(turmas[turmaEncontrada].Alunos) >= salas[salaEncontrada].Capacidade {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"erro": "A capacidade da sala foi atingida",
			})
			return
		}
	}

	novaTurma := turmas[turmaEncontrada]

	if novaTurma.DiaSemana != "" &&
		novaTurma.HoraInicio != "" &&
		novaTurma.HoraFim != "" {

		for _, turma := range turmas {
			if turma.ID == turmaID {
				continue
			}

			if !turma.Ativa {
				continue
			}

			alunoNaTurma := false

			for _, alunoID := range turma.Alunos {
				if alunoID == dados.AlunoID {
					alunoNaTurma = true
					break
				}
			}

			if !alunoNaTurma {
				continue
			}

			if turma.DiaSemana == "" ||
				turma.HoraInicio == "" ||
				turma.HoraFim == "" {
				continue
			}

			if turma.DiaSemana == novaTurma.DiaSemana &&
				novaTurma.HoraInicio < turma.HoraFim &&
				novaTurma.HoraFim > turma.HoraInicio {

				c.JSON(http.StatusConflict, gin.H{
					"erro": "O aluno já está matriculado em uma turma com conflito de horário",
				})
				return
			}
		}
	}

	turmas[turmaEncontrada].Alunos = append(
		turmas[turmaEncontrada].Alunos,
		dados.AlunoID,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Aluno matriculado com sucesso",
		"turma":    montarTurmaResposta(turmas[turmaEncontrada]),
	})
}

func ListarAlunosDaTurma(c *gin.Context) {
	turmaID, err := strconv.Atoi(c.Param("id"))

	if err != nil || turmaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID da turma inválido",
		})
		return
	}

	var turma *models.Turma

	for i := range turmas {
		if turmas[i].ID == turmaID {
			turma = &turmas[i]
			break
		}
	}

	if turma == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "Turma não encontrada",
		})
		return
	}

	alunosDaTurma := make([]models.Aluno, 0, len(turma.Alunos))

	for _, alunoID := range turma.Alunos {
		for _, aluno := range alunos {
			if aluno.ID == alunoID {
				alunosDaTurma = append(alunosDaTurma, aluno)
				break
			}
		}
	}

	c.JSON(http.StatusOK, alunosDaTurma)
}

func AlocarSala(c *gin.Context) {
	turmaID, err := strconv.Atoi(c.Param("id"))

	if err != nil || turmaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID da turma inválido",
		})
		return
	}

	var dados struct {
		SalaID     int    `json:"sala_id"`
		DiaSemana  string `json:"dia_semana"`
		HoraInicio string `json:"hora_inicio"`
		HoraFim    string `json:"hora_fim"`
	}

	if err := c.ShouldBindJSON(&dados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	dados.DiaSemana = normalizarDia(dados.DiaSemana)

	if dados.SalaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O ID da sala é obrigatório",
		})
		return
	}

	if dados.DiaSemana == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dia da semana inválido",
		})
		return
	}

	if !validarHorario(dados.HoraInicio) ||
		!validarHorario(dados.HoraFim) {

		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Os horários devem estar no formato HH:MM",
		})
		return
	}

	if dados.HoraInicio >= dados.HoraFim {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "O horário de início deve ser anterior ao horário de fim",
		})
		return
	}

	turmaEncontrada := -1

	for i, turma := range turmas {
		if turma.ID == turmaID {
			turmaEncontrada = i
			break
		}
	}

	if turmaEncontrada == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "Turma não encontrada",
		})
		return
	}

	if !turmas[turmaEncontrada].Ativa {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "A turma está inativa",
		})
		return
	}

	salaEncontrada := -1

	for i, sala := range salas {
		if sala.ID == dados.SalaID {
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

	if len(turmas[turmaEncontrada].Alunos) >
		salas[salaEncontrada].Capacidade {

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro": "A capacidade da sala é insuficiente para a turma",
		})
		return
	}

	for _, turma := range turmas {
		if turma.ID == turmaID {
			continue
		}

		if !turma.Ativa {
			continue
		}

		if turma.SalaID == nil {
			continue
		}

		if *turma.SalaID != dados.SalaID {
			continue
		}

		if turma.DiaSemana != dados.DiaSemana {
			continue
		}

		if dados.HoraInicio < turma.HoraFim &&
			dados.HoraFim > turma.HoraInicio {

			c.JSON(http.StatusConflict, gin.H{
				"erro": "A sala já está ocupada nesse horário",
			})
			return
		}
	}

	salaID := salas[salaEncontrada].ID

	turmas[turmaEncontrada].SalaID = &salaID
	turmas[turmaEncontrada].DiaSemana = dados.DiaSemana
	turmas[turmaEncontrada].HoraInicio = dados.HoraInicio
	turmas[turmaEncontrada].HoraFim = dados.HoraFim

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Sala alocada com sucesso",
		"turma":    montarTurmaResposta(turmas[turmaEncontrada]),
	})
}
