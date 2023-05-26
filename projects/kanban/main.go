package kanban

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status = int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

type Task struct {
	status      status
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type Model struct {
	focused status
	lists   []list.Model
	err     error
	loaded  bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	// Init todo list
	m.lists[todo].Title = "Todo"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "got milk?"},
		Task{status: todo, title: "buy cookies", description: "the good biscuits"},
		Task{status: todo, title: "buy coke", description: "not white powder variety"},
	})
	// Init in progress list
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "got milk?"},
		Task{status: todo, title: "buy cookies", description: "the good biscuits"},
	})
	// Init done list
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "buy cookies", description: "the good biscuits"},
		Task{status: todo, title: "buy coke", description: "not white powder variety"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loaded {
		return lipgloss.JoinHorizontal(lipgloss.Left, m.lists[todo].View(), m.lists[inProgress].View(), m.lists[done].View())
	} else {
		return "loading..."
	}
}

func Run() {
	m := New()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
