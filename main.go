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
    BorderForeground(lipgloss.Color("201")).
	Border(lipgloss.RoundedBorder())

type focusedModel int

const (
	listModel focusedModel = iota
	inputModel
)

var kanbanModel *model

type confirmDeleteMsg bool


type status int

const (
	todo status = iota
	inProgress
	done
)

type task struct {
	title string
	desc  string
	id    int
	prio  int
}

func (i task) Title() string { return i.title }
func (i task) Description() string {
	if len(i.desc) > 25 {
		return i.desc[:25] + "..."
	} else {
		return i.desc
	}
}
func (i task) FilterValue() string { return i.title }

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
	items1 := []list.Item{
		task{title: "learn Go", desc: "Try to learn some Golang", id: 1, prio: 3},
		task{title: "learn Rust", desc: "Try to learn some Rust", id: 2, prio: 2},
		task{title: "do some Leetcode", desc: "do some Leetcode Problems", id: 3, prio: 3},
	}
	items2 := []list.Item{
		task{title: "sjdkfhsd", desc: "sfjsdlfjsdklj", id: 1, prio: 3},
		task{title: "learn Rust", desc: "sdklfjsdk", id: 2, prio: 2},
		task{title: "kfsjdkjcxvb", desc: "ksldjf", id: 3, prio: 3},
	}
	items3 := []list.Item{
		task{title: "rwueor", desc: "wueiprt", id: 1, prio: 3},
		task{title: "348975348", desc: "ksdchrgb", id: 2, prio: 2},
		task{title: "348hjeskhf", desc: "dfhosbajskbcasb", id: 3, prio: 3},
	}
	defList1 := list.New(items1, list.NewDefaultDelegate(), 0, 0)
	defList2 := list.New(items2, list.NewDefaultDelegate(), 0, 0)
	defList3 := list.New(items3, list.NewDefaultDelegate(), 0, 0)
	// defList.SetShowHelp(false)
	m := model{list: []list.Model{defList1, defList2, defList3}}
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
		selItem = selItem.(task)
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
			selItem := m.list[m.focused].SelectedItem()
			if selItem != nil {
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
			// var cmds []tea.Cmd
			//
			//    for i := range m.list {
			//        var cmd tea.Cmd
			//        m.list[i], cmd = m.list[i].Update(msg)
			//        cmds = append(cmds, cmd)
			//    }
   //          return m, tea.Batch(cmds...)
        case "d":
            m.list[m.focused].RemoveItem(m.list[m.focused].Index())
		}

	case tea.WindowSizeMsg:
		height, width := focusedModelStyle.GetFrameSize()
		m.width = width
		m.height = height
        focusedModelStyle.MarginTop(height / 2)
        modelStyle.MarginTop(height / 2)
        focusedModelStyle.MarginLeft(width / 2)
        modelStyle.MarginLeft(width / 2)
		for i := range m.list {
			m.list[i].SetSize(msg.Width, msg.Height / 2)
		}

	}
	// var cmds []tea.Cmd
	//
	//    for i := range m.list {
	//        var cmd tea.Cmd
	//        m.list[i], cmd = m.list[i].Update(msg)
	//        cmds = append(cmds, cmd)
	//    }
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

	return lipgloss.Place(m.width, m.height, lipgloss.Right, lipgloss.Bottom, s)
}
