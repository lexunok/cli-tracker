package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	id          int
	description string
	status      Status
	//created at and updated at
}

type Status string

const (
	TODO       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "done"
)

func main() {

	//Получаем аргументы
	command := os.Args[1]
	var description string
	if len(os.Args) > 2 {
		description = os.Args[2]
	}
	nameOfFile := "db.json"

	//Создаем файл если его не существует
	file, err := os.OpenFile(nameOfFile, os.O_RDONLY|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	switch command {

	//Если команда добавить
	case "add":

		//Считываем файл
		bytes, err := os.ReadFile(nameOfFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		//Декодируем список задач из байтов
		tasks := getTaskList(&bytes)

		//Создаем задачу
		task := createTask(len(tasks), description)

		//Определяем отступ для записи
		offset := int64(len(bytes))

		//Если первый объект то 0 иначе -2 байта чтобы перезаписать ]
		if len(tasks) != 0 {
			offset -= 2
		}

		//Записываем
		file.WriteAt([]byte(task), offset)

		//Если команда получить список
	case "list":

		//Считываем файл
		bytes, err := os.ReadFile(nameOfFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		//Декодируем список задач из байтов
		tasks := getTaskList(&bytes)

		fmt.Println("Список задач:")
		for _, el := range tasks {
			fmt.Printf("\t ( id:%d, '%s', статус: '%s' )\n", el.id, el.description, el.status)
		}

	}
}

func createTask(id int, description string) string {
	task := Task{
		id:          id,
		description: description,
		status:      TODO,
	}
	fmt.Println("Создана задача:", task.id, "-", task.description)

	encodedTask := fmt.Sprintf("\n\t{\"id\":%d, \"description\":\"%s\", \"status\":\"%s\"}\n]", task.id, task.description, task.status)

	if task.id == 0 {
		return "[" + encodedTask
	} else {
		return "," + encodedTask
	}
}

// Можно модифицировать, например с помощью trim или проходиться по списку
func getTaskList(bytes *[]byte) []Task {

	isValue := false
	isObject := false
	listOfTasks := make([]Task, 0)
	tempValue := make([]byte, 0)
	tempTask := Task{}
	count := 0

	for i := 0; i < len(*bytes); i++ {
		symbol := (*bytes)[i]

		switch symbol {
		case '{':
			isObject = true
		case '}':
			isObject = false
			isValue = false
			count = 0

			status := strings.Trim(string(tempValue), "\"")
			tempTask.status = Status(status)

			listOfTasks = append(listOfTasks, tempTask)
			tempTask = Task{}
			tempValue = make([]byte, 0)
		case ':':
			isValue = true
			continue
		case ',':
			if isObject {
				switch count {
				case 0:
					id, err := strconv.Atoi(string(tempValue))
					if err != nil {
						log.Fatal(err.Error())
					}
					tempTask.id = id
				case 1:
					tempTask.description = strings.Trim(string(tempValue), "\"")
				}

				isValue = false
				tempValue = make([]byte, 0)
				count++
			}
		}
		if isValue {
			tempValue = append(tempValue, symbol)
		}
	}
	return listOfTasks
}
