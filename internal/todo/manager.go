// Реализация структур
package todo

import (
	"fmt"
)

var TasksJsonStorage = "..tasks.json" // не делал константой для гибкости

// Add - добавление новой задачи
func Add(tasks []Task, desc string) []Task {
	// проверка на дублирование задач
	for _, task := range tasks {
		if task.Description == desc && !task.Done { // если задача с таким описанием уже в списке И не выполнена
			fmt.Printf("у вас уже есть невыполненная задача \"%s\"\n", desc)
			return tasks
		}
	}
	newTask := NewTask(desc)       // создаем структуру новой задачи
	tasks = append(tasks, newTask) // добавляем ее в слайс задач, который возвращаем

	//необязательный вывод, больше для отладки
	fmt.Printf("Задача \"%s\" успешно добавлена!, ID задачи - %d\n", newTask.Description, newTask.ID)
	return tasks
}

// List фильтрует список задач.
func List(tasks []Task, filter string) []Task {
	switch filter {
	case "done":
		var doneTasks []Task
		for _, task := range tasks {
			if task.Done {
				doneTasks = append(doneTasks, task)
			}

		}
		return doneTasks
	case "pending":
		var pendingTasks []Task
		for _, task := range tasks {
			if !task.Done {
				pendingTasks = append(pendingTasks, task)
			}
		}
		return pendingTasks
	}
	return tasks // default if "all"
}

// Complete меняет статус задачи на "выполнено"
func Complete(tasks []Task, id int) ([]Task, error) {
	flag := false
	for idx, task := range tasks {
		if task.ID == id && !task.Done {
			tasks[idx].Done = true
			flag = true
			break
		} else if task.ID == id && task.Done {
			return nil, fmt.Errorf("задача с ID %d уже выполнена", task.ID)
		}
	}
	if flag {
		fmt.Println("Статус задачи изменен на \"выполнена\", ее id: ", id)
	} else {
		return nil, fmt.Errorf("задача с таким ID не найдена")
	}
	// не создавал копию, так как нам требуется в основном слайсе поменять статус у задачи
	return tasks, nil
}

// Delete удаляет задачу из списка задач
func Delete(tasks []Task, id int) ([]Task, error) {
	flag := false
	for idx, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:idx], tasks[idx+1:]...)
			flag = true
			break
		}
	}
	if flag {
		fmt.Println("Задача удалена из списка, ее id: ", id)
	} else {
		return nil, fmt.Errorf("задача с таким ID не найдена")
	}
	// не создавал копию, так как нам требуется в основном слайсе поменять статус у задачи
	return tasks, nil
}
