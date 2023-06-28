package gsnake

import (
	"math"
)

// Main Menu options
const (
	EXIT          string = "EXIT"
	LEADERBOARD   string = "LEADERBOARD"
	EASY          string = "EASY"
	NORMAL        string = "NORMAL"
	HARD          string = "HARD"
	INSANITY      string = "INSANITY"
	SINGLE_PLAYER string = "SINGLE_PLAYER"
	MULTI_PLAYER  string = "MULTI_PLAYER"
)

var DIFFICULTIES = []string{EASY, NORMAL, HARD, INSANITY}
var MENU_OPTIONS = []string{EASY, NORMAL, HARD, INSANITY, EXIT}
var SSH_MENU_OPTIONS = []string{SINGLE_PLAYER, MULTI_PLAYER}

type Menu struct {
	state              *State
	screen             *Screen
	selectedMenuOption int
	keypressCh         chan bool
}

func newMenu(
	state *State,
	screen *Screen,
	selectedMenuOption int,
) *Menu {
	return &Menu{
		state:              state,
		screen:             screen,
		selectedMenuOption: selectedMenuOption,
		keypressCh:         make(chan bool),
	}
}

func NewLocalMenu(state *State, screen *Screen) *Menu {
	return newMenu(state, screen, 1) // 1 defaults to NORMAL SPEED
}

func NewOnlineMenu(state *State, screen *Screen) *Menu {
	return newMenu(state, screen, 0) // 0 defaults to SINGLE PLAYER
}

func (m *Menu) Run() {
	m.screen.clear()
	m.screen.renderMainMenu(m.selectedMenuOption)
	<-m.keypressCh
}

func (m *Menu) strategy(event rune) {
	if isUp(event) {
		m.selectedMenuOption = int(math.Max(float64(0), float64(m.selectedMenuOption-1)))
	} else if isDown(event) {
		m.selectedMenuOption = int(math.Min(float64(len(MENU_OPTIONS)-1), float64(m.selectedMenuOption+1)))
	} else if isEnterKey(event) {
		selectedOpt := MENU_OPTIONS[m.selectedMenuOption]
		if selectedOpt == EXIT {
			m.state.setState(QUIT)
		} else {
			game := m.state.setState(IN_GAME).setGameMode(SINGLE)
			switch selectedOpt {
			case EASY:
				game.setDifficulty(EASY)
			case NORMAL:
				game.setDifficulty(NORMAL)
			case HARD:
				game.setDifficulty(HARD)
			case INSANITY:
				game.setDifficulty(INSANITY)
			}
		}
	} else if event == 'q' {
		m.state.setState(QUIT)
	}
	m.keypressCh <- true
}
