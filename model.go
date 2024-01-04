package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	help     help.Model
	loaded   bool
	focused  status
	cols     []column
	quitting bool
}

func NewBoard() *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, focused: todo}
}

func (b *Board) Init() tea.Cmd {
	return nil
}

func (b *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		b.help.Width = msg.Width - margin
		for i := 0; i < len(b.cols); i++ {
			var res tea.Model
			res, cmd = b.cols[i].Update(msg)
			b.cols[i] = res.(column)
			cmds = append(cmds, cmd)
		}
		b.loaded = true
		return b, tea.Batch(cmds...)
	case Form:
		return b, b.cols[b.focused].Set(msg.index, msg.CreateTask())
	case moveMsg:
		return b, b.cols[b.focused.Next()].Set(APPEND, msg.Task)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			b.quitting = true
			return b, tea.Quit
		case key.Matches(msg, keys.Left):
			b.cols[b.focused].Blur()
			b.focused = b.focused.Prev()
			b.cols[b.focused].Focus()
		case key.Matches(msg, keys.Right):
			b.cols[b.focused].Blur()
			b.focused = b.focused.Next()
			b.cols[b.focused].Focus()
		}
	}
	res, cmd := b.cols[b.focused].Update(msg)
	if _, ok := res.(column); ok {
		b.cols[b.focused] = res.(column)
	} else {
		return res, cmd
	}
	return b, cmd
}

func (b *Board) View() string {
	if b.quitting {
		return ""
	}
	if !b.loaded {
		return "loading..."
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		b.cols[todo].View(),
		b.cols[inProgress].View(),
		b.cols[done].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, b.help.View(keys))
}
