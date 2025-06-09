package game

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame()
	state := game.GetState()

	// Test initial state
	if state.SelectedColumn != GridWidth/2 {
		t.Errorf("Initial selected column should be %d, got %d", GridWidth/2, state.SelectedColumn)
	}

	if state.GameOver {
		t.Error("New game should not be over")
	}

	// Verify grid dimensions are correct
	if len(state.Grid) != GridHeight {
		t.Errorf("Grid height should be %d, got %d", GridHeight, len(state.Grid))
	}
	if len(state.Grid[0]) != GridWidth {
		t.Errorf("Grid width should be %d, got %d", GridWidth, len(state.Grid[0]))
	}
}

func TestGenerateWalls(t *testing.T) {
	game := NewGame()
	state := game.GetState()

	// Test that each column has the correct number of walls
	for col := 0; col < GridWidth; col++ {
		wallCount := 0
		for row := 0; row < GridHeight; row++ {
			if state.Grid[row][col] == Wall {
				wallCount++
			}
		}

		if wallCount < MinWalls || wallCount > MaxWalls {
			t.Errorf("Column %d has %d walls, expected between %d and %d",
				col, wallCount, MinWalls, MaxWalls)
		}
	}
}

func TestMoveColumnUp(t *testing.T) {
	game := NewGame()

	// Set up a known pattern in column 0
	game.state.SelectedColumn = 0
	game.state.Grid[0][0] = Wall
	game.state.Grid[1][0] = Empty
	game.state.Grid[2][0] = Wall

	// Store original pattern for comparison
	originalPattern := make([]CellType, GridHeight)
	for i := 0; i < GridHeight; i++ {
		originalPattern[i] = game.state.Grid[i][0]
	}

	game.ProcessAction(ActionMoveColumnUp)
	state := game.GetState()

	// Check that each cell moved up by one position (with wraparound)
	for row := 0; row < GridHeight-1; row++ {
		expected := originalPattern[row+1]
		actual := state.Grid[row][0]
		if actual != expected {
			t.Errorf("After moving up, row %d should be %v, got %v", row, expected, actual)
		}
	}

	// Check wraparound: bottom should be what was originally at top
	if state.Grid[GridHeight-1][0] != originalPattern[0] {
		t.Errorf("After moving up, bottom cell should be %v, got %v",
			originalPattern[0], state.Grid[GridHeight-1][0])
	}
}

func TestMoveColumnDown(t *testing.T) {
	game := NewGame()

	// Set up a known pattern in column 0
	game.state.SelectedColumn = 0
	game.state.Grid[0][0] = Wall
	game.state.Grid[1][0] = Empty
	game.state.Grid[2][0] = Wall

	// Store original pattern for comparison
	originalPattern := make([]CellType, GridHeight)
	for i := 0; i < GridHeight; i++ {
		originalPattern[i] = game.state.Grid[i][0]
	}

	game.ProcessAction(ActionMoveColumnDown)
	state := game.GetState()

	// Check that each cell moved down by one position (with wraparound)
	for row := 1; row < GridHeight; row++ {
		expected := originalPattern[row-1]
		actual := state.Grid[row][0]
		if actual != expected {
			t.Errorf("After moving down, row %d should be %v, got %v", row, expected, actual)
		}
	}

	// Check wraparound: top should be what was originally at bottom
	if state.Grid[0][0] != originalPattern[GridHeight-1] {
		t.Errorf("After moving down, top cell should be %v, got %v",
			originalPattern[GridHeight-1], state.Grid[0][0])
	}
}

func TestMoveSelection(t *testing.T) {
	game := NewGame()
	game.state.SelectedColumn = 5

	// Test moving right
	game.ProcessAction(ActionMoveRight)
	if game.GetState().SelectedColumn != 6 {
		t.Errorf("After moving right, selected column should be 6, got %d", game.GetState().SelectedColumn)
	}

	// Test moving left
	game.ProcessAction(ActionMoveLeft)
	if game.GetState().SelectedColumn != 5 {
		t.Errorf("After moving left, selected column should be 5, got %d", game.GetState().SelectedColumn)
	}
}

func TestMoveSelectionBounds(t *testing.T) {
	// Test left boundary
	game := NewGame()
	game.state.SelectedColumn = 0
	game.ProcessAction(ActionMoveLeft)
	if game.GetState().SelectedColumn != 0 {
		t.Errorf("Selection should not move beyond left boundary, got %d", game.GetState().SelectedColumn)
	}

	// Test right boundary
	game.state.SelectedColumn = GridWidth - 1
	game.ProcessAction(ActionMoveRight)
	if game.GetState().SelectedColumn != GridWidth-1 {
		t.Errorf("Selection should not move beyond right boundary, got %d", game.GetState().SelectedColumn)
	}
}

func TestQuitAction(t *testing.T) {
	game := NewGame()
	if game.IsGameOver() {
		t.Error("Game should not be over initially")
	}

	game.ProcessAction(ActionQuit)
	if !game.IsGameOver() {
		t.Error("Game should be over after quit action")
	}
}

func TestActionOnGameOver(t *testing.T) {
	game := NewGame()
	game.ProcessAction(ActionQuit) // End the game

	originalState := game.GetState()

	// Try to perform actions after game over
	game.ProcessAction(ActionMoveRight)
	game.ProcessAction(ActionMoveColumnUp)

	newState := game.GetState()

	// State should be unchanged
	if newState.SelectedColumn != originalState.SelectedColumn {
		t.Error("Game state should not change after game over")
	}
}

func TestReset(t *testing.T) {
	game := NewGame()

	// Make some changes
	game.ProcessAction(ActionMoveRight)
	game.ProcessAction(ActionQuit)

	// Reset the game
	game.Reset()
	state := game.GetState()

	if state.GameOver {
		t.Error("Game should not be over after reset")
	}

	if state.SelectedColumn != GridWidth/2 {
		t.Errorf("Selected column should be reset to %d, got %d", GridWidth/2, state.SelectedColumn)
	}
}

func TestWallCountPreservation(t *testing.T) {
	game := NewGame()
	state := game.GetState()

	// Count walls in selected column before moving
	originalWallCount := 0
	for row := 0; row < GridHeight; row++ {
		if state.Grid[row][state.SelectedColumn] == Wall {
			originalWallCount++
		}
	}

	// Move column and recount
	game.ProcessAction(ActionMoveColumnUp)
	newState := game.GetState()
	newWallCount := 0
	for row := 0; row < GridHeight; row++ {
		if newState.Grid[row][newState.SelectedColumn] == Wall {
			newWallCount++
		}
	}

	if newWallCount != originalWallCount {
		t.Errorf("Wall count should be preserved after move: expected %d, got %d",
			originalWallCount, newWallCount)
	}
}
