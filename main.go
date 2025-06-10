package main

import (
	"fmt"
	"os"

	"micemen/game"
	"micemen/input"
	"micemen/render"
)

// GameEngine coordinates the game components
type GameEngine struct {
	game   game.Game
	render game.Renderer
	input  game.InputHandler
}

// NewGameEngine creates a new game engine with all components
func NewGameEngine() *GameEngine {
	gameInstance := game.NewGame()
	return &GameEngine{
		game:   gameInstance,
		render: render.NewTerminalRenderer(gameInstance), // Pass game to renderer
		input:  input.NewKeyboardHandler(),
	}
}

// Run executes the main game loop
func (e *GameEngine) Run() error {
	// Initialize input handler
	if err := e.input.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize input: %w", err)
	}
	defer e.input.Close()

	// Set up terminal
	if termRender, ok := e.render.(*render.TerminalRenderer); ok {
		termRender.HideCursor()
		defer func() {
			termRender.ShowCursor()
			termRender.Clear()
		}()
	}

	// Initial render
	e.render.Render(e.game.GetState())

	// Main game loop
	for !e.game.IsGameOver() {
		action, err := e.input.GetNextAction()
		if err != nil {
			return fmt.Errorf("input error: %w", err)
		}

		if action != game.ActionNone {
			e.game.ProcessAction(action)
			if !e.game.IsGameOver() {
				e.render.Render(e.game.GetState())
			}
		}
	}

	e.render.ShowMessage("Thanks for playing Micemen!")
	return nil
}

func main() {
	engine := NewGameEngine()
	if err := engine.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
