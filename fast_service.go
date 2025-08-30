package main

import (
	"cli-tracker/models"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FastService struct {
	Service
}

// READY
func (FastService) GetList(bytes *[]byte, filter string) error {

	fmt.Println("Список задач:")

	var (
		isValue   bool
		tempValue []byte
		tempTask  models.Task
	)
	count := -1

	v := reflect.ValueOf(&tempTask).Elem()

	for i := 0; i < len(*bytes); i++ {
		symbol := (*bytes)[i]

		if symbol == ':' && !isValue {
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
					return err
				}
				field.SetInt(int64(number))
			} else if field.Kind() == reflect.TypeOf(time.Time{}).Kind() {
				dateTime, err := time.Parse(time.RFC3339, strings.Trim(string(tempValue), "\""))
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(dateTime))
			}

			if v.NumField() == (count + 1) {
				if tempTask.Status == models.Status(filter) || filter == "" {
					createdAt := tempTask.CreatedAt.Format(time.DateTime)
					updatedAt := tempTask.UpdatedAt.Format(time.DateTime)
					fmt.Printf("\t ( id:%d, '%s', статус: '%s', createdAt: '%s', updatedAt: '%s' )\n", tempTask.Id, tempTask.Description, tempTask.Status, createdAt, updatedAt)
				}
				count = -1
				tempTask = models.Task{}
			}
			tempValue = make([]byte, 0)
		}
		if isValue {
			tempValue = append(tempValue, symbol)
		}
	}
	return nil
}

// READY
func (FastService) GetLen(bytes *[]byte) int {

	var (
		task    models.Task
		count   int
		isValue bool
	)

	for i := 0; i < len(*bytes); i++ {
		symbol := (*bytes)[i]
		if symbol == ':' && !isValue {
			count++
			isValue = true
		} else if (symbol == '}' || symbol == ',') && count >= 0 {
			isValue = false
		}
	}

	return count / reflect.ValueOf(&task).Elem().NumField()
}

// READY
func (s FastService) Create(bytes *[]byte, description string) []byte {

	offset := 0
	id := s.GetLen(bytes)

	fmt.Println("Создана задача:", id, "-", description)
	time := time.Now().Format(time.RFC3339)
	encodedTask := fmt.Sprintf("\n\t{\"id\":%d, \"description\":\"%s\", \"status\":\"%s\", \"createdAt\":\"%s\", \"updatedAt\":\"%s\"}\n]", id, description, models.TODO, time, time)

	if id == 0 {
		encodedTask = "[" + encodedTask
	} else {
		encodedTask = "," + encodedTask
		offset = len(*bytes) - 2
	}

	return append((*bytes)[:offset], []byte(encodedTask)...)
}
func (FastService) Delete(bytes *[]byte, id int) []byte {

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
func (FastService) Update(bytes *[]byte, id int, description string) []byte {

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
func (FastService) Mark(bytes *[]byte, id int, status models.Status) []byte {

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
func (FastService) GetOne(bytes *[]byte, id int) (models.Task, error) {

	isValue := false
	tempValue := make([]byte, 0)
	task := models.Task{}
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
				task = models.Task{}
			}
			tempValue = make([]byte, 0)
		}
		if isValue {
			tempValue = append(tempValue, symbol)
		}
	}

	return task, nil
}
