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

func NewOnlineGsnake(
	eventPoller events.EventPoller,
	screen *Screen,
) *Gsnake {
	state := NewState()
	leaderboard := NewLeaderboard()
	rows := screen.rows
	cols := screen.cols
	return newGsnake(
		state,
		NewEventBus(state, eventPoller),
		screen,
		leaderboard,
		NewOnlineMenu(state, screen),
		NewGame(screen, leaderboard, NewFruit(rows, cols), NewSnake(screen)),
	)
}

func NewLocalGsnake(
	screen *Screen,
) *Gsnake {
	state := NewState()
	leaderboard := NewLeaderboard()
	rows := screen.rows
	cols := screen.cols
	return newGsnake(
		state,
		NewEventBus(state, NewTerm()),
		screen,
		leaderboard,
		NewOnlineMenu(state, screen),
		NewGame(screen, leaderboard, NewFruit(rows, cols), NewSnake(screen)),
	)
}

func (g *Gsnake) Run() {
	go g.eventBus.Run()
	for g.state.get() != QUIT {
		for g.state.get() == MAIN_MENU {
			g.menu.Run()
		}
		if g.state.get() == IN_GAME {
			g.screen.clear()
			g.game.setDifficulty(g.state.diff)
			g.game.Run()
			g.game.Restart()
			g.state.set(MAIN_MENU)
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
