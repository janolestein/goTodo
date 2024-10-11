package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
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

var database *sql.DB

var kanbanModel *model

type confirmDeleteMsg bool

type status int

const (
	todo status = iota
	inProgress
	done
)

type task struct {
	title         string
	desc          string
	id            int
	prio          int
	currentStatus status
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
	// Connect to database
	dirname, err := os.UserHomeDir()
	os.MkdirAll(dirname + "/todoDB/uni", 0755)
	os.Create(dirname + "/todoDB/uni/tasks.db")
	db, err := sql.Open("sqlite3", dirname+"/todoDB/tasks.db")
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}

	database = db
	// defer close
	defer db.Close()
	lists := initalModel()
	kanbanModel = &lists
	p := tea.NewProgram(lists, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf(" Alas, there has been an Error: %v", err)
		os.Exit(1)
	}
}

func initalModel() model {

	stmt, err := database.Prepare("CREATE TABLE IF NOT EXISTS tasks (task_id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, desc TEXT, prio INTEGER, status INTEGER NOT NULL)")
	if err != nil {
		fmt.Println("Database could not be created")
		fmt.Println(err)
		panic(-1)
	}
	stmt.Exec()
	defer stmt.Close()
	itemsTodo := []list.Item{}
	itemsInProgress := []list.Item{}
	itemsDone := []list.Item{}

	tasks, _ := getAllTasks(database)

	for _, v := range tasks {

		switch v.currentStatus {
		case todo:
			itemsTodo = append(itemsTodo, v)
		case inProgress:
			itemsInProgress = append(itemsInProgress, v)
		case done:
			itemsDone = append(itemsDone, v)
		}
	}
	defList1 := list.New(itemsTodo, list.NewDefaultDelegate(), 0, 0)
	defList1.SetShowHelp(false)
	defList2 := list.New(itemsInProgress, list.NewDefaultDelegate(), 0, 0)
	defList2.SetShowHelp(false)
	defList3 := list.New(itemsDone, list.NewDefaultDelegate(), 0, 0)
	defList3.SetShowHelp(false)
	m := model{list: []list.Model{defList1, defList2, defList3}}
	m.list[todo].Title = "ToDo"
	m.list[inProgress].Title = "In Progress"
	m.list[done].Title = "Done"
	m.focused = todo

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

	sel := m.list[m.focused].SelectedItem()
	if sel != nil {
		selItem := sel.(task)
		index := m.list[m.focused].Index()
		m.list[m.focused].RemoveItem(index)
		if m.focused == todo {
			selItem.currentStatus = inProgress
			go updateTask(database, selItem)
			insertcmd := m.list[inProgress].InsertItem(len(m.list[inProgress].Items())-1, selItem)
			return insertcmd
		} else if m.focused == inProgress {
			selItem.currentStatus = done
			go updateTask(database, selItem)
			insertcmd := m.list[done].InsertItem(len(m.list[done].Items())-1, selItem)
			return insertcmd
		} else if m.focused == done {
			selItem.currentStatus = todo
			go updateTask(database, selItem)
			insertCmd := m.list[todo].InsertItem(len(m.list[todo].Items())-1, selItem)
			return insertCmd
		}
	}
	return nil
}

func (m *model) ConfirmDelete() tea.Msg {
	t := m.list[m.focused].SelectedItem()
	taskToDelete := t.(task)
	m.list[m.focused].RemoveItem(m.list[m.focused].Index())
	go deleteTask(database, taskToDelete.id)
	return nil
}

func (m *model) editTask(editedTask task, index int) tea.Cmd {
	return func() tea.Msg {
		editCmd := m.list[m.focused].SetItem(index, editedTask)
		go updateTask(database, editedTask)
		return editCmd
	}
}

func (m *model) newTask(newTask task) tea.Cmd {
	return func() tea.Msg {
		go insertNewTask(database, newTask)
		insertCmd := kanbanModel.list[todo].InsertItem(len(kanbanModel.list[todo].Items()), newTask)
		return insertCmd
	}
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list[m.focused].FilterState() != list.Filtering {

			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "right", "l":
				m.goToNext()
			case "left", "h":
				m.goToPrev()
			case "n":
				kanbanModel = &m
				return NewForm(), nil
			case "m":
				return m, m.moveToNext
			case "d":
				if m.list[m.focused].SelectedItem() != nil {

					kanbanModel = &m
					return NewConfirmForm(), nil
					// m.list[m.focused].RemoveItem(m.list[m.focused].Index())
				}
			case "e":
				if m.list[m.focused].SelectedItem() != nil {
					kanbanModel = &m
					return NewEditForm(m.list[m.focused].SelectedItem()), nil
				}
			}
		}

	case tea.WindowSizeMsg:
		// height, width := focusedModelStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		for i := range m.list {
			m.list[i].SetSize(msg.Width*4, msg.Height/2)
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
	s = lipgloss.JoinVertical(lipgloss.Center, s, m.list[m.focused].Help.View(keys))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
}
