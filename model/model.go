package model

type Choice struct {
	ID   string `json:"id"`
	Gone bool   `json:"gone"`
	Come bool   `json:"come"`
}