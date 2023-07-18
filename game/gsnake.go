package gsnake

import (
	"turutupa/gsnake/events"
)

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
		if g.state.Get() == IN_GAME_SINGLE {
			switch g.state.gameMode {
			case SINGLE:
				g.runSolo()
			case MULTI:
				g.runMulti()
			}
			g.state.SetState(MAIN_MENU)
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
}

func (g *Gsnake) runSolo() {
	g.screen.Clear()
	board := NewBoard(ROWS_SOLO, COLS_SOLO)
	leaderboard := NewLeaderboard()
	fruit := NewFruit(ROWS_SOLO, COLS_SOLO)
	snake := NewSnake(board.rows/2, board.cols/5)
	game := NewSoloGame(board, g.screen, leaderboard, fruit, snake)
	g.eventBus.Subscribe(IN_GAME_SINGLE, game.Strategy)
	game.SetDifficulty(g.state.difficulty)
	game.Run()
	game.Leaderboard()
	game.Restart()
}

func (g *Gsnake) runMulti() {
	g.screen.Clear()
	player := NewPlayer("").WithUUID().WithScreen(g.screen)
	g.eventBus.Subscribe(PLAYER_NAME_SUBMIT, player.SubmitNameStrategy)
	g.state.SetState(PLAYER_NAME_SUBMIT)
	player.SetName()
	g.eventBus.Subscribe(LOADING, func(event rune) {
		if event == 'q' {
			g.state.SetState(QUIT)
		}
	})
	g.state.SetState(LOADING)
	lobby := NewLobby()
	room := lobby.Join(player)
	exitGame := make(chan bool)
	multiStrategy := func(event rune) {
		if event == 'q' {
			room.OnExit(player)
			exitGame <- true
		}
		player.MultiGameStrategy(event)
	}
	g.eventBus.Subscribe(IN_GAME_MULTI, multiStrategy)
	g.state.SetState(IN_GAME_MULTI)
	select {
	case <-room.notifyGameEnd:
	// game is over!
	case <-exitGame:
		// user wants to leave!
	}
}
