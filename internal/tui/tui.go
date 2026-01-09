package tui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type job struct {
	id     int
	name   string
	status string
}

type model struct {
	jobs   []job
	cursor int
	nextID int
}

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)
)

func Run() error {
	m := model{
		jobs: []job{
			{id: 1, name: "nmap scan example.com", status: "queued"},
			{id: 2, name: "dns lookup example.com", status: "running"},
		},
		cursor: 0,
		nextID: 3,
	}
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
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
		case "j", "down":
			if m.cursor < len(m.jobs)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "a":
			m.jobs = append(m.jobs, job{
				id:     m.nextID,
				name:   "new job " + strconv.Itoa(m.nextID),
				status: "queued",
			})
			m.nextID++
		case "enter":
			if len(m.jobs) > 0 {
				m.jobs[m.cursor].status = nextStatus(m.jobs[m.cursor].status)
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	out := headerStyle.Render("CLI Tools Dashboard") + "\n"
	out += "q: quit  a: add job  enter: toggle status\n\n"

	for i, j := range m.jobs {
		cursor := " "
		if i == m.cursor {
			cursor = cursorStyle.Render(">")
		}
		line := fmt.Sprintf("%s [%s] %s", cursor, statusStyle.Render(j.status), j.name)
		out += line + "\n"
	}
	return out
}

func nextStatus(current string) string {
	switch current {
	case "queued":
		return "running"
	case "running":
		return "done"
	default:
		return "queued"
	}
}
