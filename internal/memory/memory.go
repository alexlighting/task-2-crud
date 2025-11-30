package memory

import (
	"fmt"
	"sync"
	"tasks-api/internal/models"
	"time"
)

type Memory struct {
	sync.RWMutex
	unicID int64
	Tasks  []models.Task
}

func New() *Memory {
	tasks := make([]models.Task, 0)
	return &Memory{Tasks: tasks}
}

func (m *Memory) Create(task models.Task) (models.Task, error) {
	//генерируем уникальный ID
	m.Lock()
	defer m.Unlock()
	m.unicID++
	task.ID = m.unicID
	task.CreatedAt = (time.Now().Format(time.RFC3339))
	m.Tasks = append(m.Tasks, task)
	return task, nil
}

func (m *Memory) Get(id int64) (models.Task, bool) {
	m.RLock()
	defer m.RUnlock()
	for _, task := range m.Tasks {
		if task.ID == id { //если нашли - возвращаем задачу
			return task, true
		}
	}
	return models.Task{}, false //задача с таким id не найдена
}

func (m *Memory) List() []models.Task {
	var tasksCopy []models.Task
	m.RLock()
	defer m.RUnlock()
	return append(tasksCopy, m.Tasks...)
}

func (m *Memory) Update(id int64, task models.Task) (models.Task, error) {
	m.Lock()
	defer m.Unlock()
	for index, tmp_task := range m.Tasks {
		if tmp_task.ID == id { //если нашли - возвращаем задачу
			if task.Title != "" {
				m.Tasks[index].Title = task.Title
			}
			m.Tasks[index].Done = task.Done
			return m.Tasks[index], nil
		}
	}
	return models.Task{}, fmt.Errorf("Wrong id")
}

func (m *Memory) Delete(id int64) error {
	m.Lock()
	defer m.Unlock()
	for index, task := range m.Tasks {
		if task.ID == id { //если нашли
			m.Tasks[index] = m.Tasks[len(m.Tasks)-1] //удаляем задачу из слайса
			m.Tasks = m.Tasks[:len(m.Tasks)-1]
			return nil
		}
	}
	return fmt.Errorf("Wrong id")
}
