package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

func (s status) Next() status {
	if s == done {
		return todo
	}
	return s + 1
}

func (s status) Prev() status {
	if s == todo {
		return done
	}
	return s - 1
}

const margin = 4

var board *Board

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	board = NewBoard(db)
	p := tea.NewProgram(board)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
