package todo

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"path"
)

// Реализация структур
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// NewTask - создание новой сущности, новой задачи
func NewTask(description string) Task {
	task := Task{
		ID:          rand.IntN(9999999),
		Description: description,
		Done:        false,
	}
	return task
}

// Load - инициализация слайса задач
func Load(tasks []Task) ([]Task, error) {
	buf, err := os.ReadFile(TasksJsonStorage)
	if err != nil {
		if os.IsNotExist(err) {
			os.Create(path.Base(TasksJsonStorage))
			return []Task{}, nil
		}
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}
	err = json.Unmarshal(buf, &tasks)
	if err != nil {
		return nil, fmt.Errorf("ошибка дисериализации: %w", err)
	}
	return tasks, nil
}

// Save перезаписывает изменения задач
func Save(data []Task) error {
	buf, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}

	return os.WriteFile(TasksJsonStorage, buf, 0644)
}
