package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/zrcoder/vtea"
)

func main() {
	content := `This is a sample file
     with multiple lines
     for testing the editor`

	editor := vtea.NewEditor(
		vtea.WithContent(content),
		vtea.WithFileName("example.txt"),
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
