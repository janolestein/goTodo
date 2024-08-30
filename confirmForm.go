package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmForm struct {
	title textinput.Model
}

func NewConfirmForm() *ConfirmForm {
	form := &ConfirmForm{}
	form.title = textinput.New()
	form.title.Focus()
	form.title.Placeholder = "y/n"
	return form
}

func (form ConfirmForm) Init() tea.Cmd {
	return nil
}

func (form ConfirmForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return form, tea.Quit
		case "esc":
			return kanbanModel, nil
		case "enter":
			if form.title.Value() == "y" {
				return kanbanModel, kanbanModel.ConfirmDelete
			} else {
				return kanbanModel, nil
			}
		}
	}
	form.title, cmd = form.title.Update(msg)
	return form, cmd
}

func (form ConfirmForm) View() string {

	s := "Are your sure you want to delete this task?\n" + focusedFormStyle.Render(form.title.View())
	return lipgloss.Place(kanbanModel.width, kanbanModel.height, lipgloss.Center, lipgloss.Center, s)
}
