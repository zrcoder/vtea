// Package vtea provides a Vim-like text editor component for terminal applications
package vtea

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

// Default styles for the editor components
// These can be overridden using the With* option functions
var (
	// lineNumberStyle defines the appearance of regular line numbers
	lineNumberStyle = lipgloss.NewStyle().
			Foreground(compat.AdaptiveColor{Light: lipgloss.Color("245"), Dark: lipgloss.Color("242")}).
			Bold(false).
			PaddingRight(1)

	// currentLineNumberStyle defines the appearance of the current line number
	currentLineNumberStyle = lipgloss.NewStyle().
				Foreground(compat.AdaptiveColor{Light: lipgloss.Color("0"), Dark: lipgloss.Color("15")}).
				Bold(true).
				Background(compat.AdaptiveColor{Light: lipgloss.Color("252"), Dark: lipgloss.Color("236")}).
				PaddingRight(1)

	// textStyle defines the appearance of regular text in the editor
	textStyle = lipgloss.NewStyle()

	// statusStyle defines the appearance of the status bar
	statusStyle = lipgloss.NewStyle().
			Foreground(compat.AdaptiveColor{Light: lipgloss.Color("7"), Dark: lipgloss.Color("8")}).
			Background(compat.AdaptiveColor{Light: lipgloss.Color("8"), Dark: lipgloss.Color("7")})

	// cursorStyle defines the appearance of the cursor
	cursorStyle = lipgloss.NewStyle().
			Background(compat.AdaptiveColor{Light: lipgloss.Color("252"), Dark: lipgloss.Color("248")}).
			Foreground(compat.AdaptiveColor{Light: lipgloss.Color("0"), Dark: lipgloss.Color("0")})

	// commandStyle defines the appearance of the command line
	commandStyle = lipgloss.NewStyle().
			Foreground(compat.AdaptiveColor{Light: lipgloss.Color("3"), Dark: lipgloss.Color("3")}).
			Bold(true)

	// selectedStyle defines the appearance of selected text in visual mode
	selectedStyle = lipgloss.NewStyle().Background(
		compat.AdaptiveColor{Light: lipgloss.Color("7"), Dark: lipgloss.Color("8")},
	)
)
