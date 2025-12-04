package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"todo-app/internal/todo"
)

func LoadJSON(path string) ([]todo.Task, error) {
	tasks := []todo.Task{}
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(buf, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func SaveJSON(path string, tasks []todo.Task) error {
	buf, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}
	err = os.WriteFile(path, buf, 0666)
	if err != nil {
		return err
	}
	return nil
}
