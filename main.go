package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var modelStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 1).
	Border(lipgloss.HiddenBorder())
var focusedModelStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 1).
	Border(lipgloss.RoundedBorder())

type focusedModel int

const (
	listModel focusedModel = iota
	inputModel
)

var kanbanModel *model

type status int

const (
	todo status = iota
	inProgress
	done
)

type item struct {
	title string
	desc  string
	id    int
	prio  int
}

func (i item) Title() string { return i.title }
func (i item) Description() string {
	if len(i.desc) > 25 {
		return i.desc[:25] + "..."
	} else {
		return i.desc
	}
}
func (i item) FilterValue() string { return i.title }

type model struct {
	list    []list.Model
	focused status
	width   int
	height  int
}

func main() {
	lists := initalModel()
	kanbanModel = &lists
	p := tea.NewProgram(lists, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf(" Alas, there has been an Error: %v", err)
		os.Exit(1)
	}
}

func initalModel() model {
	items := []list.Item{
		item{title: "learn Go", desc: "Try to learn some Golang", id: 1, prio: 3},
		item{title: "learn Rust", desc: "Try to learn some Rust", id: 2, prio: 2},
		item{title: "do some Leetcode", desc: "do some Leetcode Problems", id: 3, prio: 3},
	}
	defList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	defList.SetShowHelp(false)
	m := model{list: []list.Model{defList, defList, defList}}
	m.list[todo].Title = "ToDo"
	m.list[inProgress].Title = "In Progress"
	m.list[done].Title = "Done"
	m.focused = 0

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) goToNext() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *model) goToPrev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *model) moveToNext() tea.Msg {

	selItem := m.list[m.focused].SelectedItem()
	if selItem != nil {
        selItem = selItem.(item)
		index := m.list[m.focused].Index()
		m.list[m.focused].RemoveItem(index)
		if m.focused == todo {
			m.list[inProgress].InsertItem(len(m.list[inProgress].Items())-1, selItem)
		} else if m.focused == inProgress {
			m.list[done].InsertItem(len(m.list[done].Items())-1, selItem)
		} else if m.focused == done {
			m.list[todo].InsertItem(len(m.list[todo].Items())-1, selItem)
		}
	}
    m.goToNext()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "right", "l":
			m.goToNext()
		case "left", "h":
			m.goToPrev()
		case "n":
			return NewForm(), nil
		case "m":
			return m, m.moveToNext
		}

	case tea.WindowSizeMsg:
		height, width := focusedModelStyle.GetFrameSize()
		m.width = width
		m.height = height
		for i := range m.list {
			m.list[i].SetSize(msg.Width, msg.Height-height)
		}

	}
	var cmd tea.Cmd

	m.list[m.focused], cmd = m.list[m.focused].Update(msg)
	return m, cmd
}

func (m model) View() string {
	var views []string
	for i := range m.list {
		if int(m.focused) == i {
			views = append(views, focusedModelStyle.Render(m.list[i].View()))
		} else {
			views = append(views, modelStyle.Render(m.list[i].View()))
		}
	}
	s := lipgloss.JoinHorizontal(lipgloss.Center, views...) + "\n\n"

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
}
