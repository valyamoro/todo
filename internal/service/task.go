package service

import (
	"github.com/valyamoro/TODO/internal/domain"
)

type TasksRepository interface {
	Create(task domain.Task) (domain.Task, error)
	GetByID(id int64) (domain.Task, error)
	GetAll() ([]domain.Task, error)
	Delete(id int64) (domain.Task, error)
	Update(id int64, inp domain.UpdateTaskInput) (domain.Task, error)
}

type Tasks struct {
	repo TasksRepository
}

func NewTasks(repo TasksRepository) *Tasks {
	return &Tasks{
		repo: repo,
	}
}

func (t *Tasks) Create(task domain.Task) (domain.Task, error) {
	return t.repo.Create(task)
}

func (b *Tasks) GetByID(id int64) (domain.Task, error) {
	return b.repo.GetByID(id)
}

func (b *Tasks) GetAll() ([]domain.Task, error) {
	return b.repo.GetAll()
}

func (b *Tasks) Delete(id int64) (domain.Task, error) {
	return b.repo.Delete(id)
}

func (b *Tasks) Update(id int64, inp domain.UpdateTaskInput) (domain.Task, error) {
	return b.repo.Update(id, inp)
}
