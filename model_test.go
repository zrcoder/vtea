package vtea

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelInitialization(t *testing.T) {
	// Test with default options
	model := New()

	assert.Equal(t, ModeNormal, model.mode, "Initial mode should be Normal")
	assert.NotNil(t, model.buffer, "Buffer should be initialized")
	assert.NotNil(t, model.registry, "Binding registry should be initialized")
	assert.NotNil(t, model.commands, "Command registry should be initialized")

	// Test with content option
	testContent := "Test content"
	model = New(WithContent(testContent))

	assert.Equal(t, testContent, model.buffer.text(), "Buffer content should be initialized with provided content")

	// Test with filename option
	model = New(WithContent(""), WithFileName("test.go"))

	assert.NotNil(t, model.highlighter, "Syntax highlighter should be initialized")
}

func TestModelUpdate(t *testing.T) {
	model := New(WithFullScreen())

	// Test key message handling
	keyMsg := tea.KeyPressMsg{Code: 'i'}
	updatedModel, _ := model.Update(keyMsg)

	assert.Equal(t, ModeInsert, updatedModel.mode, "After pressing 'i', mode should be Insert")

	// Test window resize message
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ = model.Update(sizeMsg)

	assert.Equal(t, 100, updatedModel.width, "Window width should be updated correctly")
	assert.Equal(t, 48, updatedModel.height, "Window height should be adjusted for status bar") // height is adjusted for status bar

	// Test setting status message
	statusCmd := model.SetStatusMessage("Test status")
	statusCmd()

	assert.Equal(t, "Test status", model.statusMessage, "Status message should be updated correctly")
}

func TestModelKeySequences(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2\nLine 3"))

	// Get the dd (delete line) binding
	binding := model.registry.FindExact("dd", ModeNormal)
	require.NotNil(t, binding, "Built-in binding for 'dd' should exist")

	// First 'd' key press
	keyMsg := tea.KeyPressMsg{Code: 'd'}
	updatedModel, _ := model.Update(keyMsg)

	// Should be collecting key sequence
	assert.Len(t, updatedModel.keySequence, 1, "Key sequence should have 1 key after first 'd'")
	assert.Equal(t, "d", updatedModel.keySequence[0], "First key in sequence should be 'd'")

	// Second 'd' key press
	updatedModel, _ = updatedModel.Update(keyMsg)

	// Sequence should be executed and cleared
	assert.Empty(t, updatedModel.keySequence, "Key sequence should be cleared after command execution")

	// Buffer should be updated (line deleted)
	assert.Equal(t, 2, updatedModel.buffer.lineCount(), "Buffer should have 2 lines after deletion")

	// First line should now be what was previously the second line
	assert.Equal(t, "Line 2", updatedModel.buffer.Line(0), "After deleting first line, new first line should be 'Line 2'")
}

func TestModelCountPrefix(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2\nLine 3\nLine 4\nLine 5"))

	// Press '3'
	keyMsg := tea.KeyPressMsg{Code: '3'}
	updatedModel, _ := model.Update(keyMsg)

	assert.Equal(t, 3, updatedModel.countPrefix, "Count prefix should be 3")

	// Press 'j' to move down 3 lines
	keyMsg = tea.KeyPressMsg{Code: 'j'}
	updatedModel, _ = updatedModel.Update(keyMsg)

	assert.Equal(t, 3, updatedModel.cursor.Row, "Cursor should move down 3 lines to row 3")

	// Count prefix should be reset
	assert.Equal(t, 1, updatedModel.countPrefix, "Count prefix should be reset after use")

	// Test multi-digit count
	keyMsg = tea.KeyPressMsg{Code: '1'}
	updatedModel, _ = updatedModel.Update(keyMsg)

	keyMsg = tea.KeyPressMsg{Code: '2'}
	updatedModel, _ = updatedModel.Update(keyMsg)

	assert.Equal(t, 12, updatedModel.countPrefix, "Count prefix should be 12")
}

func TestModelCommandMode(t *testing.T) {
	model := New()

	// Enter command mode
	keyMsg := tea.KeyPressMsg{Code: ':'}
	updatedModel, _ := model.Update(keyMsg)

	assert.Equal(t, ModeCommand, updatedModel.mode, "Mode should be Command after pressing ':'")

	// Type command
	for _, ch := range "test" {
		keyMsg = tea.KeyPressMsg{Code: ch}
		updatedModel, _ = updatedModel.Update(keyMsg)
	}

	assert.Equal(t, "test", updatedModel.commandBuffer, "Command buffer should contain 'test'")

	// Register test command
	commandCalled := false
	model.AddCommand("test", func(b Buffer, args []string) tea.Cmd {
		commandCalled = true
		return nil
	})

	// Execute command with Enter
	keyMsg = tea.KeyPressMsg{Code: tea.KeyEnter}
	updatedModel, cmd := updatedModel.Update(keyMsg)
	for cmd != nil {
		updatedModel, cmd = updatedModel.Update(cmd())
	}

	assert.True(t, commandCalled, "Command should have been called")
	assert.Equal(t, ModeNormal, updatedModel.mode, "Mode should return to Normal after command execution")

	// Test command backspace
	updatedModel.mode = ModeCommand
	updatedModel.commandBuffer = "test"

	keyMsg = tea.KeyPressMsg{Code: tea.KeyBackspace}
	updatedModel, _ = updatedModel.Update(keyMsg)

	assert.Equal(t, "tes", updatedModel.commandBuffer, "Command buffer should be 'tes' after backspace")
}

func TestModelVisualMode(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2\nLine 3"))

	// Enter visual mode
	keyMsg := tea.KeyPressMsg{Code: 'v'}
	updatedModel, _ := model.Update(keyMsg)

	assert.Equal(t, ModeVisual, updatedModel.mode, "Mode should be Visual after pressing 'v'")
	assert.Equal(t, 0, updatedModel.visualStart.Row, "Visual start row should be 0")
	assert.Equal(t, 0, updatedModel.visualStart.Col, "Visual start column should be 0")

	// Move cursor to create selection
	keyMsg = tea.KeyPressMsg{Code: 'j'}
	updatedModel, _ = updatedModel.Update(keyMsg)

	// Check selection boundaries
	start, end := updatedModel.GetSelectionBoundary()

	assert.Equal(t, 0, start.Row, "Selection start row should be 0")
	assert.Equal(t, 1, end.Row, "Selection end row should be 1")

	// Test yank in visual mode
	keyMsg = tea.KeyPressMsg{Code: 'y'}
	updatedModel, _ = updatedModel.Update(keyMsg)

	assert.Equal(t, ModeNormal, updatedModel.mode, "Mode should return to Normal after yanking")
	assert.Contains(t, updatedModel.yankBuffer, "Line 1", "Yank buffer should contain 'Line 1'")
}

func TestModelInsertMode(t *testing.T) {
	model := New(WithContent("Line 1"))

	// Enter insert mode
	keyMsg := tea.KeyPressMsg{Code: 'i'}
	updatedModel, _ := model.Update(keyMsg)

	assert.Equal(t, ModeInsert, updatedModel.mode, "Mode should be Insert after pressing 'i'")
	model.cursor = newCursor(0, 6)

	// Type some text
	for _, ch := range " inserted" {
		keyMsg = tea.KeyPressMsg{Code: ch}
		updatedModel, _ = updatedModel.Update(keyMsg)
	}

	expectedText := "Line 1 inserted"
	assert.Equal(t, expectedText, updatedModel.buffer.text(), "Buffer content should match expected after insertion")

	// Exit insert mode
	keyMsg = tea.KeyPressMsg{Code: tea.KeyEsc}
	updatedModel, _ = updatedModel.Update(keyMsg)

	assert.Equal(t, ModeNormal, updatedModel.mode, "Mode should be Normal after pressing Escape")
}

func TestEditorOptions(t *testing.T) {
	// Test multiple options
	model := New(
		WithContent("Test content"),
		WithFileName("test.go"),
		WithEnableStatusBar(false),
		WithBlinkInterval(200*time.Millisecond),
	)

	assert.Equal(t, "Test content", model.buffer.text(), "WithContent option should be applied correctly")
	assert.False(t, model.enableStatusBar, "WithEnableStatusBar(false) option should be applied correctly")
	assert.Equal(t, 200*time.Millisecond, model.blinkInterval, "WithBlinkInterval option should be applied correctly")

	// Test disabling command mode
	model = New(WithEnableModeCommand(false))

	assert.False(t, model.enableCommandMode, "WithEnableModeCommand(false) option should be applied correctly")

	// Test enabling relative line numbers
	model = New(WithRelativeNumbers(true))

	assert.True(t, model.relativeNumbers, "WithRelativeNumbers option should be applied correctly")
}
