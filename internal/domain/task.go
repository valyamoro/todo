package domain

import (
	"errors"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Time string `json:"time"`
}

type UpdateTaskInput struct {
	ID *int64 `json:"id"`
	Name *string `json:"name"`
	Description *string `json:"description"`
	Time *string `json:"time"`
}
