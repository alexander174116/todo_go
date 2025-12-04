package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-app/internal/todo"
	"todo-app/storage"
)

// default storage path in TasksJsonStorage: tasks.json
// for change - use: todo.TasksJsonStorage = "your_path"
func main() {
	var tasks []todo.Task
	tasks, err := todo.Load(tasks) // читаю текущий список задач из нашей базы, получаю слайс структур
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var usedFlagArg string
	if len(os.Args) < 2 {
		fmt.Println("введите команду для запуска программы\nдоступные команды: add, list, complete/delete, load, export")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "add":
		command := flag.NewFlagSet("add", flag.ContinueOnError)
		command.StringVar(&usedFlagArg, "desc", "none", "--desc=\"new_task\"")
		err := command.Parse(os.Args[2:])
		if err != nil {
			if err == flag.ErrHelp {
				os.Exit(0)
			}
			fmt.Println("неверный флаг или неверное значение аргумента флага")
			os.Exit(1)
		}
		if usedFlagArg == "none" || usedFlagArg == "" {
			fmt.Println("не указано описание или флаг\nдоступный флаг для данной команды: --desc")
			os.Exit(1)
		}
		tasks = todo.Add(tasks, usedFlagArg)
	case "list":
		command := flag.NewFlagSet("list", flag.ContinueOnError)
		command.StringVar(&usedFlagArg, "filter", "all", "--filter=all/done/pending")
		err := command.Parse(os.Args[2:])
		if err != nil {
			if err == flag.ErrHelp {
				os.Exit(0)
			}
			fmt.Println("неверный флаг или неверное значение аргумента флага") // информативное сообщение об ошибке уже выведено в flag.Parse
			os.Exit(1)
		}
		switch usedFlagArg {
		case "all":
			if len(tasks) == 0 {
				fmt.Println("Список задач пуст либо отсутствует!")
			} else {
				fmt.Println("Список всех задач:")
			}
			for idx, task := range tasks {
				status := "выполняется"
				if task.Done {
					status = "выполнена"
				}
				fmt.Printf("Задача №%3d: ID - %19d | Статус задачи: %11s | Описание задачи: \"%s\".\n", idx+1, task.ID, status, task.Description)
			}
		case "done":
			filtered := todo.List(tasks, usedFlagArg)
			if len(filtered) == 0 {
				fmt.Println("Выполненные задачи не найдены!")
			} else {
				fmt.Println("Список выполненных задач:")
			}
			for idx, task := range filtered {
				fmt.Printf("Задача №%3d: ID задачи - %d; описание выполненной задачи - \"%s\".\n", idx+1, task.ID, task.Description)
			}
		case "pending":
			filtered := todo.List(tasks, usedFlagArg)

			if len(filtered) == 0 {
				fmt.Println("Невыполненные задачи не найдены!")
			} else {
				fmt.Println("Задачи, находящиеся в процессе выполнения:")
			}
			for idx, task := range filtered {

				fmt.Printf("Задача №%3d: ID задачи - %d; описание невыполненной задачи - \"%s\".\n", idx+1, task.ID, task.Description)
			}
		default:
			fmt.Println("доступные команды для --filter \"all/done/pending\"")
		}
	case "complete":
		command := flag.NewFlagSet("complete", flag.ContinueOnError)
		command.StringVar(&usedFlagArg, "id", "none", "--id=[task_id]")
		err := command.Parse(os.Args[2:])
		if err != nil {
			if err == flag.ErrHelp {
				os.Exit(0)
			}
			fmt.Println("неверный флаг или неверное значение аргумента флага")
			os.Exit(1)
		}
		if usedFlagArg == "none" {
			fmt.Println("необходимо указать id задачи - --id=[task_id]")
			os.Exit(1)
		}
		if idInInt, err := strconv.Atoi(usedFlagArg); err != nil {
			fmt.Println("введен неверный id / требуется id задачи (набор цифр)")
			os.Exit(1)
		} else {
			_, err := todo.Complete(tasks, idInInt)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	case "delete":
		command := flag.NewFlagSet("delete", flag.ContinueOnError)
		command.StringVar(&usedFlagArg, "id", "none", "--id=[task_id]")
		err := command.Parse(os.Args[2:])
		if err != nil {
			if err == flag.ErrHelp {
				os.Exit(0)
			}
			fmt.Println("неверный флаг или неверное значение аргумента флага")
			os.Exit(1)
		}
		if usedFlagArg == "none" {
			fmt.Println("необходимо указать id задачи - --id=[task_id]")
			os.Exit(1)
		}
		if idInInt, err := strconv.Atoi(usedFlagArg); err != nil {
			fmt.Println("требуется id задачи (набор цифр)")
			os.Exit(1)
		} else {
			tasks, err = todo.Delete(tasks, idInInt)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	case "export":
		var outPath string
		command := flag.NewFlagSet("export", flag.ContinueOnError)
		command.StringVar(&usedFlagArg, "format", "none", "--format=json/csv")
		command.StringVar(&outPath, "out", "none", "--out=[export_path]")
		if err := command.Parse(os.Args[2:]); err != nil {
			fmt.Println("usage for export: --format [json/scv] --out [export_path]")
			os.Exit(1)
		}

		switch usedFlagArg {
		case "json":
			if outPath == "" || outPath == "none" {
				fmt.Println("введите имя файла, в который экспортировать задачи: --out=[export_path]")
				os.Exit(1)
			}
			if err := storage.SaveJSON(outPath, tasks); err != nil {
				fmt.Println("ошибка экспорта:", err)
				os.Exit(1)
			}

			fmt.Println("export done in path:", outPath)
		case "csv":
			if outPath == "" || outPath == "none" {
				fmt.Println("введите имя файла, в который экспортировать задачи: --out=[export_path]")
				os.Exit(1)
			}
			if err := storage.SaveCSV(outPath, tasks); err != nil {
				fmt.Println("ошибка экспорта:", err)
				os.Exit(1)
			}

			fmt.Println("export done in path:", outPath)
		default:
			fmt.Println("выберите формат экспорта - --format=json/csv")
			os.Exit(1)
		}
	case "load":
		command := flag.NewFlagSet("load", flag.ContinueOnError)
		command.StringVar(&usedFlagArg, "file", "none", "--file=[todo_data_file]")
		command.Parse(os.Args[2:])
		if strings.HasSuffix(usedFlagArg, ".json") {
			tasksNew, err := storage.LoadJSON(usedFlagArg)
			if err != nil {
				fmt.Println("ошибка импорта задач:", err)
				os.Exit(1)
			}
			todo.Save(tasksNew)
		} else if strings.HasSuffix(usedFlagArg, ".csv") {
			tasksNew, err := storage.LoadCSV(usedFlagArg)
			if err != nil {
				fmt.Println("ошибка импорта задач:", err)
				os.Exit(1)
			}
			todo.Save(tasksNew)
		} else {
			fmt.Println("загрузка данных доступна только из json/csv\nиспользование load --file=[path]")
			os.Exit(1)
		}

		fmt.Println("Задачи успешно загружены!")
		os.Exit(0)
	default:
		fmt.Println("доступные команды: add, list, complete/delete, load, LoadJSON/LoadCSV, export")
		os.Exit(1)
	}

	todo.Save(tasks)
}
