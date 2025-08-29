package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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
func updateTask(bytes *[]byte, id int, description string) []byte {

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
			//ВНИМАНИЕ - Сильная привязка
			if count == 1 && isFound {
				startOffset = i + 1
			} else {
				isValue = true
				continue
			}
		} else if (symbol == '}' || symbol == ',') && count >= 0 {
			isValue = false

			if count == 2 { //ВНИМАНИЕ - Сильная привязка что последнее поле по счету 2
				count = -1
			} else if count == 1 && isFound {
				endOffset = i
				break
			} else if count == 0 { //ВНИМАНИЕ - Сильная привязка к id что поле должно быть на первом месте
				isFound = string(tempValue) == strconv.Itoa(id)
			}

			tempValue = make([]byte, 0)
		}
		if isValue {
			tempValue = append(tempValue, symbol)
		}
	}

	return append((*bytes)[:startOffset], append([]byte(`"`+description+`"`), (*bytes)[endOffset:]...)...)
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
			//ВНИМАНИЕ - Сильная привязка
			if count == 2 && isFound {
				startOffset = i + 1
			} else {
				isValue = true
				continue
			}
		} else if (symbol == '}' || symbol == ',') && count >= 0 {
			isValue = false

			if count == 2 { //ВНИМАНИЕ - Сильная привязка что последнее поле по счету 2
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
