package models

type ErrorReturn struct {
	Message string `json:"message"`
	Success bool `json:"success"`
}
