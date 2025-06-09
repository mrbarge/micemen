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
	}
	g.generateWalls()
}

// GetState returns a copy of the current game state
func (g *MicemenGame) GetState() GameState {
	return g.state
}

// IsGameOver returns whether the game has ended
func (g *MicemenGame) IsGameOver() bool {
	return g.state.GameOver
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

// moveColumnUp shifts all cells in the selected column up (with wraparound)
func (g *MicemenGame) moveColumnUp() {
	if g.state.SelectedColumn < 0 || g.state.SelectedColumn >= GridWidth {
		return
	}

	// Store the top cell
	topCell := g.state.Grid[0][g.state.SelectedColumn]

	// Shift all cells up
	for row := 0; row < GridHeight-1; row++ {
		g.state.Grid[row][g.state.SelectedColumn] = g.state.Grid[row+1][g.state.SelectedColumn]
	}

	// Wrap the top cell to the bottom
	g.state.Grid[GridHeight-1][g.state.SelectedColumn] = topCell
}

// moveColumnDown shifts all cells in the selected column down (with wraparound)
func (g *MicemenGame) moveColumnDown() {
	if g.state.SelectedColumn < 0 || g.state.SelectedColumn >= GridWidth {
		return
	}

	// Store the bottom cell
	bottomCell := g.state.Grid[GridHeight-1][g.state.SelectedColumn]

	// Shift all cells down
	for row := GridHeight - 1; row > 0; row-- {
		g.state.Grid[row][g.state.SelectedColumn] = g.state.Grid[row-1][g.state.SelectedColumn]
	}

	// Wrap the bottom cell to the top
	g.state.Grid[0][g.state.SelectedColumn] = bottomCell
}

// moveSelection moves the column selection left or right
func (g *MicemenGame) moveSelection(direction int) {
	newCol := g.state.SelectedColumn + direction
	if newCol >= 0 && newCol < GridWidth {
		g.state.SelectedColumn = newCol
	}
}
