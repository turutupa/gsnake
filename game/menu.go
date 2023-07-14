package gsnake

// Main Menu options
const (
	BACK          string = "BACK"
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
var OFFLINE_MENU_OPTIONS = []string{EASY, NORMAL, HARD, INSANITY, BACK}
var ONLINE_MENU_OPTIONS = []string{SINGLE_PLAYER, MULTI_PLAYER, EXIT}

type MenuState int

const (
	SELECT_GAME_MODE   MenuState = 1
	SINGLE_MODE        MenuState = 2
	PROMPT_PLAYER_NAME MenuState = 3
)

type Menu struct {
	appState           *State
	menuState          MenuState
	screen             *Screen
	selectedMenuOption int
	keypressCh         chan bool
}

func newMenu(
	appState *State,
	menuState MenuState,
	screen *Screen,
	selectedMenuOption int,
) *Menu {
	return &Menu{
		appState:           appState,
		menuState:          menuState,
		screen:             screen,
		selectedMenuOption: selectedMenuOption,
		keypressCh:         make(chan bool),
	}
}

func NewLocalMenu(appState *State, screen *Screen) *Menu {
	return newMenu(appState, SINGLE_MODE, screen, 1) // 1 defaults to NORMAL SPEED
}

func NewOnlineMenu(appState *State, screen *Screen) *Menu {
	return newMenu(appState, SELECT_GAME_MODE, screen, 0) // 0 defaults to SINGLE PLAYER
}

func (m *Menu) Run() {
	m.screen.Clear()
	switch m.menuState {
	case SELECT_GAME_MODE:
		title := "SELECT GAME MODE"
		termSize := *m.appState.termSize
		m.screen.RenderMenu(termSize, title, ONLINE_MENU_OPTIONS, m.selectedMenuOption)
	case SINGLE_MODE:
		title := "SELECT DIFFICULTY"
		termSize := *m.appState.termSize
		m.screen.RenderMenu(termSize, title, OFFLINE_MENU_OPTIONS, m.selectedMenuOption)
	}
	<-m.keypressCh
}

func (m *Menu) Strategy(event rune) {
	var options []string
	switch m.menuState {
	case SINGLE_MODE:
		options = OFFLINE_MENU_OPTIONS
	case SELECT_GAME_MODE:
		options = ONLINE_MENU_OPTIONS
	}
	if isUp(event) {
		m.selectedMenuOption = m.selectedMenuOption - 1
		if m.selectedMenuOption < 0 {
			m.selectedMenuOption = len(options) - 1
		}
	} else if isDown(event) {
		m.selectedMenuOption = m.selectedMenuOption + 1
		if m.selectedMenuOption >= len(options) {
			m.selectedMenuOption = 0
		}
	} else if isEnterKey(event) {
		selectedOpt := options[m.selectedMenuOption]
		switch selectedOpt {
		case EXIT:
			m.appState.SetState(QUIT)
		case BACK:
			m.menuState = SELECT_GAME_MODE
			m.selectedMenuOption = 0
		case SINGLE_PLAYER:
			m.menuState = SINGLE_MODE
			m.selectedMenuOption = 1
		case MULTI_PLAYER:
			m.appState.SetState(IN_GAME_SINGLE).SetGameMode(MULTI)
		default:
			gameOpts := m.appState.SetState(IN_GAME_SINGLE).SetGameMode(SINGLE)
			switch selectedOpt {
			case EASY:
				gameOpts.SetDifficulty(EASY)
			case NORMAL:
				gameOpts.SetDifficulty(NORMAL)
			case HARD:
				gameOpts.SetDifficulty(HARD)
			case INSANITY:
				gameOpts.SetDifficulty(INSANITY)
			}
		}
	} else if event == 'q' {
		m.appState.SetState(QUIT)
	}
	m.keypressCh <- true
}
