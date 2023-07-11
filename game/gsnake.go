package gsnake

import (
	"turutupa/gsnake/events"
)

const ROWS_SOLO = 20
const COLS_SOLO = 50

const ROWS_MULTI = 30
const COLS_MULTI = 50

type Gsnake struct {
	state    *State
	eventBus *EventBus
	screen   *Screen
	menu     *Menu
}

func NewGsnake(
	state *State,
	eventBus *EventBus,
	screen *Screen,
	menu *Menu,
) *Gsnake {
	eventBus.Subscribe(MAIN_MENU, menu.Strategy)
	return &Gsnake{
		state:    state,
		eventBus: eventBus,
		screen:   screen,
		menu:     menu,
	}
}

// injection wrapper to create single player game modes
func NewOfflineGsnake(screen *Screen) *Gsnake {
	state := NewState()
	eventBus := NewEventBus(state, NewTerm())
	menu := NewLocalMenu(state, screen)

	return NewGsnake(
		state,
		eventBus,
		screen,
		menu,
	)
}

// injection wrapper to create multiplayer game modes
func NewOnlineGsnake(eventPoller events.EventPoller, screen *Screen) *Gsnake {
	state := NewState()
	eventBus := NewEventBus(state, eventPoller)
	menu := NewOnlineMenu(state, screen)

	return NewGsnake(
		state,
		eventBus,
		screen,
		menu,
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
				board := NewBoard(ROWS_SOLO, COLS_SOLO)
				leaderboard := NewLeaderboard()
				fruit := NewFruit(ROWS_SOLO, COLS_SOLO)
				snake := NewSnake(board)
				game := NewSoloGame(board, g.screen, leaderboard, fruit, snake)
				g.eventBus.Subscribe(IN_GAME, game.Strategy)
				game.SetDifficulty(g.state.difficulty)
				game.Run()
				game.Leaderboard()
				game.Restart()
				g.state.SetState(MAIN_MENU)
			case MULTI:
				g.screen.Clear()
				g.screen.PromptPlayerName()
				// g.game.SetDifficulty(g.state.difficulty) // this has to go
				// g.game.Run()
				// g.game.Leaderboard()
				// g.game.Restart()
				g.state.SetState(MAIN_MENU)
			}
		}
	}
	g.Stop()
}

func (g *Gsnake) Stop() {
	g.eventBus.Stop()
	g.screen.Clear()
	g.screen.ShowCursor()
}

func (g *Gsnake) OnWindowChange(
	wc struct {
		Width       uint32
		Height      uint32
		PixelWidth  uint32
		PixelHeight uint32
	}) {
	g.state.termSize.rows = int(wc.Height)
	g.state.termSize.cols = int(wc.Width)

	g.screen.SetSize(int(wc.Height), int(wc.Width))

	// g.screen.Clear()
	// if int(wc.Height) < g.board.rows || int(wc.Width) < g.board.cols {
	// 	g.screen.RenderWarning(int(wc.Height), int(wc.Width), "MAKE YOUR TERMINAL BIGGER TO RENDER GAME")
	// 	// g.screen.
	// } else {
	// 	g.screen.SetOffset(offsetRows, offsetCols)
	// 	g.screen.RenderBoard(g.board)
	// }
}
