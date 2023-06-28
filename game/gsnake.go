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
	game        *Game
}

func newGsnake(
	state *State,
	eventBus *EventBus,
	screen *Screen,
	leaderboard *Leaderboard,
	menu *Menu,
	game *Game,
) *Gsnake {
	eventBus.subscribe(MAIN_MENU, menu.strategy)
	eventBus.subscribe(IN_GAME, game.strategy)
	return &Gsnake{
		state:       state,
		eventBus:    eventBus,
		screen:      screen,
		leaderboard: leaderboard,
		menu:        menu,
		game:        game,
	}
}

func NewOnlineGsnake(eventPoller events.EventPoller, screen *Screen) *Gsnake {
	rows := screen.rows
	cols := screen.cols
	state := NewState()
	eventBus := NewEventBus(state, eventPoller)
	leaderboard := NewLeaderboard()
	menu := NewOnlineMenu(state, screen)
	game := NewGame(screen, leaderboard, NewFruit(rows, cols), NewSnake(screen))

	return newGsnake(
		state,
		eventBus,
		screen,
		leaderboard,
		menu,
		game,
	)
}

func NewLocalGsnake(screen *Screen) *Gsnake {
	rows := screen.rows
	cols := screen.cols
	state := NewState()
	eventBus := NewEventBus(state, NewTerm())
	leaderboard := NewLeaderboard()
	menu := NewLocalMenu(state, screen)
	game := NewGame(screen, leaderboard, NewFruit(rows, cols), NewSnake(screen))

	return newGsnake(
		state,
		eventBus,
		screen,
		leaderboard,
		menu,
		game,
	)
}

func (g *Gsnake) Run() {
	go g.eventBus.Run()
	for g.state.get() != QUIT {
		for g.state.get() == MAIN_MENU {
			g.menu.Run()
		}
		if g.state.get() == IN_GAME {
			switch g.state.gameMode {
			case SINGLE:
				g.screen.clear()
				g.game.setDifficulty(g.state.difficulty)
				g.game.Run()
				g.game.Restart()
				g.state.setState(MAIN_MENU)
			case MULTI:
				g.screen.clear()
				g.game.setDifficulty(g.state.difficulty)
				g.game.Run()
				g.game.Restart()
				g.state.setState(MAIN_MENU)
			}
		}
	}
	g.Stop()
}

func (g *Gsnake) Stop() {
	g.game.Stop()
	g.eventBus.Stop()
	g.screen.clear()
	g.screen.showCursor()
}
