package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/zrcoder/vtea"
)

func main() {
	editor := vtea.NewEditor(vtea.WithFullScreen())
	editor.AddCommand("q", func(b vtea.Buffer, _ []string) tea.Cmd {
		return tea.Quit
	})

	p := tea.NewProgram(editor)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
