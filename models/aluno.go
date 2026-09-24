package models

type Aluno struct {
	ID        int    `json:"id"`
	Matricula string `json:"matricula"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
}