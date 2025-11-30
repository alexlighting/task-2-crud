package storage

import "tasks-api/internal/models"

type Storage interface {
	List() []models.Task
	Create(models.Task) (models.Task, error)
	Get(id int64) (models.Task, bool)
	Update(id int64, tasks models.Task) (models.Task, error)
	Delete(id int64) error
}
