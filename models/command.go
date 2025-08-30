package models

type Command string

const (
	Help           Command = "help"
	Add            Command = "add"
	List           Command = "list"
	Get            Command = "get"
	Len            Command = "len"
	MarkInProgress Command = "mark-in-progress"
	MarkDone       Command = "mark-done"
	Update         Command = "update"
	Delete         Command = "delete"
)
