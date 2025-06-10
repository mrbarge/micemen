package render

import (
	"fmt"
	"micemen/game"
)

// TerminalRenderer implements the Renderer interface for terminal output
type TerminalRenderer struct {
	game game.Game // Reference to game for querying state
}

// NewTerminalRenderer creates a new terminal renderer
func NewTerminalRenderer(g game.Game) *TerminalRenderer {
	return &TerminalRenderer{game: g}
}

// Clear clears the terminal screen
func (r *TerminalRenderer) Clear() {
	fmt.Print("\033[2J\033[H")
}

// Render displays the current game state
func (r *TerminalRenderer) Render(state game.GameState) {
	r.Clear()

	// Show current player with emphasis
	playerIcon := "🔺"
	if state.CurrentPlayer == game.Blue {
		playerIcon = "🔹"
	}
	fmt.Printf("%s %s Player's Turn %s\n", playerIcon, state.CurrentPlayer.String(), playerIcon)

	// Print column indicators with validity markers
	fmt.Print("  ")
	for col := 0; col < game.GridWidth; col++ {
		if col == state.SelectedColumn {
			if r.game.CanPlayerMoveColumn(state.CurrentPlayer, col) {
				fmt.Print("🔽") // Valid selected column
			} else {
				fmt.Print("❌") // Invalid selected column
			}
		} else {
			if r.game.CanPlayerMoveColumn(state.CurrentPlayer, col) {
				fmt.Print("✓ ") // Valid column
			} else {
				fmt.Print("  ") // Invalid/empty column
			}
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
	r.showTurnInfo(state)
	r.showControls()
}

// getCellDisplay returns the appropriate emoji for a cell
func (r *TerminalRenderer) getCellDisplay(state game.GameState, pos game.Position) string {
	// Check for mice at this position
	mice := r.getMiceAt(state.Mice, pos)

	isSelected := pos.Col == state.SelectedColumn
	isValidColumn := r.game.CanPlayerMoveColumn(state.CurrentPlayer, pos.Col)

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
				if isValidColumn {
					return "🔴" // Highlighted red (valid)
				} else {
					return "🟤" // Highlighted red (invalid)
				}
			}
			return "🔺" // Red mouse
		} else {
			// Blue mice
			if isSelected {
				if isValidColumn {
					return "🔵" // Highlighted blue (valid)
				} else {
					return "🟦" // Highlighted blue (invalid)
				}
			}
			return "🔹" // Blue mouse
		}
	}

	// No mice, show the underlying cell
	switch state.Grid[pos.Row][pos.Col] {
	case game.Wall:
		if isSelected {
			if isValidColumn {
				return "🟨" // Highlighted wall (valid column)
			} else {
				return "🟫" // Highlighted wall (invalid column)
			}
		}
		return "🟫" // Normal wall
	case game.Empty:
		if isSelected {
			if isValidColumn {
				return "🔳" // Highlighted empty space (valid column)
			} else {
				return "⬜" // Highlighted empty space (invalid column)
			}
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
	fmt.Printf("🔺 Red:  %d mice | Valid columns: %s\n",
		len(redPlayer), r.getValidColumnsDisplay(game.Red))
	fmt.Printf("🔹 Blue: %d mice | Valid columns: %s\n",
		len(bluePlayer), r.getValidColumnsDisplay(game.Blue))
}

// getValidColumnsDisplay returns a display string for valid columns
func (r *TerminalRenderer) getValidColumnsDisplay(player game.PlayerColor) string {
	validCols := r.game.GetValidColumnsForPlayer(player)
	if len(validCols) == 0 {
		return "None"
	}

	result := ""
	for i, col := range validCols {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%d", col+1) // 1-based for display
	}
	return result
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

// showTurnInfo displays turn-specific information
func (r *TerminalRenderer) showTurnInfo(state game.GameState) {
	fmt.Printf("\nTurn Info:\n")

	// Check if current selection is valid
	isValidSelection := r.game.CanPlayerMoveColumn(state.CurrentPlayer, state.SelectedColumn)
	if isValidSelection {
		fmt.Printf("✅ Column %d is ready to move!\n", state.SelectedColumn+1)
		fmt.Println("   Use ↑/↓ (or W/S or K/J) to move this column")
	} else {
		fmt.Printf("❌ Column %d has no %s mice\n", state.SelectedColumn+1, state.CurrentPlayer.String())
		fmt.Println("   Use ←/→ (or A/D or H/L) to find a valid column")
	}
}

// ShowMessage displays a message to the user
func (r *TerminalRenderer) ShowMessage(msg string) {
	fmt.Println(msg)
}

// showControls displays the control instructions
func (r *TerminalRenderer) showControls() {
	fmt.Println("\nControls:")
	fmt.Println("← → (or A/D or H/L) : Select column with your mice")
	fmt.Println("↑ ↓ (or W/S or K/J)  : Move your column up/down")
	fmt.Println("q                    : Quit")
	fmt.Println("\nLegend:")
	fmt.Println("🔺 Red mice    🔹 Blue mice    🟠 Mixed")
	fmt.Println("🟫 Wall        ⬛ Empty        ✓ Valid column")
}

// HideCursor hides the terminal cursor
func (r *TerminalRenderer) HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func (r *TerminalRenderer) ShowCursor() {
	fmt.Print("\033[?25h")
}
