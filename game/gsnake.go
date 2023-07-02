package gsnake

import (
	"turutupa/gsnake/events"
)

type Gsnake struct {
	state       *State
	eventBus    *EventBus
	screen      *Screen
	leaderboard *Leaderboard
	menu        *Menu
	game        *SoloGame
}

func newGsnake(
	state *State,
	eventBus *EventBus,
	board *Board,
	screen *Screen,
	leaderboard *Leaderboard,
	menu *Menu,
	game *SoloGame,
) *Gsnake {
	eventBus.Subscribe(MAIN_MENU, menu.Strategy)
	eventBus.Subscribe(IN_GAME, game.Strategy)
	return &Gsnake{
		state:       state,
		eventBus:    eventBus,
		screen:      screen,
		leaderboard: leaderboard,
		menu:        menu,
		game:        game,
	}
}

// injection wrapper to create single player game modes
func NewLocalGsnake(screen *Screen) *Gsnake {
	const rows = 20
	const cols = 50
	state := NewState()
	board := NewBoard(rows, cols)
	eventBus := NewEventBus(state, NewTerm())
	leaderboard := NewLeaderboard()
	menu := NewLocalMenu(state, board, screen)
	game := NewGame(board, screen, leaderboard, NewFruit(rows, cols), NewSnake(board))

	return newGsnake(
		state,
		eventBus,
		board,
		screen,
		leaderboard,
		menu,
		game,
	)
}

// injection wrapper to create multiplayer game modes
func NewMultiplayerGsnake(eventPoller events.EventPoller, screen *Screen) *Gsnake {
	rows := 30
	cols := 80
	board := NewBoard(rows, cols)
	state := NewState()
	eventBus := NewEventBus(state, eventPoller)
	leaderboard := NewLeaderboard()
	menu := NewOnlineMenu(state, board, screen)
	game := NewGame(board, screen, leaderboard, NewFruit(rows, cols), NewSnake(board))

	return newGsnake(
		state,
		eventBus,
		board,
		screen,
		leaderboard,
		menu,
		game,
	)
}

func (g *Gsnake) Run() {
	g.screen.HideCursor()
	go g.eventBus.Run()
	for g.state.Get() != QUIT {
		for g.state.Get() == MAIN_MENU {
			g.menu.Run()
		}
		if g.state.Get() == IN_GAME {
			switch g.state.gameMode {
			case SINGLE:
				g.screen.Clear()
				g.game.SetDifficulty(g.state.difficulty)
				g.game.Run()
				g.game.Leaderboard()
				g.game.Restart()
				g.state.SetState(MAIN_MENU)
			case MULTI:
				g.screen.Clear()
				// g.screen.promptPlayerName()
				g.game.SetDifficulty(g.state.difficulty)
				g.game.Run()
				g.game.Leaderboard()
				g.game.Restart()
				g.state.SetState(MAIN_MENU)
			}
		}
	}
	g.Stop()
}

func (g *Gsnake) Stop() {
	g.game.Stop()
	g.eventBus.Stop()
	g.screen.Clear()
	g.screen.ShowCursor()
}
