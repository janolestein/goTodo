package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var formStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 1).
	Width(50).
	Border(lipgloss.HiddenBorder())
var focusedFormStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 1).
	Width(50).
	BorderForeground(lipgloss.Color("201")).
	Border(lipgloss.RoundedBorder())

type Form struct {
	title textinput.Model
	desc  textarea.Model
}

func NewForm() *Form {
	form := &Form{}
	form.title = textinput.New()
	form.title.Focus()
	form.desc = textarea.New()
	return form
}

func (form Form) Init() tea.Cmd {
	return nil
}

func (form Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return form, tea.Quit
		case "esc":
			return kanbanModel, nil
		case "enter":
			if form.title.Focused() {
				form.title.Blur()
				form.desc.Focus()
				return form, textarea.Blink
			} else {
				title := form.title.Value()
				desc := form.desc.Value()
				if title != "" {
					return kanbanModel, kanbanModel.newTask(task{title: title, desc: desc, prio: 1, currentStatus: todo})
				} else {
					return kanbanModel, nil
				}
			}
		}
	}
	if form.title.Focused() {
		form.title, cmd = form.title.Update(msg)
		return form, cmd
	} else {
		form.desc, cmd = form.desc.Update(msg)
		return form, cmd
	}
}

func (form Form) View() string {
	var s string
	if form.title.Focused() {
		s = lipgloss.JoinVertical(lipgloss.Center, "Please Enter a New Task\n", focusedFormStyle.Render(form.title.View()), formStyle.Render(form.desc.View()))
		return lipgloss.Place(kanbanModel.width, kanbanModel.height, lipgloss.Center, lipgloss.Center, s)
	} else {
		s = lipgloss.JoinVertical(lipgloss.Center, "Please Enter a New Task\n", formStyle.Render(form.title.View()), focusedFormStyle.Render(form.desc.View()))
		return lipgloss.Place(kanbanModel.width, kanbanModel.height, lipgloss.Center, lipgloss.Center, s)
	}
}
