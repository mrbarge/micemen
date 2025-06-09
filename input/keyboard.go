package input

import (
	"micemen/game"

	"github.com/eiannone/keyboard"
)

// KeyboardHandler implements the InputHandler interface for keyboard input
type KeyboardHandler struct {
	initialized bool
}

// NewKeyboardHandler creates a new keyboard input handler
func NewKeyboardHandler() *KeyboardHandler {
	return &KeyboardHandler{}
}

// Initialize sets up the keyboard handler
func (h *KeyboardHandler) Initialize() error {
	if h.initialized {
		return nil
	}

	if err := keyboard.Open(); err != nil {
		return err
	}

	h.initialized = true
	return nil
}

// GetNextAction waits for and returns the next player action
func (h *KeyboardHandler) GetNextAction() (game.Action, error) {
	if !h.initialized {
		if err := h.Initialize(); err != nil {
			return game.ActionNone, err
		}
	}

	char, key, err := keyboard.GetKey()
	if err != nil {
		return game.ActionNone, err
	}

	// Handle special keys first
	switch key {
	case keyboard.KeyArrowLeft:
		return game.ActionMoveLeft, nil
	case keyboard.KeyArrowRight:
		return game.ActionMoveRight, nil
	case keyboard.KeyArrowUp:
		return game.ActionMoveColumnUp, nil
	case keyboard.KeyArrowDown:
		return game.ActionMoveColumnDown, nil
	case keyboard.KeyCtrlC:
		return game.ActionQuit, nil
	}

	// Handle character input (including tmux-friendly alternatives)
	switch char {
	case 'q', 'Q':
		return game.ActionQuit, nil
	case 27: // ESC character
		return game.ActionQuit, nil
	// Alternative controls for tmux compatibility
	case 'a', 'A': // Move left
		return game.ActionMoveLeft, nil
	case 'd', 'D': // Move right
		return game.ActionMoveRight, nil
	case 'w', 'W': // Move column up
		return game.ActionMoveColumnUp, nil
	case 's', 'S': // Move column down
		return game.ActionMoveColumnDown, nil
	case 'h': // Vi-style left
		return game.ActionMoveLeft, nil
	case 'l': // Vi-style right
		return game.ActionMoveRight, nil
	case 'k': // Vi-style up
		return game.ActionMoveColumnUp, nil
	case 'j': // Vi-style down
		return game.ActionMoveColumnDown, nil
	}

	return game.ActionNone, nil
}

// Close shuts down the keyboard handler
func (h *KeyboardHandler) Close() error {
	if h.initialized {
		keyboard.Close()
		h.initialized = false
	}
	return nil
}
