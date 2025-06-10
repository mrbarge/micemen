package game

// Constants for game configuration
const (
	GridWidth      = 19
	GridHeight     = 13
	MinWalls       = 5
	MaxWalls       = 8
	MicePerPlayer  = 12
	Player1Columns = 9 // Left-most columns for Player 1
	Player2Columns = 9 // Right-most columns for Player 2
)

// CellType represents what's in a grid cell
type CellType int

const (
	Empty CellType = iota
	Wall
)

// PlayerColor represents the player colors
type PlayerColor int

const (
	Red PlayerColor = iota
	Blue
)

// String returns the string representation of a player color
func (p PlayerColor) String() string {
	switch p {
	case Red:
		return "Red"
	case Blue:
		return "Blue"
	default:
		return "Unknown"
	}
}

// Position represents a coordinate in the grid
type Position struct {
	Row int
	Col int
}

// Mouse represents a mouse with its position and owner
type Mouse struct {
	Position Position
	Player   PlayerColor
}

// Action represents player actions
type Action int

const (
	ActionNone Action = iota
	ActionMoveLeft
	ActionMoveRight
	ActionMoveColumnUp
	ActionMoveColumnDown
	ActionQuit
)

// GameState represents the current state of the game
type GameState struct {
	Grid           [GridHeight][GridWidth]CellType
	SelectedColumn int
	GameOver       bool
	CurrentPlayer  PlayerColor
	Mice           []Mouse
}

// Player represents a player in the game
type Player struct {
	Color PlayerColor
	Mice  []Mouse
}

// Game interface defines the core game operations
type Game interface {
	GetState() GameState
	ProcessAction(action Action)
	IsGameOver() bool
	Reset()
	GetPlayer(color PlayerColor) Player
	GetMiceAt(pos Position) []Mouse
	CanPlayerMoveColumn(player PlayerColor, col int) bool
	GetValidColumnsForPlayer(player PlayerColor) []int
}

// Renderer interface for displaying the game
type Renderer interface {
	Render(state GameState)
	Clear()
	ShowMessage(msg string)
}

// InputHandler interface for getting player input
type InputHandler interface {
	Initialize() error
	GetNextAction() (Action, error)
	Close() error
}
