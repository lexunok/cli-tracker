package models

type Status string

const (
	TODO       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "done"
)
