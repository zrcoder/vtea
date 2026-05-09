package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zrcoder/vtea"
)

func main() {
	lineNumberStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Background(lipgloss.Color("#222222")).
		PaddingRight(1)

	currentLineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		Background(lipgloss.Color("#444444")).
		Bold(true).
		PaddingRight(1)

	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#CC8800")).
		Foreground(lipgloss.Color("black"))

	editor := vtea.NewEditor(
		vtea.WithLineNumberStyle(lineNumberStyle),
		vtea.WithCurrentLineNumberStyle(currentLineStyle),
		vtea.WithCursorStyle(cursorStyle),
		vtea.WithRelativeNumbers(true),
		vtea.WithFullScreen(),
	)

	editor.AddCommand("q", func(b vtea.Buffer, _ []string) tea.Cmd {
		return tea.Quit
	})

	p := tea.NewProgram(editor)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
