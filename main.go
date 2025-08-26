package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Task struct {
	Id          int
	Description string
	Status      Status
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
	var secondArg string
	if len(os.Args) > 2 {
		secondArg = os.Args[2]
	}
	nameOfFile := "db.json"

	//Создаем файл если его не существует
	file, err := os.OpenFile(nameOfFile, os.O_RDONLY|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	//Если команда добавить
	if command == "add" {

		//Считываем файл
		//ВНИМАНИЕ - Мы второй раз открываем файл методом ReadFile. Нужно сделать свою реализацию
		bytes, err := os.ReadFile(nameOfFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		//Получаем количество задач в файле
		length := getLenOfTaskList(&bytes)

		//Создаем задачу
		task := createTask(length, secondArg)

		//Определяем отступ для записи
		offset := int64(len(bytes))

		//ВНИМАНИЕ - Мне не нравится это решение
		//Если первый объект то 0 иначе -2 байта чтобы перезаписать ']'
		if length != 0 {
			offset -= 2
		}

		//Записываем
		file.WriteAt([]byte(task), offset)

	} else if command == "list" { //Если команда получить список

		//Считываем файл
		bytes, err := os.ReadFile(nameOfFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		//Декодируем список задач из байтов
		tasks, _ := getTaskList(&bytes, secondArg)

		fmt.Println("Список задач:", len(tasks))
		for _, el := range tasks {
			fmt.Printf("\t ( id:%d, '%s', статус: '%s' )\n", el.Id, el.Description, el.Status)
		}

	} else if command == "get" { //Если команда получить задачу
		//Считываем файл
		bytes, err := os.ReadFile(nameOfFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		//Получаем количество задач в файле
		length := getLenOfTaskList(&bytes) - 1

		//ВНИМАНИЕ -  А что если удалить задачу с id 2, то ее же можно будет получить, но ее самой не существует. Нужно наверное возвращать что такой задачи нет
		id, err := strconv.Atoi(secondArg)
		if err != nil || id < 0 || id > length {
			log.Fatal("Id должен быть в диапазоне от 0 до ", length)
		}

		task, err := getTaskById(&bytes, id)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("( id:%d, '%s', статус: '%s' )\n", task.Id, task.Description, task.Status)

	} else if command == "delete" {
		//Считываем файл
		bytes, err := os.ReadFile(nameOfFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		//Получаем количество задач в файле
		length := getLenOfTaskList(&bytes) - 1

		//ВНИМАНИЕ -  А что если удалить задачу с id 2, то ее же можно будет получить, но ее самой не существует. Нужно наверное возвращать что такой задачи нет
		id, err := strconv.Atoi(secondArg)
		if err != nil || id < 0 || id > length {
			log.Fatal("Id должен быть в диапазоне от 0 до ", length)
		}

		data := deleteTask(&bytes, id)

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")

	} else if strings.HasPrefix(command, "mark") { //Если команда отметить новый статус задачи
		//Считываем файл
		bytes, err := os.ReadFile(nameOfFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		//Получаем количество задач в файле
		length := getLenOfTaskList(&bytes) - 1

		//ВНИМАНИЕ -  А что если удалить задачу с id 2, то ее же можно будет получить, но ее самой не существует. Нужно наверное возвращать что такой задачи нет
		id, err := strconv.Atoi(secondArg)
		if err != nil || id < 0 || id > length {
			log.Fatal("Id должен быть в диапазоне от 0 до ", length)
		}

		//ВНИМАНИЕ - Может сделать валидацию на смену статуса, а то так можно на любой поменять
		suffix, _ := strings.CutPrefix(command, "mark-")
		data := markTask(&bytes, id, Status(suffix))

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")
	}
}

func createTask(id int, description string) string {

	fmt.Println("Создана задача:", id, "-", description)

	encodedTask := fmt.Sprintf("\n\t{\"id\":%d, \"description\":\"%s\", \"status\":\"%s\"}\n]", id, description, TODO)

	if id == 0 {
		return "[" + encodedTask
	} else {
		return "," + encodedTask
	}
}
func deleteTask(bytes *[]byte, id int) []byte {

	isValue := false
	isFound := false
	tempValue := make([]byte, 0)
	count := -1
	startOffset := -1
	endOffset := -1

	for i := 0; i < len(*bytes); i++ {
		symbol := (*bytes)[i]

		if symbol == ':' {
			count++
			isValue = true
			continue
		} else if (symbol == '}' || symbol == ',') && count >= 0 {

			if !isFound && count == 2 {
				startOffset = i + 1
			}

			isValue = false

			if count == 2 { //ВНИМАНИЕ - Сильная привязка что последнее поле по счету 2
				if isFound {
					endOffset = i + 1
					break
				}
				count = -1
			} else if count == 0 { //ВНИМАНИЕ - Сильная привязка к id что поле должно быть на первом месте
				isFound = string(tempValue) == strconv.Itoa(id)
			}

			tempValue = make([]byte, 0)
		}
		if isValue {
			tempValue = append(tempValue, symbol)
		}
	}

	return append((*bytes)[:startOffset], (*bytes)[endOffset:]...)
}
func markTask(bytes *[]byte, id int, status Status) []byte {

	isValue := false
	isFound := false
	tempValue := make([]byte, 0)
	count := -1
	startOffset := -1
	endOffset := -1

	for i := 0; i < len(*bytes); i++ {
		symbol := (*bytes)[i]

		if symbol == ':' {
			count++
			if count == 2 && isFound {
				startOffset = i + 1
			} else {
				isValue = true
				continue
			}
		} else if (symbol == '}' || symbol == ',') && count >= 0 {
			isValue = false

			if (count + 1) == 3 { //ВНИМАНИЕ - Сильная привязка что последнее поле по счету 3
				if isFound {
					endOffset = i
					break
				}
				count = -1
			} else if count == 0 { //ВНИМАНИЕ - Сильная привязка к id что поле должно быть на первом месте
				isFound = string(tempValue) == strconv.Itoa(id)
			}

			tempValue = make([]byte, 0)
		}
		if isValue {
			tempValue = append(tempValue, symbol)
		}
	}

	return append((*bytes)[:startOffset], append([]byte(`"`+status+`"`), (*bytes)[endOffset:]...)...)
}
func getTaskById(bytes *[]byte, id int) (Task, error) {

	isValue := false
	tempValue := make([]byte, 0)
	task := Task{}
	count := -1

	v := reflect.ValueOf(&task).Elem()

	for i := 0; i < len(*bytes); i++ {
		symbol := (*bytes)[i]

		if symbol == ':' {
			count++
			isValue = true
			continue
		} else if (symbol == '}' || symbol == ',') && count >= 0 {
			isValue = false

			field := v.Field(count)

			if field.Kind() == reflect.String {

				field.SetString(strings.Trim(string(tempValue), "\""))

			} else if field.Kind() == reflect.Int {
				number, err := strconv.Atoi(string(tempValue))
				if err != nil {
					return task, err
				}
				field.SetInt(int64(number))
			}

			if v.NumField() == (count + 1) {
				if task.Id == id {
					break
				}
				count = -1
				task = Task{}
			}
			tempValue = make([]byte, 0)
		}
		if isValue {
			tempValue = append(tempValue, symbol)
		}
	}

	return task, nil
}

func getTaskList(bytes *[]byte, filter string) ([]Task, error) {

	isValue := false
	listOfTasks := make([]Task, 0)
	tempValue := make([]byte, 0)
	tempTask := Task{}
	count := -1

	v := reflect.ValueOf(&tempTask).Elem()

	for i := 0; i < len(*bytes); i++ {
		symbol := (*bytes)[i]

		if symbol == ':' {
			count++
			isValue = true
			continue
		} else if (symbol == '}' || symbol == ',') && count >= 0 {
			isValue = false

			field := v.Field(count)

			if field.Kind() == reflect.String {

				field.SetString(strings.Trim(string(tempValue), "\""))

			} else if field.Kind() == reflect.Int {
				number, err := strconv.Atoi(string(tempValue))
				if err != nil {
					return listOfTasks, err
				}
				field.SetInt(int64(number))
			}

			if v.NumField() == (count + 1) {
				if tempTask.Status == Status(filter) || filter == "" {
					listOfTasks = append(listOfTasks, tempTask)
				}
				count = -1
				tempTask = Task{}
			}
			tempValue = make([]byte, 0)
		}
		if isValue {
			tempValue = append(tempValue, symbol)
		}
	}

	return listOfTasks, nil
}
func getLenOfTaskList(bytes *[]byte) int {

	task := Task{}
	count := 0

	v := reflect.ValueOf(&task).Elem()

	for i := 0; i < len(*bytes); i++ {
		if (*bytes)[i] == ':' {
			count++
		}
	}

	return count / v.NumField()
}
