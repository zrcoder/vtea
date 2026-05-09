package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/zrcoder/vtea"
)

func main() {
	editor := vtea.NewEditor(vtea.WithFullScreen())

	editor.AddBinding(vtea.KeyBinding{
		Key:         "ctrl+s",
		Mode:        vtea.ModeNormal,
		Description: "Save file",
		Handler: func(b vtea.Buffer) tea.Cmd {
			return vtea.SetStatusMsg("File saved!")
		},
	})
	editor.AddBinding(vtea.KeyBinding{
		Key:         "ctrl+c",
		Mode:        vtea.ModeNormal,
		Description: "quit editor",
		Handler: func(b vtea.Buffer) tea.Cmd {
			return tea.Quit
		},
	})

	p := tea.NewProgram(editor)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
