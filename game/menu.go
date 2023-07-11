package gsnake

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
	m.screen.Clear()
	m.screen.RenderMainMenu(*m.state.termSize, m.selectedMenuOption)
	<-m.keypressCh
}

func (m *Menu) Strategy(event rune) {
	if isUp(event) {
		m.selectedMenuOption = m.selectedMenuOption - 1
		if m.selectedMenuOption < 0 {
			m.selectedMenuOption = len(MENU_OPTIONS) - 1
		}
	} else if isDown(event) {
		m.selectedMenuOption = m.selectedMenuOption + 1
		if m.selectedMenuOption >= len(MENU_OPTIONS) {
			m.selectedMenuOption = 0
		}
	} else if isEnterKey(event) {
		selectedOpt := MENU_OPTIONS[m.selectedMenuOption]
		if selectedOpt == EXIT {
			m.state.SetState(QUIT)
		} else {
			gameOpts := m.state.SetState(IN_GAME).SetGameMode(SINGLE)
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
		m.state.SetState(QUIT)
	}
	m.keypressCh <- true
}
