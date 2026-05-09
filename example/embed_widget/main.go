package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zrcoder/vtea"
)

func main() {
	m := NewModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

const code = `package main

import "fmt"

func main() {
	fmt.Println("hi")
}`

type runMsg string

func NewModel() model {
	m := model{sample: "Hello, My editor"}
	editor := vtea.New(
		vtea.WithFullScreen(),
		vtea.WithFileName("main.go"),
		vtea.WithContent(code),
	)
	editor.AddCommand("run", func(b vtea.Buffer, s []string) tea.Cmd {
		return func() tea.Msg {
			return runMsg(b.Text())
		}
	})
	editor.AddCommand("q", func(b vtea.Buffer, _ []string) tea.Cmd {
		return tea.Quit
	})
	m.editor = editor
	return m
}

type model struct {
	sample string
	editor *vtea.Model
}

func (m model) Init() tea.Cmd {
	return m.editor.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_, cmd = m.editor.Update(msg)
	switch msg := msg.(type) {
	case runMsg:
		m.sample = string(msg)
	}
	return m, cmd
}

func (m model) View() tea.View {
	view := tea.NewView(lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.sample,
		"",
		m.editor.View().Content,
	))
	view.AltScreen = true
	return view
}
