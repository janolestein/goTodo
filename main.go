package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title string
	desc  string
	id    int
	prio  int
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func main() {
	p := tea.NewProgram(initalModel(), tea.WithAltScreen())
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
    m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
    m.list.Title = "ToDo"

	return m 
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		height, width := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-height, msg.Height-width)

	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
