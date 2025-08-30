package main

import (
	"cli-tracker/models"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Service interface {
	GetOne(*[]byte, int) (models.Task, error)
	GetList(*[]byte, string) error
	GetLen(*[]byte) int
	Create(*[]byte, string) []byte
	Delete(*[]byte, int) []byte
	Update(*[]byte, int, string) []byte
	Mark(*[]byte, int, models.Status) []byte
}

func main() {

	start := time.Now()

	//Получаем аргументы
	var command models.Command
	if len(os.Args) > 1 {
		command = models.Command(os.Args[1])
	}

	var secondArg string
	if len(os.Args) > 2 {
		secondArg = os.Args[2]
	}

	var thirdArg string
	if len(os.Args) > 3 {
		thirdArg = os.Args[3]
	}

	service := Service(FastService{})
	nameOfFile := "tasks.json"

	//Создаем файл если его не существует
	file, err := os.OpenFile(nameOfFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	//Считываем размер файла
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err.Error())
	}
	size := stat.Size()

	//Считываем байты в буффер
	bytes := make([]byte, size)
	_, err = file.Read(bytes)
	if err != nil && err != io.EOF {
		log.Fatal(err.Error())
	}

	switch command {
	//Если команда Help
	case models.Help:
		fmt.Println(`
			# Просмотреть список возможных команд
			go run . help

			# Добавить задачу
			go run . add "Buy groceries"

			# Обновить задачу по id
			go run . update 1 "New description"

			# Удалить задачу по id
			go run . delete 1

			# Отметить задачу как "в процессе"
			go run . mark-in-progress 1

			# Отметить задачу как "выполнено"
			go run . mark-done 1

			# Получить кол-во задач
			go run . len

			# Показать все задачи
			go run . list

			# Показать только выполненные
			go run . list done

			# Показать только "в планах"
			go run . list todo

			# Показать только "в процессе"
			go run . list in-progress
		`)
	//Если команда добавить
	case models.Add:

		//Добавляем задачу
		data := service.Create(&bytes, secondArg)

		//Записываем
		if _, err := file.WriteAt(data, 0); err != nil {
			log.Fatal(err.Error())
		}

	//Если команда получить список
	case models.List:

		//Выводим список задач из байтов по фильтру если он есть
		if err := service.GetList(&bytes, secondArg); err != nil {
			log.Fatal(err.Error())
		}

	//Если команда получить кол-во задач
	case models.Len:

		//Получить кол-во задач
		number := service.GetLen(&bytes)

		//Вывести кол-во задач
		fmt.Println("Всего задач -", number)

	//Если команда получить задачу
	case models.Get:

		//Получаем количество задач в файле
		length := service.GetLen(&bytes) - 1

		//ВНИМАНИЕ -  А что если удалить задачу с id 2, то ее же можно будет получить, но ее самой не существует. Нужно наверное возвращать что такой задачи нет
		id, err := strconv.Atoi(secondArg)
		if err != nil || id < 0 || id > length {
			log.Fatal("Id должен быть в диапазоне от 0 до ", length)
		}

		task, err := service.GetOne(&bytes, id)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("( id:%d, '%s', статус: '%s' )\n", task.Id, task.Description, task.Status)

	//Если команда обновить задачу
	case models.Update:

		//Получаем количество задач в файле
		length := service.GetLen(&bytes) - 1

		//ВНИМАНИЕ -  А что если удалить задачу с id 2, то ее же можно будет получить, но ее самой не существует. Нужно наверное возвращать что такой задачи нет
		id, err := strconv.Atoi(secondArg)
		if err != nil || id < 0 || id > length {
			log.Fatal("Id должен быть в диапазоне от 0 до ", length)
		}

		data := service.Update(&bytes, id, thirdArg)

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")

	//Если команда удалить
	case models.Delete:

		//Получаем количество задач в файле
		length := service.GetLen(&bytes) - 1

		//ВНИМАНИЕ -  А что если удалить задачу с id 2, то ее же можно будет получить, но ее самой не существует. Нужно наверное возвращать что такой задачи нет
		id, err := strconv.Atoi(secondArg)
		if err != nil || id < 0 || id > length {
			log.Fatal("Id должен быть в диапазоне от 0 до ", length)
		}

		data := service.Delete(&bytes, id)

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")

	//Если команда отметить поменять статус задачи на "В прогрессе"
	case models.MarkInProgress:

		//Получаем количество задач в файле
		length := service.GetLen(&bytes) - 1

		//ВНИМАНИЕ -  А что если удалить задачу с id 2, то ее же можно будет получить, но ее самой не существует. Нужно наверное возвращать что такой задачи нет
		id, err := strconv.Atoi(secondArg)
		if err != nil || id < 0 || id > length {
			log.Fatal("Id должен быть в диапазоне от 0 до ", length)
		}

		data := service.Mark(&bytes, id, models.InProgress)

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")

	//Если команда отметить поменять статус задачи на "Выполнено"
	case models.MarkDone:

		//Получаем количество задач в файле
		length := service.GetLen(&bytes) - 1

		//ВНИМАНИЕ -  А что если удалить задачу с id 2, то ее же можно будет получить, но ее самой не существует. Нужно наверное возвращать что такой задачи нет
		id, err := strconv.Atoi(secondArg)
		if err != nil || id < 0 || id > length {
			log.Fatal("Id должен быть в диапазоне от 0 до ", length)
		}

		data := service.Mark(&bytes, id, models.Done)

		//ВНИМАНИЕ -Обработать ошибку
		file.Write(data)

		if err := file.Truncate(int64(len(data))); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Успешно")
	}

	fmt.Println("Время выполнения команды:", time.Since(start).Microseconds())
}
