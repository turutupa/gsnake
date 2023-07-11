package gsnake

// App AppState
type AppState string
type GameMode string

const (
	SINGLE GameMode = "SINGLE"
	MULTI  GameMode = "MULTI"
)

const (
	MAIN_MENU AppState = "MAIN_MENU" // starts here by default
	IN_GAME   AppState = "IN_GAME"
	QUIT      AppState = "QUIT"
)

type TermSize struct {
	rows int
	cols int
}

type State struct {
	state      AppState
	gameMode   GameMode
	difficulty string
	termSize   *TermSize
}

func NewState() *State {
	return &State{
		state:      MAIN_MENU, // defaults to MAIN MENU
		gameMode:   SINGLE,    // defaults to single player
		difficulty: NORMAL,    // defaults to NORMAL
		termSize:   &TermSize{0, 0},
	}
}

func (sb *State) Get() AppState {
	return sb.state
}

func (sb *State) SetState(appState AppState) *State {
	sb.state = appState
	return sb
}

func (sb *State) SetGameMode(gameMode GameMode) *State {
	sb.gameMode = gameMode
	return sb
}

func (sb *State) SetDifficulty(difficulty string) *State {
	sb.difficulty = difficulty
	return sb
}
