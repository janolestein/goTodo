package main

import (
	// "fmt"
	// "os"
	//
	// "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type Form struct {
	title textinput.Model
	desc  textarea.Model
}

func NewForm() *Form {
	form := &Form{}
	form.title = textinput.New()
	form.title.Focus()
	form.desc = textarea.New()
    models = append(models, form)
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
		case "enter":
			if form.title.Focused() {
				form.title.Blur()
				form.desc.Focus()
				return form, textarea.Blink
			} else {
				models[inputModel] = form
                // title := form.title.Value()
                // desc := form.desc.Value()
				return models[listModel], nil
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
	return lipgloss.JoinVertical(lipgloss.Center, focusedModelStyle.Render(form.title.View()), focusedModelStyle.Render(form.desc.View()))
}
