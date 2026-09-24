package models

type Turma struct {
	ID         int    `json:"id"`
	Nome       string `json:"nome"`
	Disciplina string `json:"disciplina"`
	Professor  string `json:"professor"`
	Alunos     []int  `json:"alunos"`
	SalaID     *int   `json:"sala_id"`
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
	Ativa      bool   `json:"ativa"`
}