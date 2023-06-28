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

type State struct {
	state       AppState
	gameMode    GameMode
	difficulty  string
	subscribers []func(AppState)
}

func NewState() *State {
	return &State{
		state:       MAIN_MENU, // defaults to MAIN MENU
		gameMode:    SINGLE,    // defaults to single player
		difficulty:  NORMAL,    // defaults to NORMAL
		subscribers: []func(AppState){},
	}
}

func (sb *State) setState(appState AppState) *State {
	sb.state = appState
	return sb
}

func (sb *State) setGameMode(gameMode GameMode) *State {
	sb.gameMode = gameMode
	return sb
}

func (sb *State) setDifficulty(difficulty string) *State {
	sb.difficulty = difficulty
	return sb
}

func (sb *State) get() AppState {
	return sb.state
}
