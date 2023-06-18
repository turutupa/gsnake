package gsnake

import (
	"testing"
)

func TestNewGameDefaults(t *testing.T) {
	rows := 10
	cols := 10
	term := NewTerm()
	scoreboard, _ := NewScoreboard()
	screen := NewScreen(rows, cols)
	fruit := NewFruit(rows, cols)
	snake := NewSnake(screen)
	game := NewGame(term, screen, scoreboard, fruit, snake, NORMAL)

	if game.screen == nil {
		t.Errorf("screen not initialized")
	}
	if game.fruit == nil {
		t.Errorf("snake not initialized")
	}
	if game.snake == nil {
		t.Errorf("snake not initialized")
	}
	if game.speed != int(NORMAL) {
		t.Errorf("speed incorrect")
	}
}

// TODO
func TestExecuteUserInput(t *testing.T) {
}
