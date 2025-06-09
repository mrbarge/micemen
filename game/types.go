package game

// Constants for game configuration
const (
	GridWidth  = 19
	GridHeight = 13
	MinWalls   = 5
	MaxWalls   = 8
)

// CellType represents what's in a grid cell
type CellType int

const (
	Empty CellType = iota
	Wall
)

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
}

// Game interface defines the core game operations
type Game interface {
	GetState() GameState
	ProcessAction(action Action)
	IsGameOver() bool
	Reset()
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
