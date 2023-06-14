package gsnake

import (
	"testing"
)

func TestNewGameDefaults(t *testing.T) {
	width := 10
	height := 10
	term := NewTerm()
	screen := NewScreen(width, height)
	snake := NewSnake(screen)
	game := NewGame(term, screen, snake, Normal)

	if game.screen == nil {
		t.Errorf("screen not initialized")
	}

	if game.snake == nil {
		t.Errorf("snake not initialized")
	}

	if game.started {
		t.Errorf("game should not have started")
	}

	if game.finished {
		t.Errorf("game should not have finished")
	}
}

func TestExecuteUserInput(t *testing.T) {
}
