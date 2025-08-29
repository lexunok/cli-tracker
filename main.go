package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	//Получаем аргументы
	var command Command
	if len(os.Args) > 1 {
		command = Command(os.Args[1])
	}

	var secondArg string
	if len(os.Args) > 2 {
		secondArg = os.Args[2]
	}

	var thirdArg string
	if len(os.Args) > 3 {
		thirdArg = os.Args[3]
	}

	nameOfFile := "tasks.json"

	//Создаем файл если его не существует
	file, err := os.OpenFile(nameOfFile, os.O_RDONLY|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	//Если команда добавить
	switch command {
	case Add:

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

	//Если команда получить список
	case List:

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

	//Если команда получить задачу
	case Get:
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

	//Если команда обновить задачу
	case Update:
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

		data := updateTask(&bytes, id, thirdArg)

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")

	//Если команда удалить
	case Delete:
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

	//Если команда отметить поменять статус задачи на "В прогрессе"
	case MarkInProgress:
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

		data := markTask(&bytes, id, InProgress)

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")

	//Если команда отметить поменять статус задачи на "Выполнено"
	case MarkDone:
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

		data := markTask(&bytes, id, Done)

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")
	}
}
