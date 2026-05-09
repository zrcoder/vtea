// Example application demonstrating the use of vtea
// This opens itself and provides a Vim-like interface to edit the file
package main

import (
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/zrcoder/vtea"
)

func main() {
	file, err := os.Open("example/main.go")
	if err != nil {
		log.Fatalf("Failed to open example/main.go: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file stat: %v", err)
	}
	// Read the file
	buf := make([]byte, stat.Size())
	_, err = file.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Create a new editor with the file contents
	// WithFileName is used for syntax highlighting
	editor := vtea.NewEditor(
		vtea.WithContent(string(buf)),
		vtea.WithFileName("example/main.go"),
		vtea.WithFullScreen(),
	)

	// Add a custom key binding for quitting with Ctrl+C
	editor.AddBinding(vtea.KeyBinding{
		Key:         "ctrl+c",
		Mode:        vtea.ModeNormal,
		Description: "Close the editor",
		Handler: func(b vtea.Buffer) tea.Cmd {
			return tea.Quit
		},
	})

	// Add a custom command that can be invoked with :q
	editor.AddCommand("q", func(b vtea.Buffer, _ []string) tea.Cmd {
		return tea.Quit
	})

	p := tea.NewProgram(editor)
	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
	}
}
