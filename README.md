# VTea - Vim-like Text Editor for TUIs

> Forked from [**kujtimiihoxha/vimtea**](https://github.com/kujtimiihoxha/vimtea)
---

vtea is a lightweight, Vim-inspired text editor for the terminal, built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework. It provides a modular, extensible foundation for building Vim-like text editors in your terminal applications.

## Features

- Multiple editing modes (Normal, Insert, Visual, Command)
- Vim-like keybindings and commands
- Line numbers (absolute and relative)
- Count-based movement commands (e.g. `5j`, `10k`)
- Undo/redo functionality
- Visual mode selection (character and line-wise)
- Command mode
- Clipboard operations (yank, delete, paste)
- Word operations
- Extensible architecture
- Custom key bindings
- Customizable highlighting

## Installation

```bash
go get github.com/zrcoder/vtea
```

## Code Structure

The codebase has been organized into modular components:

- **model.go**: Main editor model and public interfaces
- **buffer.go**: Text buffer with undo/redo operations
- **cursor.go**: Cursor and text range operations
- **bindings.go**: Key binding registry
- **commands.go**: Command implementations
- **view.go**: Rendering functions
- **highlight.go**: Syntax highlighting
- **styles.go**: UI style definitions

## Usage

See the [example](example/) directory for complete examples:

- [basic](example/basic/main.go) - Basic editor usage
- [load_content](example/load_content/main.go) - Load content into editor
- [custom_bindings](example/custom_bindings/main.go) - Custom key bindings
- [custom_commands](example/custom_commands/main.go) - Custom commands
- [custom_styling](example/custom_styling/main.go) - Custom styling
- [embed_widget](example/embed_widget/main.go) - Embed editor in a larger TUI application

Run any example with:

```bash
go run ./example/<example_name>
```

## Default Key Bindings

### Normal Mode

- `h`, `j`, `k`, `l`: Basic movement (left, down, up, right)
- Number prefixes: `5j`, `10k`: Move multiple lines at once
- `w`: Move to next word start
- `b`: Move to previous word start
- `0`: Move to start of line
- `^`: Move to first non-whitespace character in line
- `$`: Move to end of line
- `gg`: Move to start of document
- `G`: Move to end of document
- `i`: Enter insert mode
- `a`: Append after cursor
- `A`: Append at end of line
- `I`: Insert at start of line
- `v`: Enter visual mode
- `V`: Enter visual line mode
- `:`: Enter command mode
- `x`: Delete character at cursor
- `dd`: Delete line
- `D`: Delete from cursor to end of line
- `yy`: Yank (copy) line
- `p`: Paste after cursor
- `P`: Paste before cursor
- `u`: Undo
- `ctrl+r`: Redo
- `o`: Open line below and enter insert mode
- `O`: Open line above and enter insert mode
- `diw`: Delete inner word
- `yiw`: Yank inner word
- `ciw`: Change inner word
- `zr`: Toggle relative line numbers

### Insert Mode

- `esc`: Return to normal mode
- Arrow keys: Navigate
- Regular typing inserts text

### Visual Mode

- `esc`: Return to normal mode
- `h`, `j`, `k`, `l`: Expand selection
- `y`: Yank selection
- `d`, `x`: Delete selection
- `p`: Replace selection with yanked text

### Command Mode

- `esc`: Cancel command
- `enter`: Execute command

## Extending vtea

vtea is designed to be easily extendable. You can:

1. Add custom key bindings with `editor.AddBinding()`
2. Create new commands with `editor.AddCommand()`
3. Modify the rendering style with custom style options
4. Access buffer operations directly via the Buffer interface
5. Create custom views by implementing the View interface
6. Customize the editor appearance with style options (WithTextStyle, WithLineNumberStyle, WithCurrentLineNumberStyle, etc.)

## Contributing

Contributions are welcome! Here's how you can contribute:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please make sure to update tests as appropriate and follow the existing code style.

## Development Workflow

1. Clone the repository
2. Install dependencies: `go mod download`
3. Make your changes
4. Format code: `go fmt ./...`
5. Verify imports: `goimports -w .`
6. Run the example: `go run example/main.go`
7. Create tests for your changes

## Architecture

vtea follows a modular architecture centered around these core components:

- **Editor**: The main interface that integrates all components
- **Buffer**: Manages text content with undo/redo operations
- **Cursor**: Handles positioning and selection
- **Bindings**: Registers and manages key bindings
- **Commands**: Implements editor commands (like Vim ex commands)
- **View**: Renders the editor to the terminal

These components follow clean separation of concerns, making it easier to:

- Add new features
- Test individual components
- Understand the codebase
- Customize functionality
