package database

import (
	"database/sql"
	"fmt"
	"tasks/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `SELECT * FROM tasks ORDER BY created_at DESC`
	err := s.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskStore) GetById(id int) (*models.Task, error) {
	var task models.Task

	query := `SELECT * FROM tasks WHERE id = $1`
	err := s.db.Select(&task, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task with id %d not found", id)
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task
	
	query := `INSERT INTO tasks (title, description, completed, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5) returning *
	`
	now := time.Now()

	err := s.db.QueryRowx(query, input.Title, input.Description, input.Completed, now, now).StructScan(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
