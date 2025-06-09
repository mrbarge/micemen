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

	// Show current player
	fmt.Printf("Current Player: %s\n", state.CurrentPlayer.String())

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

	// Print the grid with mice
	for row := 0; row < game.GridHeight; row++ {
		fmt.Print("  ")
		for col := 0; col < game.GridWidth; col++ {
			cell := r.getCellDisplay(state, game.Position{Row: row, Col: col})
			fmt.Print(cell)
		}
		fmt.Println()
	}

	r.showPlayerStats(state)
	r.showControls()
}

// getCellDisplay returns the appropriate emoji for a cell
func (r *TerminalRenderer) getCellDisplay(state game.GameState, pos game.Position) string {
	// Check for mice at this position
	mice := r.getMiceAt(state.Mice, pos)

	isSelected := pos.Col == state.SelectedColumn

	if len(mice) > 0 {
		// Show mice - if multiple mice, show the top one
		// If multiple mice of different colors, show special indicator
		redCount := r.countMiceByColor(mice, game.Red)
		blueCount := r.countMiceByColor(mice, game.Blue)

		if redCount > 0 && blueCount > 0 {
			// Mixed colors
			if isSelected {
				return "🟡" // Highlighted mixed
			}
			return "🟠" // Mixed colors
		} else if redCount > 0 {
			// Red mice
			if isSelected {
				return "🔴" // Highlighted red
			}
			return "🔺" // Red mouse
		} else {
			// Blue mice
			if isSelected {
				return "🔵" // Highlighted blue
			}
			return "🔹" // Blue mouse
		}
	}

	// No mice, show the underlying cell
	switch state.Grid[pos.Row][pos.Col] {
	case game.Wall:
		if isSelected {
			return "🟨" // Highlighted wall
		}
		return "🟫" // Normal wall
	case game.Empty:
		if isSelected {
			return "🔳" // Highlighted empty space
		}
		return "⬛" // Normal empty space
	default:
		return "❓" // Unknown
	}
}

// getMiceAt returns all mice at the given position
func (r *TerminalRenderer) getMiceAt(mice []game.Mouse, pos game.Position) []game.Mouse {
	var result []game.Mouse
	for _, mouse := range mice {
		if mouse.Position.Row == pos.Row && mouse.Position.Col == pos.Col {
			result = append(result, mouse)
		}
	}
	return result
}

// countMiceByColor counts mice of a specific color
func (r *TerminalRenderer) countMiceByColor(mice []game.Mouse, color game.PlayerColor) int {
	count := 0
	for _, mouse := range mice {
		if mouse.Player == color {
			count++
		}
	}
	return count
}

// showPlayerStats displays information about each player's mice
func (r *TerminalRenderer) showPlayerStats(state game.GameState) {
	redPlayer := r.getPlayerInfo(state.Mice, game.Red)
	bluePlayer := r.getPlayerInfo(state.Mice, game.Blue)

	fmt.Printf("\nPlayer Stats:\n")
	fmt.Printf("🔺 Red:  %d mice\n", len(redPlayer))
	fmt.Printf("🔹 Blue: %d mice\n", len(bluePlayer))
}

// getPlayerInfo returns mice for a specific player
func (r *TerminalRenderer) getPlayerInfo(mice []game.Mouse, color game.PlayerColor) []game.Mouse {
	var result []game.Mouse
	for _, mouse := range mice {
		if mouse.Player == color {
			result = append(result, mouse)
		}
	}
	return result
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
	fmt.Println("\nLegend:")
	fmt.Println("🔺 Red mice    🔹 Blue mice    🟠 Mixed")
	fmt.Println("🟫 Wall        ⬛ Empty space")
}

// HideCursor hides the terminal cursor
func (r *TerminalRenderer) HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func (r *TerminalRenderer) ShowCursor() {
	fmt.Print("\033[?25h")
}
