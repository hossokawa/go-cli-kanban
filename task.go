package main

import "github.com/charmbracelet/bubbles/list"

type Task struct {
	ID              uint   `db:"id"`
	TaskTitle       string `db:"title"`
	TaskDescription string `db:"description"`
	Status          string `db:"status"`
}

func NewTask(status, title, description string) Task {
	return Task{Status: status, TaskTitle: title, TaskDescription: description}
}

func (t *Task) Next() {
	switch t.Status {
	case todo.String():
		t.Status = inProgress.String()
	case inProgress.String():
		t.Status = done.String()
	case done.String():
		t.Status = todo.String()
	}
}

func (t Task) FilterValue() string {
	return t.TaskTitle
}

func (t Task) Title() string {
	return t.TaskTitle
}

func (t Task) Description() string {
	return t.TaskDescription
}

func tasksToItems(tasks []Task) []list.Item {
	var items []list.Item
	for _, t := range tasks {
		items = append(items, t)
	}
	return items
}
