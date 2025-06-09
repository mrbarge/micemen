package game

import (
	"math/rand"
	"time"
)

// MicemenGame implements the Game interface
type MicemenGame struct {
	state GameState
}

// NewGame creates a new game instance
func NewGame() *MicemenGame {
	game := &MicemenGame{}
	game.Reset()
	return game
}

// Reset initializes a new game
func (g *MicemenGame) Reset() {
	g.state = GameState{
		SelectedColumn: GridWidth / 2,
		GameOver:       false,
		CurrentPlayer:  Red, // Red player starts
		Mice:           make([]Mouse, 0, MicePerPlayer*2),
	}
	g.generateWalls()
	g.placeMice()
}

// GetState returns a copy of the current game state
func (g *MicemenGame) GetState() GameState {
	return g.state
}

// IsGameOver returns whether the game has ended
func (g *MicemenGame) IsGameOver() bool {
	return g.state.GameOver
}

// GetPlayer returns player information for the given color
func (g *MicemenGame) GetPlayer(color PlayerColor) Player {
	var mice []Mouse
	for _, mouse := range g.state.Mice {
		if mouse.Player == color {
			mice = append(mice, mouse)
		}
	}
	return Player{
		Color: color,
		Mice:  mice,
	}
}

// GetMiceAt returns all mice at the given position
func (g *MicemenGame) GetMiceAt(pos Position) []Mouse {
	var mice []Mouse
	for _, mouse := range g.state.Mice {
		if mouse.Position.Row == pos.Row && mouse.Position.Col == pos.Col {
			mice = append(mice, mouse)
		}
	}
	return mice
}

// ProcessAction handles a player action
func (g *MicemenGame) ProcessAction(action Action) {
	if g.state.GameOver {
		return
	}

	switch action {
	case ActionMoveLeft:
		g.moveSelection(-1)
	case ActionMoveRight:
		g.moveSelection(1)
	case ActionMoveColumnUp:
		g.moveColumnUp()
	case ActionMoveColumnDown:
		g.moveColumnDown()
	case ActionQuit:
		g.state.GameOver = true
	}
}

// generateWalls randomly places walls in each column
func (g *MicemenGame) generateWalls() {
	rand.Seed(time.Now().UnixNano())

	for col := 0; col < GridWidth; col++ {
		// Random number of walls for this column
		numWalls := rand.Intn(MaxWalls-MinWalls+1) + MinWalls

		// Generate random positions for walls
		positions := make(map[int]bool)
		for len(positions) < numWalls {
			pos := rand.Intn(GridHeight)
			positions[pos] = true
		}

		// Place walls at selected positions
		for row := 0; row < GridHeight; row++ {
			if positions[row] {
				g.state.Grid[row][col] = Wall
			} else {
				g.state.Grid[row][col] = Empty
			}
		}
	}
}

// placeMice randomly places mice for both players
func (g *MicemenGame) placeMice() {
	rand.Seed(time.Now().UnixNano())

	// Place Red player's mice (left 9 columns)
	g.placeMiceForPlayer(Red, 0, Player1Columns-1)

	// Place Blue player's mice (right 9 columns)
	g.placeMiceForPlayer(Blue, GridWidth-Player2Columns, GridWidth-1)
}

// placeMiceForPlayer places mice for a specific player in the given column range
func (g *MicemenGame) placeMiceForPlayer(player PlayerColor, startCol, endCol int) {
	for i := 0; i < MicePerPlayer; i++ {
		// Find a valid position for this mouse
		pos := g.findValidMousePosition(startCol, endCol)
		if pos != nil {
			mouse := Mouse{
				Position: *pos,
				Player:   player,
			}
			g.state.Mice = append(g.state.Mice, mouse)
		}
	}
}

// findValidMousePosition finds a valid position for a mouse in the given column range
func (g *MicemenGame) findValidMousePosition(startCol, endCol int) *Position {
	maxAttempts := 1000 // Prevent infinite loops
	attempts := 0

	for attempts < maxAttempts {
		attempts++

		// Random column in range
		col := startCol + rand.Intn(endCol-startCol+1)

		// Find valid rows in this column (must be above a wall or another mouse)
		validRows := g.getValidRowsForMouse(col)
		if len(validRows) == 0 {
			continue
		}

		// Pick a random valid row
		row := validRows[rand.Intn(len(validRows))]
		return &Position{Row: row, Col: col}
	}

	return nil // Could not find valid position
}

// getValidRowsForMouse returns all valid rows where a mouse can be placed in the given column
func (g *MicemenGame) getValidRowsForMouse(col int) []int {
	var validRows []int

	for row := 0; row < GridHeight; row++ {
		if g.isValidMousePosition(Position{Row: row, Col: col}) {
			validRows = append(validRows, row)
		}
	}

	return validRows
}

// isValidMousePosition checks if a mouse can be placed at the given position
func (g *MicemenGame) isValidMousePosition(pos Position) bool {
	// Check if position is within bounds
	if pos.Row < 0 || pos.Row >= GridHeight || pos.Col < 0 || pos.Col >= GridWidth {
		return false
	}

	// Mouse must be placed directly above a wall or another mouse
	if pos.Row == GridHeight-1 {
		// Bottom row: must be above a wall
		return g.state.Grid[pos.Row][pos.Col] == Wall
	}

	// Check if there's support below (wall or mouse)
	belowPos := Position{Row: pos.Row + 1, Col: pos.Col}

	// Check for wall below
	if g.state.Grid[belowPos.Row][belowPos.Col] == Wall {
		return true
	}

	// Check for mouse below
	miceBelow := g.GetMiceAt(belowPos)
	return len(miceBelow) > 0
}

// moveColumnUp shifts all cells in the selected column up (with wraparound) and updates mouse positions
func (g *MicemenGame) moveColumnUp() {
	if g.state.SelectedColumn < 0 || g.state.SelectedColumn >= GridWidth {
		return
	}

	col := g.state.SelectedColumn

	// Store the top cell
	topCell := g.state.Grid[0][col]

	// Shift all cells up
	for row := 0; row < GridHeight-1; row++ {
		g.state.Grid[row][col] = g.state.Grid[row+1][col]
	}

	// Wrap the top cell to the bottom
	g.state.Grid[GridHeight-1][col] = topCell

	// Update mouse positions in this column
	g.updateMiceForColumnShift(col, true)
}

// moveColumnDown shifts all cells in the selected column down (with wraparound) and updates mouse positions
func (g *MicemenGame) moveColumnDown() {
	if g.state.SelectedColumn < 0 || g.state.SelectedColumn >= GridWidth {
		return
	}

	col := g.state.SelectedColumn

	// Store the bottom cell
	bottomCell := g.state.Grid[GridHeight-1][col]

	// Shift all cells down
	for row := GridHeight - 1; row > 0; row-- {
		g.state.Grid[row][col] = g.state.Grid[row-1][col]
	}

	// Wrap the bottom cell to the top
	g.state.Grid[0][col] = bottomCell

	// Update mouse positions in this column
	g.updateMiceForColumnShift(col, false)
}

// updateMiceForColumnShift updates mouse positions when a column is shifted
func (g *MicemenGame) updateMiceForColumnShift(col int, shiftUp bool) {
	for i := range g.state.Mice {
		mouse := &g.state.Mice[i]
		if mouse.Position.Col == col {
			if shiftUp {
				// Shift up: row decreases, with wraparound
				if mouse.Position.Row == 0 {
					mouse.Position.Row = GridHeight - 1
				} else {
					mouse.Position.Row--
				}
			} else {
				// Shift down: row increases, with wraparound
				if mouse.Position.Row == GridHeight-1 {
					mouse.Position.Row = 0
				} else {
					mouse.Position.Row++
				}
			}
		}
	}
}

// moveSelection moves the column selection left or right
func (g *MicemenGame) moveSelection(direction int) {
	newCol := g.state.SelectedColumn + direction
	if newCol >= 0 && newCol < GridWidth {
		g.state.SelectedColumn = newCol
	}
}
