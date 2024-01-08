package main

import (
	"errors"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1

type column struct {
	focus  bool
	status status
	list   list.Model
	height int
	width  int
	tasks  []Task
}

func (c *column) Focus() {
	c.focus = true
}

func (c *column) Blur() {
	c.focus = false
}

func (c *column) Focused() bool {
	return c.focus
}

func newColumn(items []list.Item, status status, focus bool) column {
	defaultList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return column{focus: focus, status: status, list: defaultList}
}

func (c column) Init() tea.Cmd {
	return nil
}

func (c column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var err error
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/margin, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Edit):
			if len(c.list.VisibleItems()) != 0 {
				task := c.list.SelectedItem().(Task)
				f := NewForm(db, task.TaskTitle, task.TaskDescription)
				f.index = c.list.Index()
				f.col = c
				if err := updateTask(db, task); err != nil {
					log.Fatal(err)
				}
				return f.Update(nil)
			}
		case key.Matches(msg, keys.New):
			f := newDefaultForm()
			f.index = APPEND
			f.col = c
			return f.Update(nil)
		case key.Matches(msg, keys.Delete):
			cmd, err = c.DeleteCurrent()
			if err != nil {
				log.Fatal(err)
			}
			return c, cmd
		case key.Matches(msg, keys.Enter):
			return c, c.MoveToNext()
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *column) DeleteCurrent() (tea.Cmd, error) {
	i := c.list.Index()

	if i >= 0 && i < len(c.tasks) {
		taskID := c.tasks[i].ID
		if err := delTask(db, taskID); err != nil {
			return nil, err
		}
		c.tasks = append(c.tasks[:i], c.tasks[i+1:]...)
		items, err := getTasksByStatus(db, c.status.String())
		if err != nil {
			log.Fatal(err)
		}
		c.list.SetItems(tasksToItems(items))

		var cmd tea.Cmd
		c.list, cmd = c.list.Update(nil)
		return cmd, nil
	}
	return nil, errors.New("index out of range")
}

func (c *column) Set(i int, t Task) tea.Cmd {
	if i != APPEND {
		return c.list.SetItem(i, t)
	}
	return c.list.InsertItem(APPEND, t)
}

func (c *column) setSize(width, height int) {
	c.width = width / margin
}

func (c *column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

type moveMsg struct {
	Task
}

func (c *column) MoveToNext() tea.Cmd {
	var task Task
	var ok bool
	if task, ok = c.list.SelectedItem().(Task); !ok {
		return nil
	}
	c.list.RemoveItem(c.list.Index())
	task.Status = c.status.Next().String()

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return moveMsg{task} })
}
