package addForm

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

var models *[]tea.Model

var modelStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 1).
	Border(lipgloss.HiddenBorder())
var focusedModelStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 1).
	Border(lipgloss.RoundedBorder())

type Form struct {
	title textinput.Model
	desc  textarea.Model
}

func NewForm(allModels *[]tea.Model) *Form {
	models = allModels
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
		case "enter":
			if form.title.Focused() {
				form.title.Blur()
				form.desc.Focus()
				return form, textarea.Blink
			} else {
				(*models)[1] = form
				return (*models)[0], nil
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
