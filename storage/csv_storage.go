package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"todo-app/internal/todo"
)

func LoadCSV(path string) ([]todo.Task, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	} else if len(records) < 2 {
		return nil, fmt.Errorf("пустой исходник")
	}
	tasks := []todo.Task{}
	for idx, record := range records {
		if idx == 0 {
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		done, err := strconv.ParseBool(record[2])
		if err != nil {
			return nil, err
		}
		task := todo.Task{
			ID:          id,
			Description: record[1],
			Done:        done,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func SaveCSV(path string, tasks []todo.Task) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	for idx, task := range tasks {
		forWrite := []string{
			strconv.FormatInt(int64(task.ID), 10),
			task.Description,
			strconv.FormatBool(task.Done),
		}
		if idx == 0 {
			head := []string{"id", "description", "done"}
			if err := writer.Write(head); err != nil {
				return err
			}
		}
		if err := writer.Write(forWrite); err != nil {
			return err
		}
		writer.Flush()
	}
	return nil
}
