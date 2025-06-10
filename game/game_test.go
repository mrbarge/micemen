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
	
	if state.CurrentPlayer != Red {
		t.Errorf("Initial player should be Red, got %v", state.CurrentPlayer)
	}
	
	// Verify grid dimensions are correct
	if len(state.Grid) != GridHeight {
		t.Errorf("Grid height should be %d, got %d", GridHeight, len(state.Grid))
	}
	if len(state.Grid[0]) != GridWidth {
		t.Errorf("Grid width should be %d, got %d", GridWidth, len(state.Grid[0]))
	}
}

func TestPlayerColors(t *testing.T) {
	if Red.String() != "Red" {
		t.Errorf("Red.String() should return 'Red', got '%s'", Red.String())
	}
	if Blue.String() != "Blue" {
		t.Errorf("Blue.String() should return 'Blue', got '%s'", Blue.String())
	}
}

func TestTurnBasedMovement(t *testing.T) {
	game := NewGame()
	
	// Test that only valid columns can be moved
	originalPlayer := game.GetState().CurrentPlayer
	
	// Find a column that the current player can move
	validColumns := game.getValidColumnsForPlayer(originalPlayer)
	if len(validColumns) == 0 {
		t.Fatal("No valid columns found for current player")
	}
	
	// Set selection to a valid column and try to move
	game.state.SelectedColumn = validColumns[0]
	originalGrid := game.state.Grid
	
	game.ProcessAction(ActionMoveColumnUp)
	newState := game.GetState()
	
	// Player should have switched
	if newState.CurrentPlayer == originalPlayer {
		t.Error("Player should have switched after valid move")
	}
	
	// Grid should have changed
	if game.state.Grid == originalGrid {
		t.Error("Grid should have changed after move")
	}
}

func TestInvalidColumnMovement(t *testing.T) {
	game := NewGame()
	
	// Find a column that the current player cannot move
	currentPlayer := game.GetState().CurrentPlayer
	validColumns := game.getValidColumnsForPlayer(currentPlayer)
	
	// Find an invalid column
	invalidCol := -1
	for col := 0; col < GridWidth; col++ {
		isValid := false
		for _, validCol := range validColumns {
			if col == validCol {
				isValid = true
				break
			}
		}
		if !isValid {
			invalidCol = col
			break
		}
	}
	
	if invalidCol == -1 {
		t.Skip("Could not find invalid column for test")
	}
	
	// Try to move invalid column
	game.state.SelectedColumn = invalidCol
	originalPlayer := game.GetState().CurrentPlayer
	originalGrid := game.state.Grid
	
	game.ProcessAction(ActionMoveColumnUp)
	newState := game.GetState()
	
	// Player should NOT have switched
	if newState.CurrentPlayer != originalPlayer {
		t.Error("Player should not switch after invalid move")
	}
	
	// Grid should NOT have changed
	if game.state.Grid != originalGrid {
		t.Error("Grid should not change after invalid move")
	}
}

func TestCanPlayerMoveColumn(t *testing.T) {
	game := NewGame()
	
	// Create a test scenario with known mice positions
	game.state.Mice = []Mouse{
		{Position: Position{Row: 5, Col: 3}, Player: Red},
		{Position: Position{Row: 7, Col: 10}, Player: Blue},
	}
	
	// Red player should be able to move column 3
	if !game.canPlayerMoveColumn(Red, 3) {
		t.Error("Red player should be able to move column 3")
	}
	
	// Red player should NOT be able to move column 10
	if game.canPlayerMoveColumn(Red, 10) {
		t.Error("Red player should not be able to move column 10")
	}
	
	// Blue player should be able to move column 10
	if !game.canPlayerMoveColumn(Blue, 10) {
		t.Error("Blue player should be able to move column 10")
	}
	
	// Blue player should NOT be able to move column 3
	if game.canPlayerMoveColumn(Blue, 3) {
		t.Error("Blue player should not be able to move column 3")
	}
	
	// Neither player should be able to move column with no mice
	if game.canPlayerMoveColumn(Red, 15) {
		t.Error("Red player should not be able to move empty column")
	}
	if game.canPlayerMoveColumn(Blue, 15) {
		t.Error("Blue player should not be able to move empty column")
	}
}

func TestGetValidColumnsForPlayer(t *testing.T) {
	game := NewGame()
	
	// Create test scenario
	game.state.Mice = []Mouse{
		{Position: Position{Row: 1, Col: 2}, Player: Red},
		{Position: Position{Row: 3, Col: 2}, Player: Red},  // Same column as above
		{Position: Position{Row: 5, Col: 5}, Player: Red},
		{Position: Position{Row: 7, Col: 10}, Player: Blue},
		{Position: Position{Row: 9, Col: 15}, Player: Blue},
	}
	
	redCols := game.getValidColumnsForPlayer(Red)
	blueCols := game.getValidColumnsForPlayer(Blue)
	
	// Red should have columns 2 and 5
	expectedRed := []int{2, 5}
	if len(redCols) != len(expectedRed) {
		t.Errorf("Red should have %d valid columns, got %d", len(expectedRed), len(redCols))
	}
	for i, col := range expectedRed {
		if redCols[i] != col {
			t.Errorf("Red column %d should be %d, got %d", i, col, redCols[i])
		}
	}
	
	// Blue should have columns 10 and 15
	expectedBlue := []int{10, 15}
	if len(blueCols) != len(expectedBlue) {
		t.Errorf("Blue should have %d valid columns, got %d", len(expectedBlue), len(blueCols))
	}
	for i, col := range expectedBlue {
		if blueCols[i] != col {
			t.Errorf("Blue column %d should be %d, got %d", i, col, blueCols[i])
		}
	}
}

func TestMoveSelectionToValidColumn(t *testing.T) {
	game := NewGame()
	
	// Set up test scenario
	game.state.Mice = []Mouse{
		{Position: Position{Row: 1, Col: 3}, Player: Red},
		{Position: Position{Row: 1, Col: 7}, Player: Red},
		{Position: Position{Row: 1, Col: 12}, Player: Red},
	}
	game.state.CurrentPlayer = Red
	game.state.SelectedColumn = 5 // Between valid columns 3 and 7
	
	// Move right should go to column 7
	game.moveSelectionToValidColumn(1)
	if game.state.SelectedColumn != 7 {
		t.Errorf("Moving right from 5 should select 7, got %d", game.state.SelectedColumn)
	}
	
	// Move right again should wrap to column 3
	game.moveSelectionToValidColumn(1)
	if game.state.SelectedColumn != 12 {
		t.Errorf("Moving right from 7 should select 12, got %d", game.state.SelectedColumn)
	}
	
	// Move right again should wrap to first column
	game.moveSelectionToValidColumn(1)
	if game.state.SelectedColumn != 3 {
		t.Errorf("Moving right from 12 should wrap to 3, got %d", game.state.SelectedColumn)
	}
	
	// Move left should go to column 12
	game.moveSelectionToValidColumn(-1)
	if game.state.SelectedColumn != 12 {
		t.Errorf("Moving left from 3 should wrap to 12, got %d", game.state.SelectedColumn)
	}
}

func TestSwitchPlayer(t *testing.T) {
	game := NewGame()
	
	// Set up known mice positions for both players
	game.state.Mice = []Mouse{
		{Position: Position{Row: 1, Col: 3}, Player: Red},
		{Position: Position{Row: 1, Col: 15}, Player: Blue},
	}
	
	// Start with Red player
	game.state.CurrentPlayer = Red
	game.state.SelectedColumn = 3
	
	// Switch to Blue
	game.switchPlayer()
	
	if game.state.CurrentPlayer != Blue {
		t.Errorf("Player should switch to Blue, got %v", game.state.CurrentPlayer)
	}
	
	// Selection should move to a valid column for Blue (15)
	if game.state.SelectedColumn != 15 {
		t.Errorf("Selection should move to Blue's valid column 15, got %d", game.state.SelectedColumn)
	}
	
	// Switch back to Red
	game.switchPlayer()
	
	if game.state.CurrentPlayer != Red {
		t.Errorf("Player should switch back to Red, got %v", game.state.CurrentPlayer)
	}
	
	// Selection should move to a valid column for Red (3)
	if game.state.SelectedColumn != 3 {
		t.Errorf("Selection should move to Red's valid column 3, got %d", game.state.SelectedColumn)
	}
}

func TestInitialValidColumnSelection(t *testing.T) {
	game := NewGame()
	state := game.GetState()
	
	// Initial selection should be on a valid column for the starting player
	if !game.canPlayerMoveColumn(state.CurrentPlayer, state.SelectedColumn) {
		t.Error("Initial column selection should be valid for starting player")
	}
}

// Original tests continue...

func TestMicePlacement(t *testing.T) {
	game := NewGame()
	state := game.GetState()
	
	// Count mice for each player
	redMice := 0
	blueMice := 0
	
	for _, mouse := range state.Mice {
		switch mouse.Player {
		case Red:
			redMice++
			// Red mice should be in left 9 columns
			if mouse.Position.Col >= Player1Columns {
				t.Errorf("Red mouse at column %d should be in columns 0-%d", 
					mouse.Position.Col, Player1Columns-1)
			}
		case Blue:
			blueMice++
			// Blue mice should be in right 9 columns
			if mouse.Position.Col < GridWidth-Player2Columns {
				t.Errorf("Blue mouse at column %d should be in columns %d-%d", 
					mouse.Position.Col, GridWidth-Player2Columns, GridWidth-1)
			}
		}
	}
	
	// Should have correct number of mice per player
	if redMice != MicePerPlayer {
		t.Errorf("Should have %d red mice, got %d", MicePerPlayer, redMice)
	}
	if blueMice != MicePerPlayer {
		t.Errorf("Should have %d blue mice, got %d", MicePerPlayer, blueMice)
	}
}

func TestMicePositionValidity(t *testing.T) {
	game := NewGame()
	state := game.GetState()
	
	for _, mouse := range state.Mice {
		pos := mouse.Position
		
		// Check bounds
		if pos.Row < 0 || pos.Row >= GridHeight || pos.Col < 0 || pos.Col >= GridWidth {
			t.Errorf("Mouse at invalid position: row=%d, col=%d", pos.Row, pos.Col)
		}
		
		// Check if mouse has proper support
		if !game.isValidMousePosition(pos) {
			t.Errorf("Mouse at row=%d, col=%d lacks proper support", pos.Row, pos.Col)
		}
	}
}

func TestGetPlayer(t *testing.T) {
	game := NewGame()
	
	redPlayer := game.GetPlayer(Red)
	bluePlayer := game.GetPlayer(Blue)
	
	if redPlayer.Color != Red {
		t.Errorf("Red player should have Red color, got %v", redPlayer.Color)
	}
	if bluePlayer.Color != Blue {
		t.Errorf("Blue player should have Blue color, got %v", bluePlayer.Color)
	}
	
	if len(redPlayer.Mice) != MicePerPlayer {
		t.Errorf("Red player should have %d mice, got %d", MicePerPlayer, len(redPlayer.Mice))
	}
	if len(bluePlayer.Mice) != MicePerPlayer {
		t.Errorf("Blue player should have %d mice, got %d", MicePerPlayer, len(bluePlayer.Mice))
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
	
	if state.CurrentPlayer != Red {
		t.Errorf("Current player should be reset to Red, got %v", state.CurrentPlayer)
	}
	
	// Should have proper number of mice again
	if len(state.Mice) != MicePerPlayer*2 {
		t.Errorf("Should have %d total mice after reset, got %d", MicePerPlayer*2, len(state.Mice))
	}
	
	// Initial selection should be valid for starting player
	if !game.canPlayerMoveColumn(state.CurrentPlayer, state.SelectedColumn) {
		t.Error("Initial selection after reset should be valid for starting player")
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
} %d", MicePerPlayer*2, len(state.Mice))
	}
} %d", MicePerPlayer*2, len(state.Mice))
	}
}

func TestIsValidMousePosition(t *testing.T) {
	game := NewGame()
	
	// Test bounds checking
	if game.isValidMousePosition(Position{Row: -1, Col: 0}) {
		t.Error("Position with negative row should be invalid")
	}
	if game.isValidMousePosition(Position{Row: 0, Col: -1}) {
		t.Error("Position with negative column should be invalid")
	}
	if game.isValidMousePosition(Position{Row: GridHeight, Col: 0}) {
		t.Error("Position beyond grid height should be invalid")
	}
	if game.isValidMousePosition(Position{Row: 0, Col: GridWidth}) {
		t.Error("Position beyond grid width should be invalid")
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
}package game

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