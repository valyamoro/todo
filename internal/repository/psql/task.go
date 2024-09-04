package psql

import (
	"database/sql"
	"github.com/valyamoro/TODO/internal/domain"
)

type Tasks struct {
	db *sql.DB
}

func NewTasks(db *sql.DB) *Tasks {
	return &Tasks{db}
}

func (t *Tasks) Create(task domain.Task) (domain.Task, error) {
	_, err := t.db.Exec(
		"INSERT INTO tasks (name, description, time) values ($1, $2, $3)",
		task.Name,
		task.Description,
		task.Time,
	)

	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (t *Tasks) GetByID(id int64) (domain.Task, error) {
	var task domain.Task
	err := t.db.QueryRow("SELECT id, name, description, time FROM tasks WHERE id=$1", id).
		Scan(
			&task.ID, 
			&task.Name, 
			&task.Description, 
			&task.Time,
		)

	if err == sql.ErrNoRows {
		return task, domain.ErrTaskNotFound
	}

	return task, nil
}

func (t *Tasks) GetAll() ([]domain.Task, error) {
	rows, err := t.db.Query("SELECT id, name, description, time FROM tasks")
	if err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Time,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (t *Tasks) Delete(id int64) (domain.Task, error) {
	task, err := t.GetByID(id)
	if err != nil {
		return domain.Task{}, err
	}

	_, err = t.db.Exec("DELETE FROM tasks WHERE id=$1", id)

	if err != nil {
		return domain.Task{}, err
	}

	return task, err 
}

func (t *Tasks) Update(id int64, inp domain.UpdateTaskInput) (domain.Task, error) {
	_, err := t.db.Exec(
		"UPDATE tasks SET name=$1, description=$2, time=$3 WHERE id=$4",
		inp.Name,
		inp.Description,
		inp.Time,
		id, 
	)

	if err != nil {
		return domain.Task{}, err 
	}

	return t.GetByID(id)
}
