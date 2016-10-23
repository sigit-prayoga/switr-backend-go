package model

type Todo struct {
	Label string
}

func (todo Todo) getLabel() string {
	return todo.Label
}
