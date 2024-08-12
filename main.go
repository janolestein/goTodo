package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func main(){
    p := tea.NewProgram(initalModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf(" Alas, there has been an Error: %v", err)
        os.Exit(1)
    }
}

func initalModel() model {
	return model{
		choices:  []string{"buy duplos", "buy pasta", "buy pesto"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }


		}
	}
    return m, nil
}

func (m model) View() string {
    s := "What i need to buy\n\n"

    for i, choices := range m.choices {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }

        checked := " "
        if _, ok := m.selected[i]; ok {
            checked = "x"
        }

        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choices)
    }

    s += "\nPress q to quit.\n"

    return s

}
