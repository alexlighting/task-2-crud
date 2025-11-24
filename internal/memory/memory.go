package memory

import (
	"fmt"
	"sync"
	"tasks-api/internal/models"
	"time"
)

type Memory struct {
	sync.Mutex
	Tasks []models.Task
}

func New() *Memory {
	tasks := make([]models.Task, 0)
	return &Memory{Tasks: tasks}
}

func (m *Memory) Create(task models.Task) (models.Task, error) {
	//генерируем уникальный ID
	//проблема с индексами, при удалении элементы слайса перемещаются
	//и последний в слайсе элемент может не быть с наибльшим id
	if len(m.Tasks) > 0 {
		task.ID = m.Tasks[len(m.Tasks)-1].ID + 1
	} else {
		task.ID = 1
	}
	task.CreatedAt = (time.Now().Format(time.RFC3339))

	m.Lock()
	m.Tasks = append(m.Tasks, task)
	m.Unlock()
	return task, nil
}

func (m *Memory) Get(id int) (models.Task, bool) {
	for _, task := range m.Tasks {
		if task.ID == id { //если нашли - возвращаем задачу
			return task, true
		}
	}
	return models.Task{}, false //задача с таким id не найдена
}

func (m *Memory) List() []models.Task {
	return m.Tasks
}

func (m *Memory) Update(id int, tasks models.Task) (models.Task, error) {
	for index, task := range m.Tasks {
		if task.ID == id { //если нашли - возвращаем задачу
			if tasks.Title != "" {
				m.Lock()
				m.Tasks[index].Title = tasks.Title
				m.Unlock()
			}
			m.Lock()
			m.Tasks[index].Done = tasks.Done
			m.Unlock()
			return task, nil
		}
	}
	return models.Task{}, fmt.Errorf("Задача c id - %d не найдена", id)
}

func (m *Memory) Delete(id int) error {
	for index, task := range m.Tasks {
		if task.ID == id { //если нашли
			m.Lock()
			m.Tasks[index] = m.Tasks[len(m.Tasks)-1] //удаляем задачу из слайса
			m.Tasks = m.Tasks[:len(m.Tasks)-1]
			m.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Задача c id - %d не найдена", id)
}
