package main

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

type Command string

const (
	Add            Command = "add"
	List           Command = "list"
	Get            Command = "get"
	Len            Command = "len"
	MarkInProgress Command = "mark-in-progress"
	MarkDone       Command = "mark-done"
	Update         Command = "update"
	Delete         Command = "delete"
)
