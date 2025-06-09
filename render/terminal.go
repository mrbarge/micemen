package render

import (
	"fmt"
	"micemen/game"
)

// TerminalRenderer implements the Renderer interface for terminal output
type TerminalRenderer struct{}

// NewTerminalRenderer creates a new terminal renderer
func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{}
}

// Clear clears the terminal screen
func (r *TerminalRenderer) Clear() {
	fmt.Print("\033[2J\033[H")
}

// Render displays the current game state
func (r *TerminalRenderer) Render(state game.GameState) {
	r.Clear()

	// Print column indicators
	fmt.Print("  ")
	for col := 0; col < game.GridWidth; col++ {
		if col == state.SelectedColumn {
			fmt.Print("🔽")
		} else {
			fmt.Print("  ")
		}
	}
	fmt.Println()

	// Print the grid
	for row := 0; row < game.GridHeight; row++ {
		fmt.Print("  ")
		for col := 0; col < game.GridWidth; col++ {
			switch state.Grid[row][col] {
			case game.Wall:
				if col == state.SelectedColumn {
					fmt.Print("🟨") // Highlighted wall
				} else {
					fmt.Print("🟫") // Normal wall
				}
			case game.Empty:
				if col == state.SelectedColumn {
					fmt.Print("🔳") // Highlighted empty space
				} else {
					fmt.Print("⬛") // Normal empty space
				}
			}
		}
		fmt.Println()
	}

	r.showControls()
}

// ShowMessage displays a message to the user
func (r *TerminalRenderer) ShowMessage(msg string) {
	fmt.Println(msg)
}

// showControls displays the control instructions
func (r *TerminalRenderer) showControls() {
	fmt.Println("\nControls:")
	fmt.Println("← → (or A/D or H/L) : Select column")
	fmt.Println("↑ ↓ (or W/S or K/J)  : Move column up/down")
	fmt.Println("q                    : Quit")
}

// HideCursor hides the terminal cursor
func (r *TerminalRenderer) HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func (r *TerminalRenderer) ShowCursor() {
	fmt.Print("\033[?25h")
}
