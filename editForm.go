package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type editForm struct {
	title      textinput.Model
	desc       textarea.Model
	itemToEdit list.Item
}

func NewEditForm(taskToEdit list.Item) *editForm {
	form := &editForm{}
	form.title = textinput.New()
	form.title.Focus()
	form.desc = textarea.New()
	form.itemToEdit = taskToEdit
	form.title.SetValue(taskToEdit.(task).title)
	form.desc.SetValue(taskToEdit.(task).desc)
	return form
}

func (form editForm) Init() tea.Cmd {
	return nil
}

func (form editForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				newTask := task{title: title, desc: desc, id: form.itemToEdit.(task).id, prio: form.itemToEdit.(task).prio, currentStatus: form.itemToEdit.(task).currentStatus}
				index := kanbanModel.list[kanbanModel.focused].Index()
				return kanbanModel, kanbanModel.editTask(newTask, index)
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

func (form editForm) View() string {
	var s string
	if form.title.Focused() {
		s = lipgloss.JoinVertical(lipgloss.Center, "Please Enter want you want to change\n", focusedFormStyle.Render(form.title.View()), formStyle.Render(form.desc.View()))
		return lipgloss.Place(kanbanModel.width, kanbanModel.height, lipgloss.Center, lipgloss.Center, s)
	} else {
		s = lipgloss.JoinVertical(lipgloss.Center, "Please Enter want you want to change\n", formStyle.Render(form.title.View()), focusedFormStyle.Render(form.desc.View()))
		return lipgloss.Place(kanbanModel.width, kanbanModel.height, lipgloss.Center, lipgloss.Center, s)
	}
}
