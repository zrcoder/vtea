package vtea

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditorIntegration(t *testing.T) {
	initialContent := "Hello, world!"
	model := New(WithContent(initialContent))

	assert.Equal(t, ModeNormal, model.GetMode(), "Initial mode should be Normal")

	buffer := model.GetBuffer()
	assert.Equal(t, initialContent, buffer.Text(), "Buffer content should match initial content")

	model.SetMode(ModeInsert)
	assert.Equal(t, ModeInsert, model.GetMode(), "Mode should be Insert after setting")

	model.buffer.insertAt(0, 13, " This is a test.")

	expectedContent := "Hello, world! This is a test."
	assert.Equal(t, expectedContent, buffer.Text(), "Buffer content should match expected after insertion")

	model.SetMode(ModeNormal)
	assert.Equal(t, ModeNormal, model.GetMode(), "Mode should be Normal after setting")

	testStatusMsg := "Test status"
	cmd := model.SetStatusMessage(testStatusMsg)
	cmd()

	assert.Equal(t, testStatusMsg, model.statusMessage, "Status message should match set message")
}

func TestViewportIntegration(t *testing.T) {
	var content strings.Builder
	for i := range 30 {
		content.WriteString("Line " + string(rune('A'+i%26)) + "\n")
	}

	model := New(WithContent(content.String()))

	model.width = 80
	model.height = 20
	model.viewport.SetWidth(80)
	model.viewport.SetHeight(20)

	model.cursor = newCursor(25, 0)

	model.ensureCursorVisible()

	assert.GreaterOrEqual(t, model.cursor.Row, model.viewport.YOffset(),
		"Cursor row should be within or after viewport start")
	assert.Less(t, model.cursor.Row, model.viewport.YOffset()+model.height,
		"Cursor row should be within viewport end")
}

func TestKeyBindingsIntegration(t *testing.T) {
	model := New()

	testBindingCalled := false
	model.AddBinding(KeyBinding{
		Key:         "ctrl+t",
		Mode:        ModeNormal,
		Description: "Test binding",
		Handler: func(b Buffer) tea.Cmd {
			testBindingCalled = true
			return nil
		},
	})

	iBinding := model.registry.FindExact("i", ModeNormal)
	assert.NotNil(t, iBinding, "Default binding for 'i' should exist")

	ctrlTBinding := model.registry.FindExact("ctrl+t", ModeNormal)
	require.NotNil(t, ctrlTBinding, "Custom binding for 'ctrl+t' should exist")

	_ = ctrlTBinding.Command(model)

	assert.True(t, testBindingCalled, "Custom binding command should have been executed")
}

func TestCommandsIntegration(t *testing.T) {
	model := New()

	commandCalled := false
	model.AddCommand("test", func(b Buffer, args []string) tea.Cmd {
		commandCalled = true
		if len(args) > 0 && args[0] == "arg" {
			return nil
		}
		return nil
	})

	model.commandBuffer = "test arg"
	model.Update(CommandMsg{Command: "test"})

	assert.True(t, commandCalled, "Command should have been executed")
}

func TestClearIntegration(t *testing.T) {
	// Create editor with initial content
	initialContent := "Line 1\nLine 2\nLine 3"
	model := New(WithContent(initialContent))
	buffer := model.GetBuffer()

	// Verify initial content
	assert.Equal(t, initialContent, buffer.Text(), "Buffer should have initial content")
	assert.Equal(t, 3, buffer.LineCount(), "Buffer should have 3 lines initially")

	// Set cursor to a non-zero position
	model.cursor = newCursor(1, 3)

	// Clear the buffer using the Clear method
	clearCmd := buffer.Clear()
	if clearCmd != nil {
		clearCmd()
	}

	// After clearing, the buffer should have a single empty line and cursor at 0,0
	assert.Equal(t, 1, buffer.LineCount(), "Buffer should have 1 line after clear")
	assert.Equal(t, "", buffer.Text(), "Buffer text should be empty")
	assert.Equal(t, 0, model.cursor.Row, "Cursor row should be reset to 0")
	assert.Equal(t, 0, model.cursor.Col, "Cursor column should be reset to 0")

	// Test undo functionality after clear
	undoCmd := buffer.Undo()
	undoResult := undoCmd().(UndoRedoMsg)

	assert.True(t, undoResult.Success, "Undo after clear should succeed")
	assert.Equal(t, initialContent, buffer.Text(), "Buffer should return to initial content after undo")
}

func TestResetIntegration(t *testing.T) {
	// Create editor with initial content
	initialContent := "Initial content"
	model := New(WithContent(initialContent))
	buffer := model.GetBuffer()

	// Verify initial content
	assert.Equal(t, initialContent, buffer.Text(), "Buffer should have initial content")

	// Make changes to the editor
	buffer.InsertAt(0, 0, "Modified ") // Modify the content
	model.cursor = newCursor(0, 9)     // Move cursor after "Modified "

	// Verify the changes were made
	assert.Equal(t, "Modified Initial content", buffer.Text(), "Buffer content should be modified")
	assert.Equal(t, 0, model.cursor.Row, "Cursor row should be 0")
	assert.Equal(t, 9, model.cursor.Col, "Cursor column should be 9")

	// Reset the editor
	resetCmd := model.Reset()
	if resetCmd != nil {
		resetCmd()
	}

	// Verify the editor has been reset to initial state
	assert.Equal(t, initialContent, buffer.Text(), "Buffer should be reset to initial content")
	assert.Equal(t, 0, model.cursor.Row, "Cursor row should be reset to 0")
	assert.Equal(t, 0, model.cursor.Col, "Cursor column should be reset to 0")
	assert.Equal(t, ModeNormal, model.mode, "Editor mode should be reset to Normal")
	assert.Equal(t, "", model.yankBuffer, "Yank buffer should be empty")

	// Make more changes after reset
	buffer.InsertAt(0, 0, "New ")
	assert.Equal(t, "New Initial content", buffer.Text(), "Buffer should accept changes after reset")

	// Reset again
	resetCmd = model.Reset()
	if resetCmd != nil {
		resetCmd()
	}

	// Verify reset again
	assert.Equal(t, initialContent, buffer.Text(), "Buffer should be reset to initial content again")
}
