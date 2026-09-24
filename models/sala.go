package models

type Sala struct {
	ID         int      `json:"id"`
	Nome       string   `json:"nome"`
	Capacidade int      `json:"capacidade"`
	Recursos   []string `json:"recursos"`
	Ativa      bool     `json:"ativa"`
}