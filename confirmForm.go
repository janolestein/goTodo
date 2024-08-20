package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConfirmForm struct {
	title textinput.Model
}

func NewConfirmForm() *ConfirmForm {
	form := &ConfirmForm{}
	form.title = textinput.New()
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

				// return kanbanModel, confirmDelete
			}
		}
	}
	form.title, cmd = form.title.Update(msg)
	return form, cmd
}

func (form ConfirmForm) View() string {

	return "Are your sure you want to delete this task?\n" + focusedFormStyle.Render(form.title.View())
}
