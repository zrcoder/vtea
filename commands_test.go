package vtea

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandExecution(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2\nLine 3"))

	// Test yank line command (find binding in registry)
	binding := model.registry.FindExact("yy", ModeNormal)
	require.NotNil(t, binding, "Binding for 'yy' not found")

	model.cursor = newCursor(1, 0)
	binding.Command(model)

	assert.Contains(t, model.yankBuffer, "Line 2", "yankLine should set yankBuffer to contain 'Line 2'")

	// Test delete line command
	deleteBinding := model.registry.FindExact("dd", ModeNormal)
	require.NotNil(t, deleteBinding, "Binding for 'dd' not found")

	model.cursor = newCursor(1, 0)
	deleteBinding.Command(model)

	assert.Equal(t, 2, model.buffer.lineCount(), "deleteLine should remove a line")
	assert.Equal(t, "Line 3", model.buffer.Line(1), "After deletion, line 1 should be 'Line 3'")
}

func TestPasteCommands(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2\nLine 3"))

	// Set up yankBuffer
	model.yankBuffer = "Yanked content"

	// Test paste after command
	pasteAfterBinding := model.registry.FindExact("p", ModeNormal)
	require.NotNil(t, pasteAfterBinding, "Binding for 'p' not found")

	model.cursor = newCursor(0, 5)
	pasteAfterBinding.Command(model)

	expectedContent := "Line 1Yanked content\nLine 2\nLine 3"
	assert.Equal(t, expectedContent, model.buffer.text(), "pasteAfter should insert at cursor position")

	// Test paste before command
	pasteBeforeBinding := model.registry.FindExact("P", ModeNormal)
	require.NotNil(t, pasteBeforeBinding, "Binding for 'P' not found")

	// Reset buffer
	model.buffer.lines = []string{"Line 1", "Line 2", "Line 3"}
	model.cursor = newCursor(0, 5)
	pasteBeforeBinding.Command(model)

	expectedContent = "Line Yanked content1\nLine 2\nLine 3"
	assert.Equal(t, expectedContent, model.buffer.text(), "pasteBefore should insert at cursor position")

	// Test line-wise paste
	// Reset buffer
	model.buffer.lines = []string{"Line 1", "Line 2", "Line 3"}
	model.yankBuffer = "\nYanked line"
	model.cursor = newCursor(1, 0)
	pasteAfterBinding.Command(model)

	expectedContent = "Line 1\nLine 2\nYanked line\nLine 3"
	assert.Equal(t, expectedContent, model.buffer.text(), "pasteAfter with line-wise content should insert as new line")
}

func TestInsertModeCommands(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2"))

	// Test insert at beginning of line (I command)
	insertStartBinding := model.registry.FindExact("I", ModeNormal)
	require.NotNil(t, insertStartBinding, "Binding for 'I' not found")

	model.cursor = newCursor(1, 2)
	insertStartBinding.Command(model)

	assert.Equal(t, 0, model.cursor.Col, "I command should move cursor to col 0")
	assert.Equal(t, ModeInsert, model.mode, "I command should switch to insert mode")

	// Test insert at end of line (A command)
	appendEndBinding := model.registry.FindExact("A", ModeNormal)
	require.NotNil(t, appendEndBinding, "Binding for 'A' not found")

	model.mode = ModeNormal
	model.cursor = newCursor(1, 2)
	appendEndBinding.Command(model)

	assert.Equal(t, 6, model.cursor.Col, "A command should move cursor to end of line")
	assert.Equal(t, ModeInsert, model.mode, "A command should switch to insert mode")

	// Test insert new line below (o command)
	openBelowBinding := model.registry.FindExact("o", ModeNormal)
	require.NotNil(t, openBelowBinding, "Binding for 'o' not found")

	model.mode = ModeNormal
	model.cursor = newCursor(0, 0)
	openBelowBinding.Command(model)

	assert.Equal(t, 3, model.buffer.lineCount(), "o command should add a new line")
	assert.Equal(t, 1, model.cursor.Row, "o command should position cursor at new line row")
	assert.Equal(t, 0, model.cursor.Col, "o command should position cursor at start of new line")
}

func TestOpenLineBelowIndentation(t *testing.T) {
	model := New(WithContent("    indented line\nnext"))

	openBelowBinding := model.registry.FindExact("o", ModeNormal)
	require.NotNil(t, openBelowBinding, "Binding for 'o' not found")

	model.mode = ModeNormal
	model.cursor = newCursor(0, 5)
	openBelowBinding.Command(model)

	assert.Equal(t, "    indented line", model.buffer.Line(0))
	assert.Equal(t, "    ", model.buffer.Line(1), "o should preserve leading whitespace")
}

func TestOpenLineAboveIndentation(t *testing.T) {
	model := New(WithContent("    first\n    second"))

	openAboveBinding := model.registry.FindExact("O", ModeNormal)
	require.NotNil(t, openAboveBinding, "Binding for 'O' not found")

	model.mode = ModeNormal
	model.cursor = newCursor(1, 0) // on line "    second"
	openAboveBinding.Command(model)

	assert.Equal(t, 3, model.buffer.lineCount(), "O command should add a new line")
	// O inserts new line ABOVE cursor row (row 1), using indentation from line ABOVE cursor
	// Line(0) = "    first" (unchanged), Line(1) = new line with indent from line 0, Line(2) = "    second" (shifted)
	assert.Equal(t, "    first", model.buffer.Line(0), "Line 0 unchanged")
	assert.Equal(t, "    ", model.buffer.Line(1), "O should insert line with indentation from above")
	assert.Equal(t, "    second", model.buffer.Line(2), "Line 2 shifted down")
	assert.Equal(t, 1, model.cursor.Row, "O should position cursor at new line row")
}

func TestInsertEnterKeyIndentation(t *testing.T) {
	t.Run("preserves leading whitespace when splitting line", func(t *testing.T) {
		model := New(WithContent("    some code here"))

		enterBinding := model.registry.FindExact("enter", ModeInsert)
		require.NotNil(t, enterBinding, "Enter binding not found in insert mode")

		model.mode = ModeInsert
		model.cursor = newCursor(0, 8) // cursor on the space after "some"
		enterBinding.Command(model)

		assert.Equal(t, 2, model.buffer.lineCount())
		assert.Equal(t, "    some", model.buffer.Line(0), "First line ends at cursor")
		// New line: leading whitespace (4 spaces) + remainder (" code here")
		assert.Equal(t, "     code here", model.buffer.Line(1), "New line: indentation + remainder")
		assert.Equal(t, 1, model.cursor.Row)
		assert.Equal(t, 4, model.cursor.Col, "Cursor should be positioned after leading whitespace")
	})

	t.Run("preserves tab indentation as spaces", func(t *testing.T) {
		model := New(WithContent("        if true {"))

		enterBinding := model.registry.FindExact("enter", ModeInsert)
		require.NotNil(t, enterBinding)

		model.mode = ModeInsert
		model.cursor = newCursor(0, 8) // at position 8, right after the 8 leading spaces
		enterBinding.Command(model)

		// Line(0) = "        " (8 spaces, text before cursor)
		// Line(1) = "        if true {" (8 leading spaces + "if true {")
		assert.Equal(t, "        ", model.buffer.Line(0))
		assert.Equal(t, "        if true {", model.buffer.Line(1))
		assert.Equal(t, 8, model.cursor.Col, "Cursor should be after 8 leading spaces")
	})

	t.Run("cursor at end of indented line adds only indentation", func(t *testing.T) {
		model := New(WithContent("    completed line"))

		enterBinding := model.registry.FindExact("enter", ModeInsert)
		require.NotNil(t, enterBinding)

		model.mode = ModeInsert
		model.cursor = newCursor(0, 18) // past last char = end of line (18 chars total)
		enterBinding.Command(model)

		assert.Equal(t, "    completed line", model.buffer.Line(0))
		assert.Equal(t, "    ", model.buffer.Line(1), "New line should have only the indentation")
		assert.Equal(t, 1, model.cursor.Row)
		assert.Equal(t, 4, model.cursor.Col)
	})

	t.Run("unindented line starts at column 0", func(t *testing.T) {
		model := New(WithContent("no indent here"))

		enterBinding := model.registry.FindExact("enter", ModeInsert)
		require.NotNil(t, enterBinding)

		model.mode = ModeInsert
		model.cursor = newCursor(0, 5)
		enterBinding.Command(model)

		assert.Equal(t, "no in", model.buffer.Line(0))
		assert.Equal(t, "dent here", model.buffer.Line(1), "New line should have no leading whitespace")
		assert.Equal(t, 0, model.cursor.Col)
	})
}

func TestCursorMovementCommands(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2\nLine 3"))

	// Test move down (j)
	downBinding := model.registry.FindExact("j", ModeNormal)
	require.NotNil(t, downBinding, "Binding for 'j' not found")

	model.cursor = newCursor(0, 0)
	downBinding.Command(model)

	assert.Equal(t, 1, model.cursor.Row, "j command should increase row by 1")

	// Test move up (k)
	upBinding := model.registry.FindExact("k", ModeNormal)
	require.NotNil(t, upBinding, "Binding for 'k' not found")

	upBinding.Command(model)

	assert.Equal(t, 0, model.cursor.Row, "k command should decrease row by 1")

	// Test move right (l)
	rightBinding := model.registry.FindExact("l", ModeNormal)
	require.NotNil(t, rightBinding, "Binding for 'l' not found")

	rightBinding.Command(model)

	assert.Equal(t, 1, model.cursor.Col, "l command should increase col by 1")

	// Test move left (h)
	leftBinding := model.registry.FindExact("h", ModeNormal)
	require.NotNil(t, leftBinding, "Binding for 'h' not found")

	leftBinding.Command(model)

	assert.Equal(t, 0, model.cursor.Col, "h command should decrease col by 1")
}

func TestWrappedMovementCommands(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2\nLine 3"))

	// Test move to beginning of line (0)
	startBinding := model.registry.FindExact("0", ModeNormal)
	require.NotNil(t, startBinding, "Binding for '0' not found")

	model.cursor = newCursor(1, 3)
	startBinding.Command(model)

	assert.Equal(t, 0, model.cursor.Col, "0 command should set col to 0")

	// Test move to end of line ($)
	endBinding := model.registry.FindExact("$", ModeNormal)
	require.NotNil(t, endBinding, "Binding for '$' not found")

	endBinding.Command(model)

	assert.Equal(t, 5, model.cursor.Col, "$ command should move to end of line")

	// Test space (advance cursor)
	spaceBinding := model.registry.FindExact(" ", ModeNormal)
	require.NotNil(t, spaceBinding, "Binding for space not found")

	// Position cursor at second-to-last position of first line
	model.cursor = newCursor(0, 4)
	spaceBinding.Command(model)

	// Should move to last column
	assert.Equal(t, 5, model.cursor.Col, "Space should move right by 1")

	// One more space should wrap to next line
	spaceBinding.Command(model)

	assert.Equal(t, 1, model.cursor.Row, "Space at end of line should wrap to next line (row)")
	assert.Equal(t, 0, model.cursor.Col, "Space at end of line should wrap to col 0")
}

func TestJumpCommands(t *testing.T) {
	model := New(WithContent("Line 1\nLine 2\nLine 3\nLine 4\nLine 5"))

	// Test move to first line (gg)
	startDocBinding := model.registry.FindExact("gg", ModeNormal)
	require.NotNil(t, startDocBinding, "Binding for 'gg' not found")

	model.cursor = newCursor(3, 0)
	startDocBinding.Command(model)

	assert.Equal(t, 0, model.cursor.Row, "gg command should set row to 0")

	// Test move to last line (G)
	endDocBinding := model.registry.FindExact("G", ModeNormal)
	require.NotNil(t, endDocBinding, "Binding for 'G' not found")

	endDocBinding.Command(model)

	assert.Equal(t, 4, model.cursor.Row, "G command should move to last line (4)")
}

func TestCommandLineCommands(t *testing.T) {
	model := New()

	// Register test command
	cmdExecuted := false
	model.commands.Register("test", func(m *Model) tea.Cmd {
		cmdExecuted = true
		return nil
	})

	// Set up command mode
	model.mode = ModeCommand
	model.commandBuffer = "test"

	// Get the execute command binding
	execBinding := model.registry.FindExact("enter", ModeCommand)
	require.NotNil(t, execBinding, "Binding for 'enter' in command mode not found")

	cmd := execBinding.Command(model)
	model.Update(cmd())

	assert.True(t, cmdExecuted, "Command execution should run registered command")
	assert.Equal(t, ModeNormal, model.mode, "After command execution, mode should be Normal")
}

func TestReplaceCharCommand(t *testing.T) {
	model := New(WithContent("Hello, World!"))

	// Test binding exists
	rBinding := model.registry.FindExact("r", ModeNormal)
	require.NotNil(t, rBinding, "Binding for 'r' not found")
	assert.Equal(t, "r", rBinding.Key, "Expected binding key 'r'")

	// Test that replaceChar sets replacePending
	model.cursor = newCursor(0, 4) // cursor at 'o' in "Hello"
	model.replacePending = false
	rBinding.Command(model)
	assert.True(t, model.replacePending, "replaceChar should set replacePending to true")

	// Next keypress: replace 'o' with 'a'
	model.handleKeypress(tea.KeyPressMsg(tea.Key{Text: "a", Code: 'a'}))
	assert.False(t, model.replacePending, "replacePending should be cleared after replacement")
	assert.Equal(t, "Hella, World!", model.buffer.text(), "Character at cursor should be replaced")
	assert.Equal(t, 4, model.cursor.Col, "Cursor column should remain at replaced position")
	assert.Equal(t, 0, model.cursor.Row, "Cursor row should remain unchanged")
}

func TestReplaceCharEdgeCases(t *testing.T) {
	t.Run("replace at start of line", func(t *testing.T) {
		model := New(WithContent("Hello"))
		rBinding := model.registry.FindExact("r", ModeNormal)
		require.NotNil(t, rBinding)

		model.cursor = newCursor(0, 0)
		rBinding.Command(model)
		model.handleKeypress(tea.KeyPressMsg(tea.Key{Text: "X", Code: 'X'}))

		assert.Equal(t, "Xello", model.buffer.text())
	})

	t.Run("replace at last character", func(t *testing.T) {
		model := New(WithContent("Hello"))
		rBinding := model.registry.FindExact("r", ModeNormal)
		require.NotNil(t, rBinding)

		model.cursor = newCursor(0, 4)
		rBinding.Command(model)
		model.handleKeypress(tea.KeyPressMsg(tea.Key{Text: "!", Code: '!'}))

		assert.Equal(t, "Hell!", model.buffer.text())
	})

	t.Run("replace on empty line does nothing", func(t *testing.T) {
		model := New(WithContent(""))
		rBinding := model.registry.FindExact("r", ModeNormal)
		require.NotNil(t, rBinding)

		model.cursor = newCursor(0, 0)
		rBinding.Command(model)
		assert.True(t, model.replacePending, "replaceChar should set replacePending")

		model.handleKeypress(tea.KeyPressMsg(tea.Key{Text: "a", Code: 'a'}))
		assert.False(t, model.replacePending, "replacePending should be cleared")
		assert.Equal(t, "", model.buffer.text(), "Empty line should remain unchanged")
	})

	t.Run("pressing esc after r cancels replace", func(t *testing.T) {
		model := New(WithContent("Hello"))
		rBinding := model.registry.FindExact("r", ModeNormal)
		require.NotNil(t, rBinding)

		model.cursor = newCursor(0, 2)
		rBinding.Command(model)
		assert.True(t, model.replacePending)

		model.handleKeypress(tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))
		assert.False(t, model.replacePending, "replacePending should be cleared by esc")
		assert.Equal(t, "Hello", model.buffer.text(), "Content should not change")
	})

	t.Run("replace with space", func(t *testing.T) {
		model := New(WithContent("Hello-World"))
		rBinding := model.registry.FindExact("r", ModeNormal)
		require.NotNil(t, rBinding)

		model.cursor = newCursor(0, 5) // cursor at '-'
		rBinding.Command(model)
		model.handleKeypress(tea.KeyPressMsg(tea.Key{Code: tea.KeySpace}))

		assert.Equal(t, "Hello World", model.buffer.text(), "Hyphen should be replaced with space")
	})

	t.Run("chained replace operations", func(t *testing.T) {
		model := New(WithContent("cat"))
		rBinding := model.registry.FindExact("r", ModeNormal)
		require.NotNil(t, rBinding)

		model.cursor = newCursor(0, 0)
		rBinding.Command(model)
		model.handleKeypress(tea.KeyPressMsg(tea.Key{Text: "b", Code: 'b'}))
		assert.Equal(t, "bat", model.buffer.text())

		model.cursor = newCursor(0, 2)
		rBinding.Command(model)
		model.handleKeypress(tea.KeyPressMsg(tea.Key{Text: "g", Code: 'g'}))
		assert.Equal(t, "bag", model.buffer.text())
	})

	t.Run("replace supports undo", func(t *testing.T) {
		model := New(WithContent("Hello"))
		rBinding := model.registry.FindExact("r", ModeNormal)
		require.NotNil(t, rBinding)

		model.cursor = newCursor(0, 0)
		rBinding.Command(model)
		model.handleKeypress(tea.KeyPressMsg(tea.Key{Text: "Y", Code: 'Y'}))

		assert.Equal(t, "Yello", model.buffer.text())

		// Undo
		undoBinding := model.registry.FindExact("u", ModeNormal)
		require.NotNil(t, undoBinding)
		cmd := undoBinding.Command(model)
		if cmd != nil {
			model.Update(cmd())
		}

		assert.Equal(t, "Hello", model.buffer.text(), "Undo should revert the replace")
	})
}
